package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucasefe/weader/gh"
	"github.com/lucasefe/weader/weather"
)

// TODO: Handle errors from the API
// TODO: Handle multiple errors returned from gorequest
// TODO: Add http caching layer
// TODO: Add weather cache.
// TODO: Research if weather api supports batch

// Result represents the server unique result, for now.
type Result struct {
	AvgTemperature int    `json:"avg_temperature"`
	Location       string `json:"location"`
	ReposCount     int    `json:"repos_count"`
	Username       string `json:"username"`
}

func main() {

	router := httprouter.New()
	router.GET("/:username", byUsername)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", router)
}

func byUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")

	user := gh.FetchUser(username)
	repos := gh.FetchRepos(user.Login)

	temperatures := []int{}
	for _, repo := range repos {
		temperature, err := weather.FetchTemperature(user.Location, repo.CreatedAt)

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
