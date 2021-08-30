package api

import "github.com/spf13/viper"

func init() {
	// Server config
	viper.SetDefault("SERVER_ADDRESS", ":8000")
	viper.SetDefault("DEVELOPMENT_ENVIRONMENT", true)

	// OCR Engine config
	viper.SetDefault("OCR_ENGINE_TESSDATA_PREFIX", "/home/flavioltonon/Documents/dev/tools/tesseract/tessdata_best")
	viper.SetDefault("OCR_ENGINE_LANGUAGE", "por")
}
