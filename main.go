package main

import (
	_ "embed"
	"encoding/json/v2"
	"log"
	"net/http"
)

//go:embed index.html
var homePage []byte

type service struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Status int `json:"status"`
}

func main() {
	services := []service{
		service{Name: "kaamkiya.dev", Url: "https://kaamkiya.dev/"},
		service{Name: "Git",          Url: "https://git.kaamkiya.dev/"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /data", func(w http.ResponseWriter, req *http.Request) {
		for i, s := range services {
			res, err := http.Get(s.Url)
			if err != nil {
				services[i].Status = 0
			} else {
				services[i].Status = res.StatusCode
			}
		}

		data, err := json.Marshal(services)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(homePage)); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":29278", mux)
}
