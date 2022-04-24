package imagehost

import (
	"github.com/spf13/viper"
)

var (
	imageConfig = viper.New()
)

func init() {
	imageConfig.SetDefault("jpeg_quality", 80)
	imageConfig.SetDefault("gif_num_colors", 256)
	imageConfig.SetDefault("cache_as_jpeg", true)
	imageConfig.SetDefault("save_as_jpeg", false)

	imageConfig.SetDefault("upload.async", false)
	imageConfig.SetDefault("upload.throttling.enable", true)
	imageConfig.SetDefault("upload.throttling.burst", 20)
	imageConfig.SetDefault("upload.throttling.rate", 5)
	imageConfig.SetDefault("upload.throttling.purge", "1m")
	imageConfig.SetDefault("upload.throttling.expire", "5m")
	imageConfig.SetDefault("upload.max_file_size", 10485760) // 10M
	imageConfig.SetDefault("upload.max_pixels", 15000000)    // 15M pixels

	imageConfig.SetDefault("cache.driver", "local")
	imageConfig.SetDefault("cache.limit", 1073741824) // 1GB

	imageConfig.SetDefault("storage.driver", "local")
	imageConfig.SetDefault("storage.local.path", "./images")
	imageConfig.SetDefault("storage.s3.bucket", "BUCKET")
	imageConfig.SetDefault("storage.cache.clean", true)

	imageConfig.SetDefault("transformations", []map[string]any{
		{
			"name":   "square",
			"params": "w_256,h_256,c_p,g_c",
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
}
