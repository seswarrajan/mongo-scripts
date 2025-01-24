package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	credentials := options.Credential{
		Username: "",
		Password: "",
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetAuth(credentials)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Define database and collection
	dbName := "data-pipeline"
	collectionName := "kubearmor_alerts_529"
	collection := client.Database(dbName).Collection(collectionName)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Query the collection (fetch all documents)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error querying collection:", err)
		return
	}
	defer cursor.Close(ctx)

	// Prepare file for writing
	outputFile := "kubearmor_alerts_529.json"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write documents to file in JSON format
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print

	var count int
	for cursor.Next(ctx) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			fmt.Println("Error decoding document:", err)
			continue
		}
		if err := encoder.Encode(document); err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}
		count++
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return
	}

	fmt.Printf("Dumped %d documents from collection '%s' to file '%s'\n", count, collectionName, outputFile)
}
