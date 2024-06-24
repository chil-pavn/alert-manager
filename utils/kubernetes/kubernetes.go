package kubernetes

import (
    "fmt"
    "os/exec"
)

func GetPodCPUUsage(namespace, podName string) string {
    // Logic to get CPU usage of the pod
    // Example using kubectl
    output, err := exec.Command("kubectl", "top", "pod", podName, "-n", namespace, "--no-headers").Output()
    if err != nil {
        return fmt.Sprintf("Error fetching CPU usage: %v", err)
    }
    return string(output)
}

func GetPodMemoryUsage(namespace, podName string) string {
    // Logic to get memory usage of the pod
    // Example using kubectl
    output, err := exec.Command("kubectl", "top", "pod", podName, "-n", namespace, "--no-headers").Output()
    if err != nil {
        return fmt.Sprintf("Error fetching memory usage: %v", err)
    }
    return string(output)
}

func GetPodDetails(namespace, podName string) string {
    // Logic to get details of the pod
    // Example using kubectl
    output, err := exec.Command("kubectl", "describe", "pod", podName, "-n", namespace).Output()
    if err != nil {
        return fmt.Sprintf("Error fetching pod details: %v", err)
    }
    return string(output)
}

func DecreasePodCPU(namespace, podName string) bool {
    // Logic to decrease CPU limits of the pod
    // Example using kubectl
    err := exec.Command("kubectl", "set", "resources", "pod", podName, "-n", namespace, "--limits=cpu=100m").Run()
    return err == nil
}

func DecreasePodMemory(namespace, podName string) bool {
    // Logic to decrease memory limits of the pod
    // Example using kubectl
    err := exec.Command("kubectl", "set", "resources", "pod", podName, "-n", namespace, "--limits=memory=128Mi").Run()
    return err == nil
}

func CorrectPodImage(namespace, podName string) bool {
    // Logic to correct the image of the pod
    // Example using kubectl
    err := exec.Command("kubectl", "set", "image", "pod", podName, "-n", namespace, fmt.Sprintf("%s=correct-image:latest", podName)).Run()
    return err == nil
}
