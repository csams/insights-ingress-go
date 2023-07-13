package version

import (
	"encoding/json"
	"net/http"
)

// IngressVersion is the json structure for the /version endpoint
type IngressVersion struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func NewHandler(c CompletedConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData, err := json.Marshal(c.IngressVersion)
		if err != nil {
			c.Log.Error("Unable to get version")
			w.Write([]byte(`{"version": "unavailable"}`))
		} else {
			w.Write(jsonData)
		}
	}
}
