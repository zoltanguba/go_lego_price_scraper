package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Parse bson file?
func get_set_ids_from_mongo_sets(setsFromMongo []bson.M) []int32 {
	var idsFromMongo []int32
	for i := 0; i < len(setsFromMongo); i++ {
		var i interface{} = setsFromMongo[i]["SetId"]
		t := i.(int32)
		idsFromMongo = append(idsFromMongo, t)
	}
	return idsFromMongo
}

// Parse bson file?
func check_if_item_in_json_is_missing_from_mongo(idsFromMongo []int32, setsFromJson JsonLegoSets, client *mongo.Client, ctx context.Context) {
	for i := 0; i < len(setsFromJson.JsonSetsSlice); i++ {
		SetId, SetName, Link := parse_lego_json_data(setsFromJson.JsonSetsSlice[i])

		// add new item if any
		if contains(idsFromMongo, SetId) {
			fmt.Println(SetId, " is on Mongo")
		} else {
			data := bson.D{
				{Key: "SetId", Value: SetId},
				{Key: "SetName", Value: SetName},
				{Key: "Link", Value: Link},
			}
			insert_lego_document(client, ctx, "legoItemData", data)
		}
	}
}
