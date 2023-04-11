package pkg

type BrightnessConfig struct {
	Set            bool `json:"set"`
	TargetBrightness int  `json:"target_brightness"`
}

type Config struct {
	Sunrise BrightnessConfig `json:"sunrise"`
	Sunset  BrightnessConfig `json:"sunset"`
	Lat     float64          `json:"lat"`
	Long    float64          `json:"long"`
}
