package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get DB connection instance
func GetDBConn(
	ctx context.Context,
	uri string,
	dbName string,
) (dbConn *mongo.Database, err error) {
	var (
		options    = options.Client().ApplyURI(uri)
		client     *mongo.Client
		dbInstance *mongo.Database
	)

	if client, err = mongo.Connect(ctx, options); err != nil {
		return nil, err
	}
	dbInstance = client.Database(dbName)

	return dbInstance, nil
}

// Get default DB connection instance
func GetDBConnDefault(ctx context.Context) (dbConn *mongo.Database, err error) {
	return GetDBConn(ctx, getDbUri(), getDbName())
}

// Get databse URI
func getDbUri() (dbUri string) {
	var (
		protocol    = "mongodb"
		host        = "localhost"
		port        = "27017"
		user        = "root"
		password    = ""
		envProtocol string
		envHost     string
		envPort     string
		envUser     string
		envPassword string
		ok          bool
	)

	if envProtocol, ok = os.LookupEnv("DB_PROTOCOL"); ok {
		protocol = envProtocol
	}
	if envHost, ok = os.LookupEnv("DB_HOST"); ok {
		host = envHost
	}
	if envPort, ok = os.LookupEnv("DB_PORT"); ok {
		port = envPort
	}
	if envUser, ok = os.LookupEnv("DB_USER"); ok {
		user = envUser
	}
	if envPassword, ok = os.LookupEnv("DB_PASS"); ok {
		password = envPassword
	}

	return protocol + "://" +
		user + ":" +
		password + "@" +
		host + ":" +
		port
}

// Get database name
func getDbName() (dbName string) {
	dbName = "goblog"
	if envDbName, ok := os.LookupEnv("DB_NAME"); ok {
		dbName = envDbName
	}

	return dbName
}
