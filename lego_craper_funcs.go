package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func parse_price_string(price_string string) float64 {
	scraped_price := price_string[:len(price_string)-3]
	scraped_price = strings.Replace(scraped_price, " ", "", -1)
	int_price, err := strconv.ParseInt(scraped_price, 0, 0)
	float_price := float64(int_price)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	return float_price
}

func get_lego_prices(link string) []float64 {
	c := colly.NewCollector()

	price_list := []float64{}

	c.OnHTML("div.row-price", func(h *colly.HTMLElement) {
		scraped_price := h.Text
		int_price := parse_price_string(scraped_price)
		price_list = append(price_list, int_price)
	})


	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping %s\n", e.Error())
	})

	c.Visit(link)

	sort.Float64s(price_list)

	return price_list

}

func calculate_average_of_first_n_prices(n int, price_list []float64) float64 {
	shortened_slice := price_list[:n]

	var sum float64 = 0
	for _, price := range shortened_slice {
		sum += price
	}

	return sum / float64(n)
}

func calculate_min_price_ratio(num_of_prices_to_check int, price_list []float64) float64 {
	min_price := price_list[0]
	average_price := calculate_average_of_first_n_prices(num_of_prices_to_check, price_list)

	return min_price / average_price
}

func check_if_price_ratio_is_interesting(actual_price_ratio float64, treshold float64) bool {
	return actual_price_ratio < treshold
}

func manage_find_lego_price_data_query_result(query_result []bson.D, docs *[]bson.D, data bson.D){
	
	if len(query_result) == 0 {
		fmt.Printf("query_result returned empty, preparing document to insert\n")
		*docs = append(*docs, data)
	}
}

func scrape_links_from_mongo_for_price_data_(client *mongo.Client, ctx context.Context, setsFromMongo []bson.M) []bson.D{

	var docs []bson.D

	for i := 0; i < len(setsFromMongo); i++ {
		var link interface{} = setsFromMongo[i]["Link"]
		set_link_from_mongo := link.(string)

		var set interface{} = setsFromMongo[i]["SetId"]
		SetId := set.(int32)

		var _id interface{} = setsFromMongo[i]["_id"]
		set_data_id := _id.(primitive.ObjectID)

		price_list := get_lego_prices(set_link_from_mongo)

		price_ratio := calculate_min_price_ratio(4, price_list)

		price_ratio_is_interesting := check_if_price_ratio_is_interesting(price_ratio, 0.95)
		fmt.Println("price_ratio_is_interesting: ", price_ratio_is_interesting)

		today := time.Now()
		date_of_scraping := today.Format("2006-01-02")

		query := bson.D{{"SetId", SetId}, {"date_of_scraping", date_of_scraping}}

		query_result := find_lego_price_data_by_id_and_date(client, ctx, query)

		data := bson.D{
			{"SetId", SetId}, 
			{"set_data_id", set_data_id}, 
			{"price_list", price_list}, 
			{"price_ratio", price_ratio}, 
			{"price_ratio_is_interesting", price_ratio_is_interesting}, 
			{"date_of_scraping", date_of_scraping},
		}

		if len(query_result) == 0 {
			fmt.Printf("query_result returned empty, preparing document to insert\n")
			docs = append(docs, data)
		} else {
			update_existing_daily_price_entry(client, ctx, SetId, price_list, price_ratio_is_interesting)
		}
	}

	return docs
}
