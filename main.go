package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucasefe/weader/gh"
	"github.com/lucasefe/weader/util"
	"github.com/lucasefe/weader/weather"
)

// TODO: Tests, please.
// TODO: Handle multiple errors returned from gorequest
// TODO: Add http caching layer
// TODO: Add weather cache.
// TODO: Research if weather api supports batch
// TODO: Error handling at http level is very repetitive. DRY it up.
// TODO: Paginate the repos.

// Result represents the server unique result, for now.
type Result struct {
	AvgTemperature int    `json:"avg_temperature"`
	Location       string `json:"location"`
	ReposCount     int    `json:"repos_count"`
	Username       string `json:"username"`
}

func main() {
	router := httprouter.New()
	router.GET("/:username", getByUsername)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", util.NewTimer(router))
}

func getByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("username")

	user, err := gh.FetchUser(username)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repos, err := gh.FetchRepos(user.Login)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	temperatures, err := fetchTemperatures(user, repos)
	if err != nil {
		log.Printf("error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := Result{
		Username:       user.Login,
		Location:       user.Location,
		ReposCount:     len(repos),
		AvgTemperature: avg(temperatures),
	}

	fmt.Fprint(w, render(res))
}

func render(res Result) string {
	b, _ := json.Marshal(res)
	s := string(b)

	return s
}

func avg(numbers []int) int {
	count := len(numbers)

	if count == 0 {
		return 0
	}

	var sum int
	for _, t := range numbers {
		sum += t
	}

	return sum / count
}

func fetchTemperatures(user *gh.User, repos []*gh.Repository) ([]int, error) {
	temperatures := []int{}

	for _, repo := range repos {
		temperature, err := weather.FetchTemperature(user.Location, repo.CreatedAt)

		if err != nil {
			return nil, err
		}

		temperatures = append(temperatures, temperature)
	}

	return temperatures, nil
}
