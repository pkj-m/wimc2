package main

import (
	"context"

	"github.com/pkj-m/wimc/jobs"

	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	//config.Parse(config{}, configFilePath)

	ctx := context.Background()
	youtube, err := youtube.NewService(ctx, option.WithAPIKey("my-secret-key"))
	if err != nil {
		log.Fatal("failed to instantiate youtube client")
	}

	channelID := "UCpUS8oR-IJSP3gFYrGAWEuQ"

	searchResults, err := jobs.GetVideosFromChannel(youtube, channelID)
	if err != nil {
		log.Fatal(err.Error())
	}

	// compare latest list of videoIDs with previous day to identify how many new videos need to be fetched
	newResults, err := jobs.CompareVideoListForChannel(searchResults)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, result := range newResults {
		video, err := jobs.DownloadVideoFromID(result)
	}

	// for _, videoID := range newVideos {
	// 	video, err := jobs.DownloadVideoFromID(videoID)
	// 	if err != nil {
	// 		log.Println("failed to download video from youtube: ", err.Error())
	// 		continue
	// 	}

	// 	// video --> audio  my cannon--> text -->

	// 	audio, err := jobs.ExtractAudioFromVideo(video)
	// 	if err != nil {
	// 		log.Println("failed to extract audio from video: ", err.Error())
	// 		continue
	// 	}

	// 	// do we accept multiple timestamps ??
	// 	eventTimestamp, err := jobs.FindTimestamp(audio)
	// 	if err != nil {
	// 		log.Println("error in finding event timestamp: ", err.Error())
	// 		continue
	// 	}

	// 	err = jobs.TrimVideo(eventTimestamp)
	// 	if err != nil {
	// 		log.Println("error in trimming video: ", err.Error())
	// 		continue
	// 	}
	// }

}
