package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type JsonLegoSets struct {
	JsonSetsSlice []LegoData `json:"sets"`
}

type LegoData struct {
	SetId   int32  `json:"set_id"`
	SetName string `json:"set_name"`
	Link    string `json:"link"`
}

func read_json_data() JsonLegoSets {
	jsonFile, err := os.Open("lego_sets.json")
	if err != nil {
		fmt.Println("Error when opening file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var payload JsonLegoSets
	err = json.Unmarshal(byteValue, &payload)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
	}

	return payload
}

func parse_lego_json_data(legoData LegoData) (int32, string, string) {
	return legoData.SetId, legoData.SetName, legoData.Link
}










