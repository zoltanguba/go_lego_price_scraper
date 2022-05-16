package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// query all lego idem data in the collection
func get_lego_item_data_from_mongo(client *mongo.Client, ctx context.Context) []bson.M {

	goLegoDatabase := client.Database("go_lego")
	goLegoCollection := goLegoDatabase.Collection("legoItemData")

	cursor, err := goLegoCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	var sets []bson.M
	if err = cursor.All(ctx, &sets); err != nil {
		fmt.Printf("Error: %s", err)
	}

	return sets
}

// add document to selected collection
func insert_lego_document(client *mongo.Client, ctx context.Context, collection string, data bson.D) {
	goLegoDatabase := client.Database("go_lego")
	legoCollection := goLegoDatabase.Collection(collection)

	legoScrapeResult, err := legoCollection.InsertOne(ctx, data)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(legoScrapeResult.InsertedID)
}

func find_lego_price_data_by_id_and_date(client *mongo.Client, ctx context.Context, query bson.D) []bson.D {
	goLegoDatabase := client.Database("go_lego")
	legoCollection := goLegoDatabase.Collection("priceScraping")

	cursor, err := legoCollection.Find(ctx, query)

	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results
}

func insert_newly_scraped_price_data(client *mongo.Client, ctx context.Context, data_to_insert []interface{}) {
	goLegoDatabase := client.Database("go_lego")
	legoCollection := goLegoDatabase.Collection("priceScraping")

	result, err := legoCollection.InsertMany(ctx, data_to_insert)
	if err != nil {
		panic(err)
	}
	fmt.Println(result, " has been inserted")
}

func update_existing_daily_price_entry(client *mongo.Client, ctx context.Context, SetId int32, price_list []float64, price_ratio_is_interesting bool) {
	goLegoDatabase := client.Database("go_lego")
	legoCollection := goLegoDatabase.Collection("priceScraping")

	result, err := legoCollection.UpdateOne(
		ctx,
		bson.M{"SetId": SetId},
		bson.D{
			{"$set", bson.D{{"price_list", price_list}}},
			{"$set", bson.D{{"price_ratio_is_interesting", price_ratio_is_interesting}}},
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

/*
func connect_to_mongo1() {

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


	sets := get_lego_item_data_from_mongo(client, ctx)
	payload := read_json_data()
	idsFromMongo := get_set_ids_from_mongo_sets(sets)
	check_if_item_in_json_is_missing_from_mongo(idsFromMongo, payload, client, ctx)

}
*/
