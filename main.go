package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// connect to mongo
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Could not load .env file.")
		os.Exit(1)
	}

	password := os.Getenv("MONGO_PUCKOS_PASS")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://puckos:" + password + "@puckosdb.wqf3k.gcp.mongodb.net/myFirstDatabase?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// check json data for new sets, save new sets
	sets := get_lego_item_data_from_mongo(client, ctx)
	payload := read_json_data()
	idsFromMongo := get_set_ids_from_mongo_sets(sets)
	check_if_item_in_json_is_missing_from_mongo(idsFromMongo, payload, client, ctx)

	docs := scrape_links_from_mongo_for_price_data_(client, ctx, sets)

	if len(docs) > 0 {
		interface_docs := convert_scraped_price_data_to_interface(docs)
		insert_newly_scraped_price_data(client, ctx, interface_docs)
	}


}
