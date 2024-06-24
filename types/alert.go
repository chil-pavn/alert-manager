package types

type Alert struct {
    Annotations map[string]string `json:"annotations"`
    Labels      map[string]string `json:"labels"`
    // Additional fields can be added here
}

// Additional methods for Alert can be defined here if needed
