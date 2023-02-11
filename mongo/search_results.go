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
	today := currentTime.Format("09-07-2017")

	collectionName := fmt.Sprintf("%s-%s", cfg.Mongo.CollectionPrefix, today)
	err := Client.Database(cfg.Mongo.Database).CreateCollection(context.TODO(), collectionName)
	if err != nil {
		logger.Error("failed to create mongo collection", zap.Error(err), zap.String("collection", collectionName))
		return err
	}

	resultsInterface := convertToInterface(results)
	res, err := Client.Database(cfg.Mongo.Database).Collection(collectionName).InsertMany(context.TODO(), resultsInterface)
	if err != nil {
		logger.Error("failed to save results in collection", zap.Error(err))
		return err
	}
	logger.Info("successfully saved search results in mongo collection", zap.Int("count", len(res.InsertedIDs)))
	return nil
}

// FetchSearchResults accepts a date and returns the list of search results that were received that day from YouTube
func FetchSearchResults(cfg *config.AppConfig, logger *zap.Logger, target *time.Time) ([]*youtube.SearchResult, error) {
	if target == nil {
		logger.Error("nil target passed")
		return nil, errors.New("invalid argument received")
	}

	targetDate := target.Format("09-07-2017")
	collectionName := fmt.Sprintf("%s-%s", cfg.Mongo.CollectionPrefix, targetDate)

	coll := Client.Database(cfg.Mongo.Database).Collection(collectionName)
	if coll == nil {
		logger.Error("could not find mongo collection", zap.String("collection", collectionName))
		return []*youtube.SearchResult{}, nil
	}

	cur, err := Client.Database(cfg.Mongo.Database).Collection(collectionName).Find(context.TODO(), bson.D{})
	if err != nil {
		logger.Error("error while fetching documents from collection", zap.String("collection", collectionName))
		return nil, err
	}

	var results []*youtube.SearchResult
	for cur.Next(context.TODO()) {
		var result youtube.SearchResult
		err := cur.Decode(&result)
		if err != nil {
			logger.Error("failed to decode mongo response into searchResult", zap.Error(err), zap.Any("result", cur.Current))
			continue
		}
		results = append(results, &result)
	}

	return results, nil
}

func convertToInterface(results []*youtube.SearchResult) []interface{} {
	interfaces := make([]interface{}, 0)
	for _, result := range results {
		if result == nil {
			continue
		}
		var inter interface{}
		txt, _ := json.Marshal(*result)
		err := json.Unmarshal(txt, inter)
		if err != nil {
			// handle error
			continue
		}
		interfaces = append(interfaces, inter)
	}
	return interfaces
}
