package jobs

import (
	"errors"
	"fmt"
	"github.com/pkj-m/wimc/config"
	"go.uber.org/zap"
	"google.golang.org/api/youtube/v3"
	"io"
	"net/http"
)

const (
	englishSubtitleCode = "en"
)

// GetSubtitlesForVideo accepts a video ID and returns the english subtitles if found
func GetSubtitlesForVideo(yt *youtube.Service, logger *zap.Logger, cfg *config.AppConfig, videoID string) (*string, error) {
	newLogger := logger.With(
		zap.String("callSiteTag", "jobs::GetSubtitlesForVideo"),
		zap.String("videoID", videoID))
	newLogger.Info("called")

	call := yt.Captions.List(cfg.Youtube.CaptionsPart, videoID)
	// costs 50 quota points
	resp, err := call.Do()

	if err != nil {
		newLogger.Error("server error", zap.Error(err))
		return nil, err
	}

	if resp == nil {
		newLogger.Error("nil response")
		return nil, errors.New("nil response from server")
	}

	newLogger.Info("received response", zap.Int("count", len(resp.Items)))
	// iterate over response to extract english subs
	var englishSubs *youtube.Caption
	for _, res := range resp.Items {
		logger.Info(fmt.Sprintf("caption response: %v", res))
		if res == nil {
			continue
		}
		englishSubs = res
		break
	}

	if englishSubs == nil {
		newLogger.Error("no english captions found")
		return nil, errors.New("no english captions found")
	}

	return DownloadSubtitlesFromID(yt, logger, cfg, englishSubs)
}

func DownloadSubtitlesFromID(yt *youtube.Service, logger *zap.Logger, cfg *config.AppConfig, cap *youtube.Caption) (*string, error) {
	newLogger := logger.With(
		zap.String("callSiteTag", "jobs::DownloadSubtitlesFromID"),
		zap.String("captionID", cap.Id))
	newLogger.Info("called")

	if cap == nil {
		newLogger.Error("nil caption argument")
		return nil, errors.New("nil caption argument")
	}

	call := yt.Captions.Download(cap.Id)
	resp, err := call.Download()
	if err != nil {
		newLogger.Error("error in downloading caption", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var subs *string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			//log.Fatal(err)
			newLogger.Error("error in reading response", zap.Error(err))
		}
		bodyString := string(bodyBytes)
		subs = &bodyString
	}

	return subs, nil
}
