package jobs

import (
	"errors"

	"google.golang.org/api/youtube/v3"
)

// GetVideosFromChannel accepts a youtube channel link, and downloads list of all videos on the channel
func GetVideosFromChannel(yt *youtube.Service, channelID string) ([]*youtube.SearchResult, error) {
	call := yt.Search.List([]string{"snippet"}).ChannelId(channelID).Type("video").MaxResults(20).Order("date")
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("nil response received")
	}

	// we could extract just the video ID from the result and return a clean string, but in doing so
	// we also lose out on a bunch of extra information (such as the video title, published at, etc)
	// which might later come in handy when we try to enrich the FE with more details about the clip
	var videoIDs []*youtube.SearchResult
	for _, item := range resp.Items {
		if item.Id.Kind == "youtube#video" {
			videoIDs = append(videoIDs, item)
		}
	}

	return videoIDs, nil
}
