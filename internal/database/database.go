package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get DB connection instance
func GetDBConn(ctx context.Context, uri string, dbName string) (*mongo.Database, error) {
	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	dbInstance := client.Database(dbName)

	return dbInstance, nil
}

// Get default DB connection instance
func GetDBConnDefault(ctx context.Context) (*mongo.Database, error) {
	return GetDBConn(ctx, getDbUri(), getDbName())
}

// Get databse URI
func getDbUri() string {
	protocol := "mongodb"
	if envProtocol, ok := os.LookupEnv("DB_PROTOCOL"); ok {
		protocol = envProtocol
	}
	host := "localhost"
	if envHost, ok := os.LookupEnv("DB_HOST"); ok {
		host = envHost
	}
	port := "27017"
	if envPort, ok := os.LookupEnv("DB_PORT"); ok {
		port = envPort
	}
	user := "root"
	if envUser, ok := os.LookupEnv("DB_USER"); ok {
		user = envUser
	}
	password := ""
	if envPassword, ok := os.LookupEnv("DB_PASS"); ok {
		password = envPassword
	}
	return protocol + "://" + user + ":" + password + "@" + host + ":" + port
}

// Get database name
func getDbName() string {
	dbName := "goblog"
	if envDbName, ok := os.LookupEnv("DB_NAME"); ok {
		dbName = envDbName
	}
	return dbName
}
