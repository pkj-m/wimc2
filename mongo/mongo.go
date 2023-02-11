package mongo

import (
	"context"
	"fmt"
	"github.com/pkj-m/wimc/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Client *mongo.Client

func Setup(cfg *config.AppConfig, logger *zap.Logger) {
	connectionURI := getConnectionURI(cfg)
	var err error
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionURI))
	if err != nil {
		logger.Fatal("could not connect to mongo", zap.Error(err))
	}
}

func CloseConnection() error {
	if err := Client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}

func getConnectionURI(cfg *config.AppConfig) string {
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.1llsgcv.mongodb.net/test", cfg.Mongo.Username, cfg.Mongo.Password)
	return uri
}
