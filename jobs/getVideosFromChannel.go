package jobs

import (
	"errors"
	"github.com/pkj-m/wimc/config"
	"github.com/pkj-m/wimc/mongo"
	"go.uber.org/zap"

	"google.golang.org/api/youtube/v3"
)

// GetVideosFromChannel accepts a YouTube channel link, and downloads list of all videos on the channel
func GetVideosFromChannel(yt *youtube.Service, logger *zap.Logger, cfg *config.AppConfig) ([]*youtube.SearchResult, error) {
	call := yt.Search.List(cfg.Youtube.ChannelListParts).
		ChannelId(cfg.Youtube.ChannelID).Type("video").
		MaxResults(int64(cfg.Youtube.MaxSearchResults)).Order("date")

	resp, err := call.Do()
	if err != nil {
		logger.Error("error in making call to youtube", zap.Error(err))
		return nil, err
	}

	if resp == nil {
		logger.Error("nil response received from youtube")
		return nil, errors.New("nil response received")
	}

	logger.Info("received response from YouTube server")

	// we could extract just the video ID from the result and return a clean string, but in doing so
	// we also lose out on a bunch of extra information (such as the video title, published at, etc)
	// which might later come in handy when we try to enrich the FE with more details about the clip
	var searchResults []*youtube.SearchResult
	for _, item := range resp.Items {
		if item.Id.Kind == "youtube#video" {
			searchResults = append(searchResults, item)
		} else {
			logger.Info("non video response received from server", zap.Any("item", item))
		}
	}

	logger.Info("processed search results", zap.Int("count", len(searchResults)))

	// save the search results in mongo collection
	err = mongo.SaveSearchResults(cfg, logger, searchResults)
	if err != nil {
		return nil, err
	}

	return searchResults, nil
}
