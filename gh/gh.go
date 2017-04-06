package gh

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

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

// FetchUser fetches user data from github
func FetchUser(username string) *User {
	user := &User{}
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	gorequest.New().Get(url).EndStruct(&user)

	return user
}

// FetchRepos fetches repositories data for a give username from github
func FetchRepos(username string) []Repository {
	repos := []Repository{}
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	gorequest.New().Get(url).EndStruct(&repos)

	return repos
}
