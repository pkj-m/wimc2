package jobs

import (
	"github.com/pkj-m/wimc/entity"
	"google.golang.org/api/youtube/v3"
)

// DownloadVideoFromURL accepts a video URL, downloads the video and saves it to bucket
func DownloadVideoFromID(ytVideo *youtube.SearchResult) (*entity.Video, error) {
	video := &entity.Video{}
	video.MapVideoFromYoutubeResult(ytVideo)

	// how do we download a youtube video using APIs?

}
