package receiver

import (
	"encoding/json"
	"net/http"

	"github.com/chil-pavn/alert-manager/actor"
	"github.com/chil-pavn/alert-manager/enricher"
	"github.com/chil-pavn/alert-manager/types"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var alert types.Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enrichedAlert := enricher.EnrichAlert(alert)
	actionResult := actor.TakeAction(enrichedAlert)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "result": actionResult})
}
