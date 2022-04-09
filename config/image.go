package config

import (
	"github.com/spf13/viper"
)

const ImageConfigVersion = "1.0.0"

var (
	ImageConfig *viper.Viper
)

func init() {
	ImageConfig = viper.New()
	ImageConfig.SetConfigName("image")
	ImageConfig.SetConfigType("yaml")
	ImageConfig.AddConfigPath(".")
	ImageConfig.AddConfigPath("./config")
	ImageConfig.AddConfigPath("/etc/srs_wrappper/")
	ImageConfig.AddConfigPath("$HOME/.srs_wrappper/")

	ImageConfig.SetDefault("jpeg_quality", 80)

	ImageConfig.SetDefault("upload.async", false)
	ImageConfig.SetDefault("upload.throttling.burst", 20)
	ImageConfig.SetDefault("upload.throttling.rate", 1)
	ImageConfig.SetDefault("upload.throttling.purge", "1m")
	ImageConfig.SetDefault("upload.throttling.expire", "1m")
	ImageConfig.SetDefault("upload.max_file_size", 10485760) // 10M
	ImageConfig.SetDefault("upload.max_pixels", 8000000)     // 8M pixels

	ImageConfig.SetDefault("cache.driver", "local")
	ImageConfig.SetDefault("cache.limit", 1073741824) // 1GB
	ImageConfig.SetDefault("transformations", []map[string]any{
		{
			"name":   "sw-corner",
			"params": "w_100,h_50,c_k,g_sw",
		},
		{
			"name":   "square",
			"params": "w_200,h_200",
			"eager":  true,
		},
		{
			"name":    "watermarked",
			"params":  "w_600",
			"default": true,
			"texts": []map[string]any{
				{
					"content": "{{.Name}}@MaintainMan",
					"gravity": "ne",
					"x-pos":   30,
					"y-pos":   20,
					"color":   "#fff",
					"font":    "fonts/SourceHanSans-Regular.ttf",
					"size":    12,
				},
			},
		},
	})

	ReadAndUpdateConfig(ImageConfig, "image", ImageConfigVersion)
}
