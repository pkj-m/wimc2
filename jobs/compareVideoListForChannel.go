package jobs

import (
	"github.com/pkj-m/wimc/config"
	"github.com/pkj-m/wimc/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/youtube/v3"
	"time"
)

type StrippedSearchResult struct {
	Title       string
	PublishedAt string
}

// CompareVideoListForChannel compares the search results found today with yesterday's results
// and returns the list of new videos that were not present in the previous day's search results
func CompareVideoListForChannel(cfg *config.AppConfig, logger *zap.Logger, searchResults []*youtube.SearchResult) ([]*youtube.SearchResult, error) {
	// parse today's search results in a map to identify new entries
	todaysResults := indexResultsByVideoID(searchResults)

	yesterday := time.Now().AddDate(0, 0, -1)
	previousResults, err := mongo.FetchSearchResults(cfg, logger, &yesterday)
	if err != nil {
		return nil, err
	}
	yesterdaysResults := indexResultsByVideoID(previousResults)

	newResults := findNewVideos(todaysResults, yesterdaysResults)

	resp := make([]*youtube.SearchResult, 0)
	for _, item := range newResults {
		resp = append(resp, item)
	}
	return resp, nil
}

// indexResultsByVideoID accepts a list of search results and returns a map of same results
// indexed by the video ID of each result
func indexResultsByVideoID(results []*youtube.SearchResult) map[string]*youtube.SearchResult {
	response := make(map[string]*youtube.SearchResult)
	for _, result := range results {
		response[result.Id.VideoId] = result
	}
	return response
}

func findNewVideos(new map[string]*youtube.SearchResult, old map[string]*youtube.SearchResult) map[string]*youtube.SearchResult {
	newVideos := make(map[string]*youtube.SearchResult)
	for id, res := range new {
		// if id of video is not found in yesterday's response, it must be a new video
		if _, ok := old[id]; !ok {
			newVideos[id] = res
		}
	}
	return newVideos
}
