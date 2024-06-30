package receiver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chil-pavn/alert-manager/alertmanager"
	"github.com/chil-pavn/alert-manager/types"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var alert types.Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// only few alerts are handled from the received
	if contains(AlertNamesToHandle, alert.Name()) {
		func() {
			log.Print("Encriching")
			actionResult := alertmanager.ManageAlert(alert)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "result": actionResult})
		}()
	} else {
		// Return a success response for alerts that are not handled
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "result": "Alert not handled"})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
