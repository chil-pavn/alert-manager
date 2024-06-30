package alertmanager

import (
	"fmt"
	"log"

	"github.com/chil-pavn/alert-manager/actor"
	"github.com/chil-pavn/alert-manager/enricher"
	"github.com/chil-pavn/alert-manager/types"
	"github.com/chil-pavn/alert-manager/utils/kubernetes"
)

func ManageAlert(a types.Alert) string {
    if a.Name() == "KubePodCrashLooping"{
		podName := a.Labels["pod"]
		namespace := a.Labels["namespace"]
		lastStateReason, currStateReason, err := kubernetes.CheckCrashLoopBackOff(podName, namespace)
		if err!=nil {
			return fmt.Sprintf("unable to check crashloop backoff reason: %v", err)
		}

		log.Printf("crashloop backoff reason: %s", currStateReason)

		if lastStateReason == "Error"{
			actor.SendToSlack(fmt.Sprintf("pod %s is in crashloopback off. Last state terminated reason : %s. Current Waiting state reason: %s", podName, lastStateReason, currStateReason))
			return "successfully alerted in slack"
		}

		if lastStateReason == "OOMKilled"{
			enricher.GenericEnricher(a)
			fmt.Printf("current memory limits: %s ", a.Annotations["memory_usage"])
			fmt.Printf("Taking action")
			actionResult := actor.TakeAction(a)
			return actionResult
		}
		// add more ways to manage this particular alert
    }
	return "Alert management done"
}
