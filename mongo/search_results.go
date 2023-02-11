package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkj-m/wimc/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"google.golang.org/api/youtube/v3"
	"time"
)

// SaveSearchResults accepts a list of YouTube search results and saves it in a mongo collection tagged with today's date
func SaveSearchResults(cfg *config.AppConfig, logger *zap.Logger, results []*youtube.SearchResult) error {
	currentTime := time.Now()
	today := currentTime.Format("01-02-2006")

	collectionName := fmt.Sprintf("%s-%s", cfg.Mongo.CollectionPrefix, today)
	newLogger := logger.With(
		zap.String("callSiteTag", "mongo::SaveSearchResults"),
		zap.String("collection", collectionName),
		zap.Int("input count", len(results)))
	newLogger.Info("called")

	err := Client.Database(cfg.Mongo.Database).CreateCollection(context.TODO(), collectionName)
	if err != nil {
		newLogger.Error("failed to create mongo collection", zap.Error(err))
		return err
	}

	resultsInterface := convertToInterface(newLogger, results)
	res, err := Client.Database(cfg.Mongo.Database).Collection(collectionName).InsertMany(context.TODO(), resultsInterface)
	if err != nil {
		newLogger.Error("failed to save results in collection", zap.Error(err))
		return err
	}
	newLogger.Info("success", zap.Int("saved count", len(res.InsertedIDs)))
	return nil
}

// FetchSearchResults accepts a date and returns the list of search results that were received that day from YouTube
func FetchSearchResults(cfg *config.AppConfig, logger *zap.Logger, target *time.Time) ([]*youtube.SearchResult, error) {
	targetDate := target.Format("01-02-2006")
	collectionName := fmt.Sprintf("%s-%s", cfg.Mongo.CollectionPrefix, targetDate)

	newLogger := logger.With(
		zap.String("callSiteTag", "mongo::FetchSearchResults"),
		zap.String("collection", collectionName),
		zap.String("target time", target.String()))
	newLogger.Info("called")

	if target == nil {
		newLogger.Error("nil target passed")
		return nil, errors.New("invalid argument received")
	}

	coll := Client.Database(cfg.Mongo.Database).Collection(collectionName)
	if coll == nil {
		newLogger.Error("failed to fetch collection")
		return []*youtube.SearchResult{}, nil
	}

	cur, err := Client.Database(cfg.Mongo.Database).Collection(collectionName).Find(context.TODO(), bson.D{})
	if err != nil {
		newLogger.Error("failed to fetch documents")
		return nil, err
	}

	var results []*youtube.SearchResult
	for cur.Next(context.TODO()) {
		var result youtube.SearchResult
		err := cur.Decode(&result)
		if err != nil {
			newLogger.Error("failed to decode response into searchResult", zap.Error(err), zap.Any("result", cur.Current))
			continue
		}
		results = append(results, &result)
	}

	newLogger.Info("success", zap.Int("count", len(results)))
	return results, nil
}

func convertToInterface(logger *zap.Logger, results []*youtube.SearchResult) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, result := range results {
		if result == nil {
			logger.Error("nil result pointer in results list")
			continue
		}
		var inter interface{}
		txt, _ := json.Marshal(*result)
		err := json.Unmarshal(txt, &inter)
		if err != nil {
			logger.Error("failed to unmarshal search result", zap.Any("encoded", string(txt)))
			continue
		}
		interfaces = append(interfaces, inter)
	}
	return interfaces
}
