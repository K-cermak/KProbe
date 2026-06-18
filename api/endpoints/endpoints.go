package endpoints

import (
	"encoding/json"
	"net/http"
	"os"

	"KProbeAPI/db"
	"KProbeAPI/helpers"

	"github.com/go-chi/chi/v5"
)

type StatusResponse struct {
	ProbeName string `json:"probe_name"`
	Time      string `json:"time"`
	ScanName  string `json:"scan_name"`
	ScanTime  string `json:"check"`
	Status    string `json:"status"`
}

func ServeEditor(w http.ResponseWriter, r *http.Request) {
	// if html file is not found, return 404
	_, err := os.Stat("/opt/kprobe/editor.html")
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := map[string]interface{}{
			"error":   404,
			"message": "Editor not found",
		}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	htmlPath := "/opt/kprobe/editor.html"
	http.ServeFile(w, r, htmlPath)
}

func ServeStatus(w http.ResponseWriter, r *http.Request, probeName string) {
	scanName := chi.URLParam(r, "scan_name")
	data, correct := db.GetScanNewest(scanName)
	if correct == db.DB_CONNECTION_FAILED {
		helpers.PrintError("Failed to get scan from database")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse := map[string]interface{}{
			"probe_name": probeName,
			"time":      helpers.GetCurrTime(),
			"error":     500,
			"message":   "Internal server error",
		}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	if correct == db.DB_SCAN_NEWEST_FAILED {
		helpers.PrintWarning("Failed to get scan from database with name " + scanName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		jsonResponse := map[string]interface{}{
			"probe_name": probeName,
			"time":      helpers.GetCurrTime(),
			"error":     404,
			"message":   "Not found, probably scan with this name does not exist",
		}
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	resp := StatusResponse{
		ProbeName: probeName,
		Time:      helpers.GetCurrTime(),
		ScanName:  scanName,
		ScanTime:  data.Generated,
		Status:    helpers.BoolToString(data.Passed),
	}

	helpers.PrintSuccess("Served status for scan " + scanName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
