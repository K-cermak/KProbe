package main

import (
	"encoding/json"
	"log"
	"net/http"

	"KProbeAPI/db"
	"KProbeAPI/endpoints"
	"KProbeAPI/helpers"

	"github.com/go-chi/chi/v5"
)

var ProbeName string

func main() {
	port, ok := db.GetValue("api_port")
	if ok != db.RES_OK {
		helpers.PrintError("Failed to get API port from database (" + ok + ")")
		helpers.PrintWarning("Using default port 8080")
		port = "8080"
	}

	ProbeName, ok = db.GetValue("probe_name")
	if ok != db.RES_OK {
		helpers.PrintError("Failed to get probe name from database (" + ok + ")")
		helpers.PrintWarning("Using default probe name UnknownProbe")
		ProbeName = "UnknownProbe"
	}

	allowedEditor, ok := db.GetValue("editor_endpoint")
	if ok != db.RES_OK {
		helpers.PrintError("Failed to get if editor is allowed from database (" + ok + ")")
		helpers.PrintWarning("Using default value true")
		allowedEditor = "true"
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse := map[string]interface{}{
			"probe_name": ProbeName,
			"time":      helpers.GetCurrTime(),
		}
		json.NewEncoder(w).Encode(jsonResponse)
	})

	if allowedEditor == "true" {
		r.Get("/editor", endpoints.ServeEditor)
	}
	r.Get("/status/{scan_name}", func(w http.ResponseWriter, r *http.Request) {
		endpoints.ServeStatus(w, r, ProbeName)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		jsonResponse := map[string]interface{}{
			"probe_name": ProbeName,
			"time":      helpers.GetCurrTime(),
			"error":     404,
			"message":   "Not found",
		}
		json.NewEncoder(w).Encode(jsonResponse)
	})

	helpers.PrintSuccess("API server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
