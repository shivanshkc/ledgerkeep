package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shivanshkc/ledgerkeep/src/configs"
	"github.com/shivanshkc/ledgerkeep/src/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	clientOnce      = &sync.Once{}
	clientSingleton *mongo.Client
)

// GetClient returns the MongoDB client singleton.
func GetClient() *mongo.Client {
	clientOnce.Do(func() {
		logger.Get().Info(context.Background(), &logger.Entry{Payload: "Attempting connection with MongoDB..."})
		clientSingleton = getClient()
		logger.Get().Info(context.Background(), &logger.Entry{Payload: "Connected with MongoDB."})
	})

	return clientSingleton
}

// getClient is a pure function (except for configs) to generate a MongoDB client.
func getClient() *mongo.Client {
	conf := configs.Get()
	connectOpts := options.Client().ApplyURI(conf.Mongo.Addr)

	client, err := mongo.NewClient(connectOpts)
	if err != nil {
		panic(fmt.Errorf("failed to create mongodb client: %w", err))
	}

	timeoutDuration := time.Duration(conf.Mongo.OperationTimeoutSec) * time.Second
	connectCtx, connectCancelFunc := context.WithTimeout(context.Background(), timeoutDuration)
	defer connectCancelFunc()

	if err := client.Connect(connectCtx); err != nil {
		panic(fmt.Errorf("failed to connect to mongodb: %w", err))
	}

	pingCtx, pingCancelFunc := context.WithTimeout(context.Background(), timeoutDuration)
	defer pingCancelFunc()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		panic(fmt.Errorf("failed to ping mongodb: %w", err))
	}

	return client
}
