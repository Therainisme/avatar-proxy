package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func HandleProxyAvatar(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/\S*.png`)
	matchArr := re.FindStringSubmatch(r.URL.Path)
	if len(matchArr) != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "url format err: example <https://"+r.Host+"/Therainisme.png>\n")
		return
	}

	githubName := matchArr[0][1 : len(matchArr[0])-4]

	avatar, err := avatarMemo.Get(githubName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	avatarBytes := avatar.([]byte)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(avatarBytes)))
	w.Header().Set("access-control-allow-origin", "*")

	w.Write(avatarBytes)
}

func HandleProxyAvatarLogs(w http.ResponseWriter, r *http.Request) {
	serverResponse := ServerResponse{
		Data: logStore.Logs,
		Code: 0,
		Msg:  "success",
	}
	serverResponseByte, _ := json.Marshal(serverResponse)

	w.Header().Set("access-control-allow-origin", "*")
	fmt.Fprintf(w, "%s", string(serverResponseByte))
}
