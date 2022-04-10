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
	ImageConfig.SetDefault("gif_num_colors", 256)
	ImageConfig.SetDefault("cache_as_jpeg", true)
	ImageConfig.SetDefault("save_as_jpeg", false)

	ImageConfig.SetDefault("upload.async", false)
	ImageConfig.SetDefault("upload.throttling.burst", 20)
	ImageConfig.SetDefault("upload.throttling.rate", 1)
	ImageConfig.SetDefault("upload.throttling.purge", "1m")
	ImageConfig.SetDefault("upload.throttling.expire", "1m")
	ImageConfig.SetDefault("upload.max_file_size", 10485760) // 10M
	ImageConfig.SetDefault("upload.max_pixels", 15000000)    // 15M pixels

	ImageConfig.SetDefault("cache.driver", "local")
	ImageConfig.SetDefault("cache.limit", 1073741824) // 1GB

	ImageConfig.SetDefault("storage.driver", "local")
	ImageConfig.SetDefault("storage.local.path", "./images/cache")
	ImageConfig.SetDefault("storage.local.clean", true)
	ImageConfig.SetDefault("storage.s3.bucket", "BUCKET")

	ImageConfig.SetDefault("transformations", []map[string]any{
		{
			"name":   "sw-corner",
			"params": "w_100,h_100,c_k,g_sw",
		},
		{
			"name":   "square",
			"params": "w_200,h_200",
			"eager":  true,
		},
		{
			"name":    "watermarked",
			"params":  "w_800",
			"default": true,
			"texts": []map[string]any{
				{
					"content": "{{.Name}}@MaintainMan",
					"gravity": "se",
					"x-pos":   10,
					"y-pos":   0,
					"color":   "#808080CC",
					"font":    "fonts/SourceHanSans-Regular.ttf",
					"size":    14,
				},
			},
		},
	})

	ReadAndUpdateConfig(ImageConfig, "image", ImageConfigVersion)
}
