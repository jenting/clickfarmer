package webserver

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

func (a *APIServer) getClicksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusBadRequest)
	}

	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	js, err := json.Marshal(a.cache.values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (a *APIServer) clickColorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, color string) {
	a.cache.mu.Lock()
	defer a.cache.mu.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/plain")
		count := int64(0)
		switch color {
		case "red":
			count = a.cache.values.Red
		case "green":
			count = a.cache.values.Green
		case "blue":
			count = a.cache.values.Blue
		}
		w.Write([]byte(strconv.FormatInt(count, 10)))
		return

	case http.MethodPut:
		switch color {
		case "red":
			a.cache.values.Red++
		case "green":
			a.cache.values.Green++
		case "blue":
			a.cache.values.Blue++
		}
		return

	default:
		http.Error(w, "only GET and PUT allowed", http.StatusBadRequest)
	}
}

func (a *APIServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(a.webdir, "index.html")
	tmpl := template.Must(template.ParseFiles(indexPath))
	tmpl.Execute(w, struct {
		Name   string
		Colors []string
	}{
		Name:   "Jupiter",
		Colors: []string{"red", "green", "blue"},
	})
}
