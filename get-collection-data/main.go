package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

func main() {
	username := flag.String("user", "", "mongodb username")
	password := flag.String("password", "", "mongodb password")
	url := flag.String("host", "mongodb://127.0.0.1:27017/?directConnection=true", "mongodb host")
	database := flag.String("db", "test", "mongodb database name")

	flag.Parse()

	connect(*username, *password, *url, *database)

	// Get a list of collections in the database
	collections, err := DB.ListCollectionNames(context.Background(), mongo.Pipeline{})
	if err != nil {
		fmt.Printf("Failed to list collections: %v", err)
	}

	// Traverse all collections and get the count
	for _, collectionName := range collections {
		collection := DB.Collection(collectionName)

		// Count the documents in the collection
		count, err := collection.CountDocuments(context.Background(), mongo.Pipeline{})
		if err != nil {
			log.Printf("Failed to count documents in collection %s: %v", collectionName, err)
			continue
		}

		fmt.Printf("Collection: %s, Document Count: %d\n", collectionName, count)
	}

}

func connect(username string, password string, url string, database string) {
	credentials := options.Credential{
		Username: username,
		Password: password,
	}

	clientOptions := options.Client().ApplyURI(url)
	if username != "" && password != "" {
		clientOptions.SetAuth(credentials)
	}

	log.Info().Msgf("Connecting to MongoDB - %s", url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal().Msgf("Failed to connect. %v", err)
	}

	log.Info().Msg("Pinging DB...")
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal().Msgf("Failed to connect. %v", err)
	}

	log.Info().Msg("Connected Successfully")
	DB = client.Database(database)
}
