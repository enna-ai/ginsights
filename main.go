package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
)

type User struct {
	Login string `json:"login"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")

	client := &http.Client{}

	followers, err := GetUsers(client, username, "followers", token)
	if err != nil {
		log.Fatal("Error fetching followers:", err)
		return
	}

	following, err := GetUsers(client, username, "following", token)
	if err != nil {
		log.Fatal("Error fetching following list:", err)
		return
	}

	FormatTable(following, followers)
}

func MakeRequest(method, url, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

    req.Header.Add("Accept", "application/vnd.github+json")
    req.Header.Add("Authorization", "Bearer "+token)
    req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
    req.Header.Set("Content-Type", "application/json")
    return req, nil
}

func GetUsers(client *http.Client, username, relation, token string) ([]User, error) {
	url := fmt.Sprintf("%s/users/%s/%s", "https://api.github.com", username, relation)
	req, err := MakeRequest("GET", url, token)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func FormatTable(followings, followers []User) {
	followerStatus := make(map[string]bool)
	for _, follower := range followers {
		followerStatus[follower.Login] = true
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Username", "Following Me?"})

	for _, following := range followings {
		isFollower := followerStatus[following.Login]

		followerStatusStr := fmt.Sprintf("%v", isFollower)

		table.Append([]string{following.Login, followerStatusStr})
	}

	table.Render()
}
