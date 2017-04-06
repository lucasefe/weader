package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type repository struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type result struct {
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

	location := fetchLocation(username)
	repos := fetchRepos(username)

	temperatures := []int{}
	for _, repo := range repos {
		temperature := fetchTemperature(location, repo.CreatedAt)
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

	res := result{
		AvgTemperature: avgTemperature,
		Username:       username,
		Location:       location,
		ReposCount:     len(repos),
	}

	fmt.Fprint(w, render(res))
}

func render(res result) string {
	b, _ := json.Marshal(res)
	s := string(b)
	return s
}

func fetchLocation(username string) string {
	return "Buenos Aires, Argentina"
}

func fetchRepos(username string) []repository {
	repos := []repository{}

	return repos
}

func fetchTemperature(location string, date time.Time) int {
	return 12
}
