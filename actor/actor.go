package actor

import (
	"fmt"
	"log"
	"os"

	"github.com/chil-pavn/alert-manager/types"
	"github.com/chil-pavn/alert-manager/utils"
	"github.com/chil-pavn/alert-manager/utils/kubernetes"
)

func TakeAction(alert types.Alert) string {
	if alert.Labels["alertname"] == "KubePodCrashLooping" {
		podName := alert.Labels["pod"]
		namespace := alert.Labels["namespace"]
		cpuDecreased := kubernetes.DecreasePodCPU(namespace, podName)
		memoryDecreased := kubernetes.DecreasePodMemory(namespace, podName)
		imageCorrected := kubernetes.CorrectPodImage(namespace, podName)

		message := fmt.Sprintf("Alert received: %s\nAction taken:\nCPU decreased: %t\nMemory decreased: %t\nImage corrected: %t",
			alert.Annotations["description"], cpuDecreased, memoryDecreased, imageCorrected)

		SendToSlack(message)
		return message
	}

	return "No action taken"
}

func DecreaseMemory(alert types.Alert) string{
	podName := alert.Labels["pod"]
	namespace := alert.Labels["namespace"]
	err := kubernetes.DecreasePodMemory(namespace, podName)
	if err != nil{
		fmt.Printf("unable to decrease memory: %v", err)
	}
	message := fmt.Sprintf("Alert received: %s\nAction taken:\nMemory decreased: %s",
			alert.Annotations["description"], "true")
			SendToSlack(message)
	return ""
}
func SendToSlack(message string){
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err := utils.SendSlackMessage(webhookURL, message)
	if err!=nil{
		log.Printf("unable to send slack message: %v", err)
	}
}

func IncreaseMemory(alert types.Alert) error{
	// find deployment or statefulset from the pod name
	podName := alert.Labels["pod"]
	namespace := alert.Labels["namespace"]
	containerName := alert.Labels["container"]
	kind, name, err := kubernetes.GetPodKindAndName(namespace, podName)

	if err != nil{
		return err
	}
	switch kind {
	case "Deployment":
		kubernetes.EditDeployment(name, namespace, containerName)
	case "StatefulSet":
		kubernetes.EditStatefulset(name, namespace, containerName)
	// increase the memory
	default:
		return fmt.Errorf("owner resource Kind not supported: [%s]", kind)
	}
	// rollout deplotment or statefulset 
	return nil
}
