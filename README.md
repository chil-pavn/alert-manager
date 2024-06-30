# alert manager
 custom alert manager

we have a single type alert
# main.go

1. Observing alertmanager which sends the alerts to the service and sends slack messages based on the health check on it

2. starting the webhook 

# receiver.go

receives all the alerts and sends only alerts which needs to be managed to the alertmanager module and sends the response based on the result from the alertmanager

# alertmanager.go

contains the core logic of how to handle the specific alerts by utilizing the enrichment and actor methods available and sends the result back to the server. 

Add any other logic to handle a different alert here. Eg: Pending, if pending, get reason -> Memory unavailable -> get node and pod metrics if needed -> autoscale the node 
or 
Pending -> failed to pull image -> check the network settings or the image availability -> send slack message

# enricher.go

contains various enrichment functions which can be called by the alertmanager to take a decision. Add any other enrichment functions needed.

# actor.go

contains actions that can be possible which the alertmanager will call to take action. Add actions that can be taken by the alertmanager

# utils

contains common cloud and kubernetes functions to be used by the enrichment and actor methods

# kubernetes and test

has the yaml files to deploy the alert-manager (this service), prometheus and alertmanager to test on minikube. 




# example command used for testing, replace the pod, clsuter and other details as required

curl -X POST -H "Content-Type: application/json" -d '{
  "annotations": {
    "description": "Pod customer/customer-rs-transformer-9b75b488c-cpfd7 (rs-transformer) is restarting 2.11 times / 10 minutes.",
    "runbook_url": "https://github.com/kubernetes-monitoring/kubernetes-mixin/tree/master/runbook.md#alert-name-kubepodcrashlooping",
    "summary": "Pod is crash looping."
  },
  "labels": {
    "alertname": "KubePodCrashLooping",
    "cluster": "minikube",
    "container": "crashloop",
    "endpoint": "http",
    "job": "kube-state-metrics",
    "namespace": "default",
    "pod": "test-crashloop-pod-7fb96ccb6c-h7t7m",
    "priority": "P0",
    "prometheus": "monitoring/kube-prometheus-stack-prometheus",
    "region": "us-west-1",
    "replica": "0",
    "service": "kube-prometheus-stack-kube-state-metrics",
    "severity": "CRITICAL"
  },
  "startsAt": "2022-03-02T07:31:57.339Z",
  "status": "firing"
}' http://localhost:5001/webhook

# enable metrics server in minikube if not enabled
error fetching CPU usage: the server could not find the requested resource -> minikube addons enable metrics-server