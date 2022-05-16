package main

import (
	"go.mongodb.org/mongo-driver/bson"
)

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


func convert_scraped_price_data_to_interface(scraped_price_data []bson.D) []interface{} {
	price_data_interface := make([]interface{}, len(scraped_price_data))
	for i := range scraped_price_data {
		price_data_interface[i] = scraped_price_data[i]
	}

	return price_data_interface
}