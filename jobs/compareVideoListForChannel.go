package jobs

import "google.golang.org/api/youtube/v3"

type StrippedSearchResult struct {
	Title       string
	PublishedAt string
}

// CompareVideoListForChannel compares the search results found today with yesterday's results
// and returns the list of new videos that were not present in the previous day's search results
func CompareVideoListForChannel(searchResults []*youtube.SearchResult) ([]*youtube.SearchResult, error) {
	// first we fetch the list of results from previous day
	latestResults := extractTitleAndDateFromResults(searchResults)

	// should be implemented in a mongo gateway
	previousResults := fetchPreviousResults()

	newResults := findNewVideos(latestResults, previousResults)

	resp := make([]*youtube.SearchResult, 0)
	for _, item := range newResults {
		resp = append(resp, item)
	}
	return resp, nil
}

func extractTitleAndDateFromResults(results []*youtube.SearchResult) map[string]*youtube.SearchResult {
	var response map[string]*youtube.SearchResult
	for _, result := range results {
		response[result.Snippet.Title] = result
	}
	return response
}

func fetchPreviousResults() map[string]*youtube.SearchResult {
	return nil
}

func findNewVideos(new map[string]*youtube.SearchResult, old map[string]*youtube.SearchResult) map[string]*youtube.SearchResult {
	var newVideos map[string]*youtube.SearchResult
	for id, res := range new {
		// if id of video is not found in yesterday's response, it must be a new video
		if _, ok := old[id]; !ok {
			newVideos[id] = res
		}
	}
	return newVideos
}
