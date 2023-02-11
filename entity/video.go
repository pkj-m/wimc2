package entity

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"url"`
	Length      int    `yaml:"length"` // length of the video in secs
	Size        int    `yaml:"size"`   // size of the video in bytes
	URL         string `yaml:"url"`
	PublishedAt string `yaml:"published_at"` // ISO8601 format time of publishing video
	Description string `yaml:"description"`  // description of the youtube video
	S3Link      string `yaml:"s3_link"`      // link of s3 bucket where video is present
}

func (v *Video) GenerateURLFromID() string {
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.ID)
	return url
}

func (v *Video) MapVideoFromYoutubeResult(yt *youtube.SearchResult) {
	v.ID = yt.Id.VideoId
	v.Title = yt.Snippet.Title
	v.Description = yt.Snippet.Description
	v.PublishedAt = yt.Snippet.PublishedAt
	v.URL = v.GenerateURLFromID()
}
