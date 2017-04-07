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
func FetchUser(username string) (*User, error) {
	user := &User{}

	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	_, _, errors := gorequest.New().Get(url).EndStruct(&user)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return user, nil
}

// FetchRepos fetches repositories data for a give username from github
func FetchRepos(username string) ([]*Repository, error) {
	repos := []*Repository{}

	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	_, _, errors := gorequest.New().Get(url).EndStruct(&repos)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return repos, nil
}
