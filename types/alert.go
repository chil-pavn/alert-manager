package types

type Alert struct {
    Annotations map[string]string `json:"annotations"`
    Labels      map[string]string `json:"labels"`
    
    // Additional fields can be added here
}

func (a *Alert) Name () string{
    return a.Labels["alertname"]
}

// Additional methods for Alert can be defined here if needed
