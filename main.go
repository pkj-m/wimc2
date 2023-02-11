package main

import (
	"context"
	"fmt"
	"github.com/pkj-m/wimc/config"
	"github.com/pkj-m/wimc/mongo"

	"github.com/pkj-m/wimc/jobs"

	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
)

func main() {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	logger.Info("loaded config successfully...")

	mongo.Setup(cfg, logger)
	logger.Info("connected to mongo...")

	yt, err := youtube.NewService(ctx, option.WithAPIKey(cfg.Youtube.ApiKey))
	if err != nil {
		logger.Fatal("failed to instantiate youtube client", zap.Error(err))
	}

	searchResults, err := jobs.GetVideosFromChannel(yt, logger, cfg)
	if err != nil {
		logger.Fatal("failed to fetch videos from youtube", zap.Error(err))
	}
	logger.Info("fetched videos from youtube", zap.Any("count", len(searchResults)))

	newResults, err := jobs.CompareVideoListForChannel(cfg, logger, searchResults)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("new videos: ")
	for _, result := range newResults {
		fmt.Printf("[%v] %v\n", result.Id.VideoId, result.Snippet.Title)
		//video, err := jobs.DownloadVideoFromID(result)
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
