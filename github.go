package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GithubResponse struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

var avatarMemo = NewMemo(func(githubName string) (interface{}, error) {
	log.Println("requesting avatar for:", githubName)
	return getAvatar(githubName)
})

func getAvatar(githubName string) ([]byte, error) {
	githubUserInfo, err := getGithubUserInfo(githubName)
	if err != nil {
		return nil, err
	}

	avatar, err := getGithubUserAvatar(githubUserInfo.AvatarUrl)
	if err != nil {
		return nil, err
	}

	return avatar, nil
}

func getGithubUserInfo(githubName string) (*GithubResponse, error) {
	userResp, err := http.Get("https://api.github.com/users/" + githubName)
	if err != nil || userResp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("github name: %s, fetch user info err: %v", githubName, err)
	}
	defer userResp.Body.Close()

	var githubResponse GithubResponse
	respData, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		return nil, fmt.Errorf("github name: %s, read user api response body err: %v", githubName, err)
	}

	err = json.Unmarshal(respData, &githubResponse)
	if err != nil {
		return nil, fmt.Errorf("github name: %s, parse user api response body err: %v", githubName, err)
	}

	return &githubResponse, nil
}

func getGithubUserAvatar(url string) ([]byte, error) {
	avatarResp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("avatar url: '%s', request github avatar api err: %v", url, err)
	}
	defer avatarResp.Body.Close()

	data, err := ioutil.ReadAll(avatarResp.Body)
	if err != nil {
		return nil, fmt.Errorf("avatar url: '%s', read github avatar api response body err: %v", url, err)
	}

	return data, nil
}
