package entity

type Audio struct {
	VideoID string `yaml:"video_id"`
	Length  string `yaml:"length"`
	Size    int    `yaml:"size"`
	URL     string `yaml:"url"`
}
