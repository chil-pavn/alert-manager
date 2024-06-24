package actor

import (
	"fmt"
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

		sendToSlack(message)
		return message
	}

	return "No action taken"
}

func sendToSlack(message string) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	utils.SendSlackMessage(webhookURL, message)
}
