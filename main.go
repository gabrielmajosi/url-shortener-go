package main

import (
	"embed"
	"github.com/dchest/uniuri"
	"io/fs"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

const serverAddr = ":8080"
const storeFile = "store.csv"

//go:embed web
var web embed.FS

// TODO: use dirty flag to avoid reading from disk periodically
func periodicSaveToCsv(store *map[string]string) {
	for {
		// sleep
		time.Sleep(12 * time.Second)

		var currentDisk = LoadStore()
		if !reflect.DeepEqual(currentDisk, *store) {
			slog.Info("New entries written to file", "location", storeFile, "mem", *store, "disk", currentDisk)
			SaveStore(*store)
		}
	}
}

func main() {
	var store = LoadStore()
	go periodicSaveToCsv(&store)
	slog.Info("Loaded store", "data", store)

	webRoot, _ := fs.Sub(web, "web")
	http.Handle("/", http.FileServer(http.FS(webRoot)))

	http.HandleFunc("GET /l/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		destinationUrl := store[slug]

		http.Redirect(w, r, destinationUrl, http.StatusPermanentRedirect)
	})

	http.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := r.Form.Get("url")
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// generate the slug
		slug := uniuri.New()

		store[slug] = url
		slog.Info("Stored", "slug", slug, "destination", url)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(slug))
	})

	slog.Info("Listening on " + serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		slog.Error("Couldn't start server", "addr", serverAddr, "error", err)
		return
	}

}
