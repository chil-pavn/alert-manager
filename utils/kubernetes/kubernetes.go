package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func getKubernetesClient() (*kubernetes.Clientset, error) {
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
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return clientset, nil
}

func getMetricsClient() (*metrics.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Debugf("Failed to create in-cluster config: %v", err)
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			panic(err.Error())
		}
	}
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return metricsClient, nil
}

func GetPodCPUUsage(namespace, podName string) (string, error) {
	metricsClient, err := getMetricsClient()
	if err != nil {
		return "", err
	}

	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error fetching CPU usage: %w", err)
	}

	if len(podMetrics.Containers) == 0 {
		return "No metrics available for the pod", nil
	}

	cpuUsage := podMetrics.Containers[0].Usage[v1.ResourceCPU]
	return cpuUsage.String(), nil
}

func GetPodMemoryUsage(namespace, podName string) (string, error) {
	metricsClient, err := getMetricsClient()
	if err != nil {
		return "", err
	}

	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error fetching memory usage: %w", err)
	}

	if len(podMetrics.Containers) == 0 {
		return "No metrics available for the pod", nil
	}

	memoryUsage := podMetrics.Containers[0].Usage[v1.ResourceMemory]
	return memoryUsage.String(), nil
}

func GetPodDetails(namespace, podName string) (string, error) {
	clientset, err := getKubernetesClient()
	if err != nil {
		return "", err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("error fetching pod details: %w", err)
	}

	return fmt.Sprintf("Pod details:\nName: %s\nNamespace: %s\nNode: %s\nStatus: %s\n", pod.Name, pod.Namespace, pod.Spec.NodeName, pod.Status.Phase), nil
}

func DecreasePodCPU(namespace, podName string) error {
	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching pod: %w", err)
	}

	if len(pod.Spec.Containers) == 0 {
		return fmt.Errorf("no containers in pod")
	}
    log.Print(pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU])
	pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU] = resource.MustParse("100m")
	_, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating pod CPU limit: %w", err)
	}

	return nil
}

func DecreasePodMemory(namespace, podName string) error {
	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching pod: %w", err)
	}

	if len(pod.Spec.Containers) == 0 {
		return fmt.Errorf("no containers in pod")
	}

	pod.Spec.Containers[0].Resources.Limits[v1.ResourceMemory] = resource.MustParse("128Mi")
	_, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating pod memory limit: %w", err)
	}

	return nil
}

func CorrectPodImage(namespace, podName string) error {
	clientset, err := getKubernetesClient()
	if err != nil {
		return err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error fetching pod: %w", err)
	}

	if len(pod.Spec.Containers) == 0 {
		return fmt.Errorf("no containers in pod")
	}

	pod.Spec.Containers[0].Image = "busybox:latest"
	_, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating pod image: %w", err)
	}

	return nil
}

func CheckCrashLoopBackOff(podName, namespace string) (string, string, error) {
	clientset, err := getKubernetesClient()
	if err != nil {
		return "","", err
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return "","", err
	}

	if pod.Status.ContainerStatuses != nil {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason == "CrashLoopBackOff" {

				log.Printf("errors: %s, %s", containerStatus.LastTerminationState.Terminated.Reason, containerStatus.State.Waiting.Reason )
				return containerStatus.LastTerminationState.Terminated.Reason ,containerStatus.State.Waiting.Reason, nil
			}
		}
	}

	return "", "", nil
}