package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
)

var weatherAPIKey string

func init() {
	weatherAPIKey = getWeatherAPIKey()
}

// RootResult represents the root weather data response
type RootResult struct {
	Data struct {
		Weather []struct {
			MaxTempC string `json:"maxtempC"`
			MaxTempF string `json:"maxtempF"`
		} `json:"weather"`
	} `json:"data"`
}

// FetchTemperature retrieves temperature from api for a given location and date
func FetchTemperature(location string, date time.Time) (int, error) {
	dateString := date.Format("2006-01-02")

	url := fmt.Sprintf("http://api.worldweatheronline.com/premium/v1/past-weather.ashx?key=%s&format=json&q=%s&date=%s",
		weatherAPIKey,
		location,
		dateString)

	_, body, errs := gorequest.New().Get(url).End()
	if len(errs) > 0 {
		return 0, fmt.Errorf("Could not get weather information for %s: %+v", location, errs[0])
	}

	var root RootResult
	err := json.Unmarshal([]byte(body), &root)
	if err != nil {
		return 0, err
	}

	if len(root.Data.Weather) > 0 {
		temp, err := strconv.Atoi(root.Data.Weather[0].MaxTempC)

		if err != nil {
			return 0, err
		}

		return temp, nil
	}

	return 0, nil
}

func getWeatherAPIKey() string {
	key, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		log.Fatal("Provide WEATHER_API_KEY!")
		os.Exit(2)
	}

	return key
}
