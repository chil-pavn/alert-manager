package enricher

import (
	"log"

	"github.com/chil-pavn/alert-manager/types"
	"github.com/chil-pavn/alert-manager/utils/cloud"
	"github.com/chil-pavn/alert-manager/utils/kubernetes"
)


func GenericEnricher(a types.Alert) {
	podName := a.Labels["pod"]
	namespace := a.Labels["namespace"]

	// Get additional data from Kubernetes
	cpuUsage, err := kubernetes.GetPodCPUUsage(namespace, podName)
	if err != nil {
		log.Print(err.Error())
	}
	memoryUsage, err := kubernetes.GetPodMemoryUsage(namespace, podName)
	if err != nil {
		log.Print(err.Error())
	}
	podDetails, err := kubernetes.GetPodDetails(namespace, podName)
	if err != nil {
		log.Print(err.Error())
	}

	// Enrich the alert with the additional data
	a.Annotations["cpu_usage"] = cpuUsage
	a.Annotations["memory_usage"] = memoryUsage
	a.Annotations["pod_details"] = podDetails
}

func CpuEnricher(a types.Alert) {
	podName := a.Labels["pod"]
	namespace := a.Labels["namespace"]

	// Get additional data from Kubernetes
	cpuUsage, err := kubernetes.GetPodCPUUsage(namespace, podName)
	if err != nil{
		log.Printf(err.Error())
	}
	// Get cloud details also
	NodeCpuUsage := cloud.GetCpuUsage(namespace, podName)

	a.Annotations["cpu_usage"] = cpuUsage
	a.Annotations["NodeCpuUsage"] = NodeCpuUsage
}
