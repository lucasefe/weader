package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/parnurzeal/gorequest"
)

// TODO: Handle errors from the API
// TODO: Handle multiple errors returned from gorequest
// TODO: Big file. Split.
// TODO: Add http caching layer
// TODO: Add weather cache.
// TODO: Research if weather api supports batch

// User represents the github user
type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Location string `json:"location"`
}

// Repository represents the github repository
type Repository struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Result represents the server unique result, for now.
type Result struct {
	AvgTemperature int    `json:"avg_temperature"`
	Location       string `json:"location"`
	ReposCount     int    `json:"repos_count"`
	Username       string `json:"username"`
}

var weatherAPIKey string

// WeatherRoot represents the root weather data response
type WeatherRoot struct {
	Data struct {
		Weather []struct {
			MaxTempC string `json:"maxtempC"`
			MaxTempF string `json:"maxtempF"`
		} `json:"weather"`
	} `json:"data"`
}

func main() {
	weatherAPIKey = getWeatherAPIKey()

	router := httprouter.New()
	router.GET("/:username", byUsername)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", router)
}

func byUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	user := fetchUser(username)
	repos := fetchRepos(user.Login)

	temperatures := []int{}
	for _, repo := range repos {
		temperature, err := fetchTemperature(user.Location, repo.CreatedAt)

		if err != nil {
			log.Printf("error: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		temperatures = append(temperatures, temperature)
	}

	var sum int
	for _, t := range temperatures {
		sum += t
	}

	var avgTemperature int
	if len(temperatures) > 0 {
		avgTemperature = sum / len(temperatures)
	}

	res := Result{
		AvgTemperature: avgTemperature,
		Username:       user.Login,
		Location:       user.Location,
		ReposCount:     len(repos),
	}

	fmt.Fprint(w, render(res))
}

func render(res Result) string {
	b, _ := json.Marshal(res)
	s := string(b)

	return s
}

func fetchUser(username string) *User {
	user := &User{}
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	gorequest.New().Get(url).EndStruct(&user)

	return user
}

func fetchRepos(username string) []Repository {
	repos := []Repository{}
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	gorequest.New().Get(url).EndStruct(&repos)

	return repos
}

func fetchTemperature(location string, date time.Time) (int, error) {
	dateString := date.Format("2006-01-02")

	url := fmt.Sprintf("http://api.worldweatheronline.com/premium/v1/past-weather.ashx?key=%s&format=json&q=%s&date=%s",
		weatherAPIKey,
		location,
		dateString)

	_, body, errs := gorequest.New().Get(url).End()
	if len(errs) > 0 {
		return 0, fmt.Errorf("Could not get weather information for %s: %+v", location, errs[0])
	}

	var root WeatherRoot
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
