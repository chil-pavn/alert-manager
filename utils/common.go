package utils

import (
    "bytes"
    "net/http"
    "encoding/json"
)

func SendSlackMessage(webhookURL, message string) error {
    payload := map[string]string{"text": message}
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    _, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
    return err
}
