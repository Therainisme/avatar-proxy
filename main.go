package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type GithubResponse struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

func HandleProxyAvatar(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/\S*.png`)
	matchArr := re.FindStringSubmatch(r.URL.Path)
	if len(matchArr) != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "url format err: example <https://avatar.therainisme.com/Therainisme.png>\n")
		return
	}

	githubName := matchArr[0][1 : len(matchArr[0])-4]

	userResp, err := http.Get("https://api.github.com/users/" + githubName)
	if err != nil || userResp.StatusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("github name: %s, fetch user info err: %v\n", githubName, err)
		fmt.Fprintf(w, "fetch user info err: %v\n", err)
		return
	}
	defer userResp.Body.Close()

	var githubResponse GithubResponse
	respData, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("github name: %s, read user api response body err: %v\n", githubName, err)
		return
	}

	err = json.Unmarshal(respData, &githubResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("github name: %s, parse user api response body err: %v\n", githubName, err)
		return
	}

	avatarResp, err := http.Get(githubResponse.AvatarUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("github name: %s, read avatar api response body err: %v\n", githubName, err)
		return
	}
	defer avatarResp.Body.Close()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", avatarResp.ContentLength))
	w.Header().Set("access-control-allow-origin", "*")
	io.Copy(w, avatarResp.Body)
}

func main() {
	http.HandleFunc("/", HandleProxyAvatar)
	http.ListenAndServe(":8080", nil)
}
