package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/xaxys/maintainman/core/config"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/util"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/spf13/viper"
)

var (
	Storage IStorage
	s3Conn  *s3.S3
)

type IStorage interface {
	Path() string
	Exist(id string) bool
	Load(id string, fn func(io.Reader) error) error
	Save(id, format string, fn func(io.Writer) error) error
	LoadBytes(id string) ([]byte, error)
	SaveBytes(id, format string, data []byte) error
	Delete(id string) error
	Sub(path string, clean bool) IStorage
}

func init() {
	s3Conn, _ = initS3Conn(config.AppConfig)
	Storage = InitStorage(config.AppConfig)
}

func InitStorage(config *viper.Viper) (storage IStorage) {
	storageType := config.GetString("storage.driver")
	switch storageType {
	case "":
		return nil
	case "local":
		path := config.GetString("storage.local.path")
		clean := config.GetBool("storage.local.clean")
		storage = newLocalStorage(path, clean)
	case "s3":
		conn, err := initS3Conn(config)
		if err != nil {
			if s3Conn == nil {
				panic(fmt.Errorf("no s3 connection specified in both config and env: %+v", err))
			}
			fmt.Printf("no s3 connection specified, use default connection: %+v", err)
			conn = s3Conn
		}
		bucket := config.GetString("storage.s3.bucket")
		path := config.GetString("storage.s3.path")
		clean := config.GetBool("storage.s3.clean")
		storage = newS3Storage(conn, bucket, path, clean)
	default:
		panic(fmt.Errorf("support local and s3 only"))
	}
	return storage
}

func initS3Conn(config *viper.Viper) (*s3.S3, error) {
	accessKey := config.GetString("storage.s3.access_key")
	secretKey := config.GetString("storage.s3.secret_key")
	regionInfo := config.GetString("storage.s3.region")
	auth, err := aws.GetAuth(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("aws cannot get auth: %v", err)
	}
	region, ok := aws.Regions[regionInfo]
	if !ok {
		return nil, fmt.Errorf("invalid region: %v", regionInfo)
	}
	return s3.New(auth, region), nil
}

type LocalStorage struct {
	path string
}

func newLocalStorage(path string, clean bool) *LocalStorage {
	path = filepath.Clean(path)
	if clean {
		os.RemoveAll(path)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	return &LocalStorage{
		path: path,
	}
}

func (s *LocalStorage) Path() string {
	return s.path
}

func (s *LocalStorage) Exist(id string) bool {
	fullPath := filepath.Join(s.path, id)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *LocalStorage) Load(id string, fn func(io.Reader) error) (err error) {
	fullPath := filepath.Join(s.path, id)
	reader, err := os.Open(fullPath)
	defer reader.Close()
	if err != nil {
		return fmt.Errorf("failed to load %s: %v", id, err)
	}
	err = fn(reader)
	return err
}

func (s *LocalStorage) LoadBytes(id string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := s.Load(id, func(reader io.Reader) error {
		_, err := io.Copy(buffer, reader)
		return err
	}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *LocalStorage) Save(id, format string, fn func(io.Writer) error) error {
	// Open file for writing, overwrite if it already exists
	fullPath := filepath.Join(s.path, id)
	writer, err := os.Create(fullPath)
	defer writer.Close()
	if err != nil {
		return err
	}
	if err := fn(writer); err != nil {
		return err
	}
	return nil
}

func (s *LocalStorage) SaveBytes(id, format string, data []byte) error {
	return s.Save(id, format, func(writer io.Writer) error {
		_, err := writer.Write(data)
		return err
	})
}

func (s *LocalStorage) Delete(id string) error {
	fullPath := filepath.Join(s.path, id)
	return os.Remove(fullPath)
}

func (s *LocalStorage) Sub(path string, clean bool) IStorage {
	subPath := filepath.Join(s.path, path)
	return newLocalStorage(subPath, clean)
}

// S3Storage is a storage implementation using Amazon S3
type S3Storage struct {
	path   string
	bucket *s3.Bucket
}

func newS3Storage(conn *s3.S3, bucket, path string, clean bool) *S3Storage {
	bucketObj := conn.Bucket(bucket)
	path = filepath.Clean(path)
	storage := &S3Storage{
		path:   path,
		bucket: bucketObj,
	}
	if clean {
		storage.Clean()
	}
	return storage
}

func (s *S3Storage) Path() string {
	return s.path
}

func (s *S3Storage) Exist(id string) bool {
	fullPath := s.path + "/" + id
	resp, err := s.bucket.List(fullPath, "/", "", 10)
	if err != nil {
		logger.Logger.Errorf("Error while listing S3 bucket: %v\n", err)
		return false
	}
	if resp == nil {
		logger.Logger.Error("Error while listing S3 bucket: empty response")
	}

	for _, element := range resp.Contents {
		if element.Key == id {
			return true
		}
	}

	return false
}

func (s *S3Storage) Load(id string, fn func(io.Reader) error) (err error) {
	fullPath := s.path + "/" + id
	rc, err := s.bucket.GetReader(fullPath)
	defer rc.Close()
	if err != nil {
		return err
	}
	err = fn(rc)
	return err
}

func (s *S3Storage) LoadBytes(id string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := s.Load(id, func(reader io.Reader) error {
		_, err := io.Copy(buffer, reader)
		return err
	}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *S3Storage) Save(id, format string, fn func(io.Writer) error) error {
	buffer := bytes.NewBuffer(nil)
	if err := fn(buffer); err != nil {
		return err
	}
	return s.SaveBytes(id, format, buffer.Bytes())
}

func (s *S3Storage) SaveBytes(id, format string, data []byte) error {
	fullPath := s.path + "/" + id
	return s.bucket.Put(fullPath, data, format, s3.Private)
}

func (s *S3Storage) Delete(id string) error {
	fullPath := s.path + "/" + id
	return s.bucket.Del(fullPath)
}

func (s *S3Storage) Sub(path string, clean bool) IStorage {
	subPath := s.path + "/" + path
	return newS3Storage(s.bucket.S3, s.bucket.Name, subPath, clean)
}

func (s *S3Storage) Clean() error {
	// Delete all objects in the bucket
	resp, err := s.bucket.List(s.path, "/", "", 1000)
	if err != nil {
		return err
	}
	for len(resp.Contents) > 0 {
		keys := util.TransSlice(resp.Contents, func(k s3.Key) string { return k.Key })
		if err := s.bucket.MultiDel(keys); err != nil {
			return err
		}
		resp, err = s.bucket.List(s.path, "/", "", 1000)
		if err != nil {
			return err
		}
	}
	return nil
}
