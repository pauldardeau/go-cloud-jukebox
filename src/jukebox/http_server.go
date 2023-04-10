package jukebox

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	jukebox *Jukebox
	port    int
}

func NewHttpServer(jukebox *Jukebox, port int) *HttpServer {
	var httpServer HttpServer
	httpServer.jukebox = jukebox
	httpServer.port = port
	return &httpServer
}

func (httpServer *HttpServer) Run() bool {
	http.HandleFunc("/songAdvance", func(w http.ResponseWriter, r *http.Request) {
		if httpServer.jukebox != nil {
			httpServer.jukebox.AdvanceToNextSong()
			fmt.Fprintf(w, "<html><body>advanced to next song</body></html>")
		}
	})

	http.HandleFunc("/togglePausePlay", func(w http.ResponseWriter, r *http.Request) {
		if httpServer.jukebox != nil {
			httpServer.jukebox.TogglePausePlay()
			fmt.Fprintf(w, "<html><body>toggled pause/play</body></html>")
		}
	})

	/*
		http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
			httpServer.jukebox.TogglePausePlay()
		})

		http.HandleFunc("/api/memoryUsage", func(w http.ResponseWriter, r *http.Request) {
			httpServer.jukebox.TogglePausePlay()
		})
	*/

	listenAddress := fmt.Sprintf(":%d", httpServer.port)

	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		return false
	} else {
		fmt.Println("Http server listening...")
		return true
	}
}
