package enricher

import (
	"github.com/chil-pavn/alert-manager/types"
	"github.com/chil-pavn/alert-manager/utils/kubernetes"
)

func EnrichAlert(alert types.Alert) types.Alert {
	podName := alert.Labels["pod"]
	namespace := alert.Labels["namespace"]

	// Get additional data from Kubernetes
	cpuUsage := kubernetes.GetPodCPUUsage(namespace, podName)
	memoryUsage := kubernetes.GetPodMemoryUsage(namespace, podName)
	podDetails := kubernetes.GetPodDetails(namespace, podName)

	// Enrich the alert with the additional data
	alert.Annotations["cpu_usage"] = cpuUsage
	alert.Annotations["memory_usage"] = memoryUsage
	alert.Annotations["pod_details"] = podDetails

	return alert
}
