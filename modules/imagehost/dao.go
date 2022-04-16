package imagehost

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/xaxys/maintainman/core/storage"
	"github.com/xaxys/maintainman/core/util"

	"github.com/spf13/viper"
	"golang.org/x/image/bmp"
)

var (
	transformationPO *transformationPersistence
)

type transformationPersistence struct {
	data  []TransformationInfo
	index map[string]*Transformation
	eager []*Transformation
	def   *Transformation
}

func newTransformationPersistence(config *viper.Viper) (s *transformationPersistence) {
	s = &transformationPersistence{
		index: make(map[string]*Transformation),
		eager: []*Transformation{},
	}

	config.UnmarshalKey("transformations", &s.data)
	for _, info := range s.data {
		if info.Name == "origin" {
			panic("origin transformation is reserved")
		}
		trans := info.ToTransformation()
		if s.index[info.Name] != nil {
			panic(fmt.Errorf("duplicate transformation name %s", info.Name))
		}
		s.index[info.Name] = trans
		if info.Eager {
			s.eager = append(s.eager, trans)
		}
		if info.Default {
			if s.def != nil {
				panic("default transformation already set")
			}
			s.def = trans
		}
	}

	return
}

func getTransformation(name string) (*Transformation, bool) {
	if name == "" {
		return transformationPO.def, true
	}
	if name == "origin" {
		return nil, true
	}
	trans, ok := transformationPO.index[name]
	return trans, ok
}

func getEagerTransformation() []*Transformation {
	return transformationPO.eager
}

// Storage API

var (
	imageStorage      storage.IStorage
	imageCacheStorage storage.IStorage
)

func onEvict(a any) error {
	if id, ok := a.(string); ok {
		return deleteImage(id, true)
	}
	return nil
}

func existImage(id string, cached bool) bool {
	store := util.Tenary(cached, imageCacheStorage, imageStorage)
	return stExistImage(store, id)
}

func loadImage(id string, cached bool) (img image.Image, data []byte, format string, err error) {
	store := util.Tenary(cached, imageCacheStorage, imageStorage)
	return stLoadImage(store, id)
}

func saveImageBytes(id, format string, data []byte, cached bool) error {
	store := util.Tenary(cached, imageCacheStorage, imageStorage)
	return stSaveImageBytes(store, id, format, data)
}

func saveImage(id, format string, img image.Image, cached bool) ([]byte, error) {
	store := util.Tenary(cached, imageCacheStorage, imageStorage)
	return stSaveImage(store, id, format, img)
}

func deleteImage(id string, cached bool) error {
	store := util.Tenary(cached, imageCacheStorage, imageStorage)
	return stDeleteImage(store, id)
}

// St Functions

func stExistImage(store storage.IStorage, id string) bool {
	return store.Exist(id)
}

func stLoadImage(store storage.IStorage, id string) (img image.Image, data []byte, format string, err error) {
	if data, err = store.LoadBytes(id); err != nil {
		return nil, nil, "", err
	}
	img, format, err = image.Decode(bytes.NewReader(data))
	return img, data, format, err
}

func stSaveImageBytes(store storage.IStorage, id, format string, data []byte) error {
	return store.SaveBytes(id, format, data)
}

func stSaveImage(store storage.IStorage, id, format string, img image.Image) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	imgType := ""

	switch format {
	case "jpeg", "jpg":
		imgType = "image/jpeg"
		options := &jpeg.Options{Quality: imageConfig.GetInt("jpeg_quality")}
		if err := jpeg.Encode(buffer, img, options); err != nil {
			return nil, err
		}
	case "png":
		imgType = "image/png"
		if err := png.Encode(buffer, img); err != nil {
			return nil, err
		}
	case "gif":
		imgType = "image/gif"
		options := &gif.Options{NumColors: imageConfig.GetInt("gif_num_colors")}
		if err := gif.Encode(buffer, img, options); err != nil {
			return nil, err
		}
	case "bmp":
		imgType = "image/bmp"
		if err := bmp.Encode(buffer, img); err != nil {
			return nil, err
		}
	default:
		// TODO: support webp format
		// data, _ := ioReadFile(file_name)
		// m, err := webp.Decode(bytes.NewReader(data))
		// if err == nil {
		//     var buf bytes.Buffer
		//     webp.Encode(&buf, m, nil)
		//     ioWriteFile(`./_test.webp`, buf.Bytes(), 0666)
		// }
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	data := buffer.Bytes()
	return data, stSaveImageBytes(store, id, imgType, data)
}

func stDeleteImage(store storage.IStorage, id string) error {
	return store.Delete(id)
}
