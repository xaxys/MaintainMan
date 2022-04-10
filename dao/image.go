package dao

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"maintainman/cache"
	"maintainman/config"
	"maintainman/storage"
	"maintainman/util"

	"github.com/spf13/viper"
	"golang.org/x/image/bmp"
)

var (
	TransformationPO = NewTransformationPersistence(config.ImageConfig)
)

func init() {
	cache.CreateImageCache(func(a any) error {
		if id, ok := a.(string); ok {
			return DeleteImage(id)
		}
		return nil
	})
}

type TransformationPersistence struct {
	data  []util.TransformationInfo
	index map[string]*util.Transformation
	eager []*util.Transformation
	def   *util.Transformation
}

func NewTransformationPersistence(config *viper.Viper) (s *TransformationPersistence) {
	s = &TransformationPersistence{
		index: make(map[string]*util.Transformation),
		eager: []*util.Transformation{},
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

func GetTransformation(name string) (*util.Transformation, bool) {
	if name == "" {
		return TransformationPO.def, true
	}
	if name == "origin" {
		return nil, true
	}
	trans, ok := TransformationPO.index[name]
	return trans, ok
}

func GetEagerTransformation() []*util.Transformation {
	return TransformationPO.eager
}

func ExistImage(id string) bool {
	return storage.Storage.Exist(id)
}

func LoadImage(id string) (img image.Image, data []byte, format string, err error) {
	if data, err = storage.Storage.LoadBytes(id); err != nil {
		return nil, nil, "", err
	}
	img, format, err = image.Decode(bytes.NewReader(data))
	return img, data, format, err
}

func SaveImageBytes(id, format string, data []byte) error {
	return storage.Storage.SaveBytes(id, format, data)
}

func SaveImage(id, format string, img image.Image) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	imgType := ""

	switch format {
	case "jpeg", "jpg":
		imgType = "image/jpeg"
		options := &jpeg.Options{Quality: config.ImageConfig.GetInt("jpeg_quality")}
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
		options := &gif.Options{NumColors: config.ImageConfig.GetInt("gif_num_colors")}
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
		// data, _ := ioutil.ReadFile(file_name)
		// m, err := webp.Decode(bytes.NewReader(data))
		// if err == nil {
		//     var buf bytes.Buffer
		//     webp.Encode(&buf, m, nil)
		//     ioutil.WriteFile(`./_test.webp`, buf.Bytes(), 0666)
		// }
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	data := buffer.Bytes()
	return data, SaveImageBytes(id, imgType, data)
}

func DeleteImage(id string) error {
	return storage.Storage.Delete(id)
}
