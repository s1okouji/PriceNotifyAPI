package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/s1okouji/price_notify_api/model"
)

type appdetails struct {
	Data struct {
		Name           string `json:"name"`
		Price_overview struct {
			Initial          int `json:"initial"`
			Final            int `json:"final"`
			Discount_percent int `json:"discount_percent"`
		} `json:"price_overview"`
	} `json:"data"`
}

func GetGameModel(appId int) model.Game {
	resp, err := http.Get("https://store.steampowered.com/api/appdetails?currency=JP&appids=" + strconv.Itoa(appId))
	if err != nil {
		log.Fatal("cannot get price")
	}

	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	var jsonData appdetails
	// fmt.Println(string(byteArray))
	// fmt.Println(GetData(string(byteArray)))
	json.Unmarshal([]byte(GetData(string(byteArray))), &jsonData)

	result := model.Game{
		AppId:           appId,
		AppName:         jsonData.Data.Name,
		BasePrice:       jsonData.Data.Price_overview.Initial,
		FinalPrice:      jsonData.Data.Price_overview.Final,
		DiscountPercent: jsonData.Data.Price_overview.Discount_percent,
		CreatedAt:       time.Now(),
	}

	return result
}
