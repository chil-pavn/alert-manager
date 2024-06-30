package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chil-pavn/alert-manager/receiver"
	"github.com/chil-pavn/alert-manager/utils"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Start the health check goroutine
    go startHealthCheck("default", "alertmanager")

    http.HandleFunc("/webhook", receiver.HandleWebhook)
    port := ":5001"
    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func startHealthCheck(namespace, serviceName string) {
    config, err := rest.InClusterConfig()
    if err != nil {
        logrus.Debugf("Failed to create in-cluster config: %v", err)
        config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
        if err != nil {
            panic(err.Error())
    }
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")

    for {
        checkServiceHealth(clientset, namespace, serviceName, slackWebhookURL)
        time.Sleep(30 * time.Second) // Adjust the interval as needed
    }
}

func checkServiceHealth(clientset *kubernetes.Clientset, namespace, serviceName, slackWebhookURL string) {
    _, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get service %s: %v", serviceName, err)
        err  = utils.SendSlackMessage(slackWebhookURL, fmt.Sprintf("Warning: Failed to get service %s status", serviceName))
        if err != nil{
            log.Printf("unable to send slack message %s:", err)
        }
        return
    }

    endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get endpoints for service %s: %v", serviceName, err)
        err = utils.SendSlackMessage(slackWebhookURL, fmt.Sprintf("Warning: Failed to get endpoints for service %s", serviceName))
		if err != nil{
            log.Printf("unable to send slack message %s:", err)
        }
        return
    }

    if len(endpoints.Subsets) == 0 {
        log.Printf("No endpoints available for service %s", serviceName)
        err := utils.SendSlackMessage(slackWebhookURL, fmt.Sprintf("Warning: No endpoints available for service %s", serviceName))
		if err != nil{
            log.Printf("unable to send slack message %s:", err)
        }
        return
    }

    for _, subset := range endpoints.Subsets {
        for _, address := range subset.Addresses {
            url := fmt.Sprintf("http://%s:%d", address.IP, subset.Ports[0].Port)
            if !isServiceAlive(url) {
                message := fmt.Sprintf("Warning: Service %s is not healthy. Endpoint %s is not responsive", serviceName, url)
                log.Println(message)
                err := utils.SendSlackMessage(slackWebhookURL, message)
				if err != nil{
                    log.Printf("unable to send slack message %s:", err)
                }
            }
        }
    }
}

func isServiceAlive(url string) bool {
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        return false
    }
    return true
}
  