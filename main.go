package main

import (
	"embed"
	"encoding/json/v2"
	"log"
	"net/http"
	"io/fs"
)

//go:embed static/*
var siteContent embed.FS

type service struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Status int `json:"status"`
}

func main() {
	services := []service{
		service{Name: "kaamkiya.dev", Url: "https://kaamkiya.dev/"},
		service{Name: "Git",          Url: "https://git.kaamkiya.dev/"},
		service{Name: "Trilium",      Url: "http://192.168.2.196:13809"},
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

	siteContentFS, err := fs.Sub(siteContent, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /", http.FileServer(http.FS(siteContentFS)))

	http.ListenAndServe(":29278", mux)
}
