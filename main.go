package main

import (
	"embed"
	"github.com/dchest/uniuri"
	"io/fs"
	"log/slog"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const serverAddr = ":8080"
const storeFile = "store.csv"

//go:embed web
var web embed.FS

type RecordsStore struct {
	records  map[string]string // slug -> url
	mut      sync.RWMutex
	modified bool
}

func periodicSaveToCsv(store *RecordsStore) {
	for {
		// sleep
		time.Sleep(12 * time.Second)
		if !store.modified {
			continue
		}

		var currentDiskRecords = LoadStore().records
		store.mut.RLock()
		if !reflect.DeepEqual(currentDiskRecords, store.records) {
			store.mut.RUnlock()
			slog.Info("New entries written to file", "location", storeFile, "mem", &store.records, "disk", currentDiskRecords)
			SaveStore(store)
		} else {
			store.mut.RUnlock()
		}
	}
}

func main() {
	var store = LoadStore()
	go periodicSaveToCsv(&store)
	slog.Info("Loaded store", "data", store.records)

	webRoot, _ := fs.Sub(web, "web")
	http.Handle("/", http.FileServer(http.FS(webRoot)))

	http.HandleFunc("GET /l/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		store.mut.RLock()
		destinationUrl := store.records[slug]
		store.mut.RUnlock()

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

		store.mut.Lock()
		store.records[slug] = url
		store.modified = true
		store.mut.Unlock()

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
