package entity

import (
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"url"`
	Length      int    `yaml:"length"`
	Size        int    `yaml:"size"`
	URL         string `yaml:"url"`
	PublishedAt string `yaml:"published_at"`
}

func (v *Video) GenerateURLFromID() {
	url := v.ID
	v.URL = url
}

func (v *Video) MapVideoFromYoutubeResult(yt *youtube.SearchResult) {
	v.ID = yt.Sni
}
