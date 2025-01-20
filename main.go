package main

import (
	"github.com/dchest/uniuri"
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

const serverAddr = ":8080"
const storeFile = "store.csv"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//_, _ = fmt.Fprintf(w, "SEE API DOCS")

		templ, err := template.ParseFiles("home.html")
		if err != nil {
			slog.Error("Error parsing template", "error", err)
			return
		}

		err = templ.Execute(w, nil)
		if err != nil {
			slog.Error("Error executing template", "error", err)
		}
	})

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
	})

	slog.Info("Listening on " + serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		slog.Error("Couldn't start server", "addr", serverAddr, "error", err)
		return
	}

}
