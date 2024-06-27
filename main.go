package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getEnv retrieves the value of an environment variable by key
func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

func main() {
	// Set up logging
	log := logrus.New()
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Initialize Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Retrieve necessary environment variables
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		log.Fatalf("HOSTNAME environment variable is not set")
	}
	checkIntervalStr, err := getEnv("CHECK_INTERVAL")
	if err != nil {
		log.Fatalf("Failed to get CHECK_INTERVAL: %v", err)
	}
	checkInterval, err := time.ParseDuration(checkIntervalStr)
	if err != nil {
		log.Fatalf("Failed to parse CHECK_INTERVAL: %v", err)
	}
	ipAddress, err := getEnv("IP_ADDRESS")
	if err != nil {
		log.Fatalf("Failed to get IP_ADDRESS: %v", err)
	}
	nodeLabelKey, err := getEnv("NODE_LABEL_KEY")
	if err != nil {
		log.Fatalf("Failed to get NODE_LABEL_KEY: %v", err)
	}
	nodeLabelValueActive, err := getEnv("NODE_LABEL_VALUE_ACTIVE")
	if err != nil {
		log.Fatalf("Failed to get NODE_LABEL_VALUE_ACTIVE: %v", err)
	}
	nodeLabelValueInactive, err := getEnv("NODE_LABEL_VALUE_INACTIVE")
	if err != nil {
		log.Fatalf("Failed to get NODE_LABEL_VALUE_INACTIVE: %v", err)
	}
	waitTimeStr, err := getEnv("WAIT_TIME")
	if err != nil {
		log.Fatalf("Failed to get WAIT_TIME: %v", err)
	}
	waitTime, err := time.ParseDuration(waitTimeStr)
	if err != nil {
		log.Fatalf("Failed to parse WAIT_TIME: %v", err)
	}
	resourceNamespace, err := getEnv("RESOURCE_NAMESPACE")
	if err != nil {
		log.Fatalf("Failed to get RESOURCE_NAMESPACE: %v", err)
	}
	resourceName, err := getEnv("RESOURCE_NAME")
	if err != nil {
		log.Fatalf("Failed to get RESOURCE_NAME: %v", err)
	}
	resourceType, err := getEnv("RESOURCE_TYPE")
	if err != nil {
		log.Fatalf("Failed to get RESOURCE_TYPE: %v", err)
	}

	// Main loop for checking and updating resources
	for {
		// Refresh node information from Kubernetes API
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), hostname, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Failed to get node information: %v", err)
			time.Sleep(checkInterval)
			continue
		}

		ipPresent := checkIPPresence(ipAddress, log)
		labelValue := nodeLabelValueInactive
		if ipPresent {
			labelValue = nodeLabelValueActive
		}

		currentLabelValue, exists := node.Labels[nodeLabelKey]
		if !exists || currentLabelValue != labelValue {
			log.Debugf("Current label value is '%s', updating to '%s'", currentLabelValue, labelValue)
			errUpdateLabel := updateNodeLabel(clientset, hostname, nodeLabelKey, labelValue, log)
			if errUpdateLabel != nil {
				log.Errorf("Failed to update node label: %v", errUpdateLabel)
			} else {
				log.Infof("Node label updated successfully: %s=%s", nodeLabelKey, labelValue)
				if ipPresent {
					time.Sleep(waitTime)
					if resourceType == "deployment" {
						errRestartDeployment := restartDeployment(clientset, resourceNamespace, resourceName, log)
						if errRestartDeployment != nil {
							log.Errorf("Failed to restart deployment: %v", errRestartDeployment)
						} else {
							log.Infof("Deployment restarted successfully: %s/%s", resourceNamespace, resourceName)
						}
					} else if resourceType == "statefulset" {
						errRestartStatefulSet := restartStatefulSet(clientset, resourceNamespace, resourceName, log)
						if errRestartStatefulSet != nil {
							log.Errorf("Failed to restart statefulset: %v", errRestartStatefulSet)
						} else {
							log.Infof("StatefulSet restarted successfully: %s/%s", resourceNamespace, resourceName)
						}
					}
				}
			}
		} else {
			log.Debugf("Label %s is already set to %s, no update needed", nodeLabelKey, labelValue)
		}

		time.Sleep(checkInterval)
	}
}

// checkIPPresence checks if the given IP address is present in any network interface
func checkIPPresence(ip string, log *logrus.Logger) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Errorf("Failed to get network interfaces: %v", err)
		return false
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Errorf("Failed to get addresses for interface %s: %v", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			if strings.Contains(addr.String(), ip) {
				return true
			}
		}
	}
	return false
}

// updateNodeLabel updates the label of a Kubernetes node
func updateNodeLabel(clientset *kubernetes.Clientset, nodeName, labelKey, labelValue string, log *logrus.Logger) error {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node %s: %v", nodeName, err)
	}

	if node.Labels == nil {
		node.Labels = make(map[string]string)
	}
	node.Labels[labelKey] = labelValue

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update node %s: %v", nodeName, err)
	}
	return nil
}

// getCurrentNodeLabelValue retrieves the current value of a node label
func getCurrentNodeLabelValue(clientset *kubernetes.Clientset, nodeName, labelKey string) (string, error) {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get node %s: %v", nodeName, err)
	}

	value, exists := node.Labels[labelKey]
	if !exists {
		return "", nil
	}
	return value, nil
}

// restartDeployment restarts a Kubernetes deployment
func restartDeployment(clientset *kubernetes.Clientset, namespace, name string, log *logrus.Logger) error {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s: %v", name, err)
	}
	if deployment == nil {
		return fmt.Errorf("deployment %s is nil", name)
	}

	// Add a check to ensure the deployment exists before proceeding
	if deployment.Status.Replicas == 0 {
		return fmt.Errorf("deployment %s has no replicas", name)
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment %s: %v", name, err)
	}
	return nil
}

// restartStatefulSet restarts a Kubernetes statefulset
func restartStatefulSet(clientset *kubernetes.Clientset, namespace, name string, log *logrus.Logger) error {
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get statefulset %s: %v", name, err)
	}
	if statefulSet == nil {
		return fmt.Errorf("statefulset %s is nil", name)
	}

	// Add a check to ensure the statefulset exists before proceeding
	if statefulSet.Status.Replicas == 0 {
		return fmt.Errorf("statefulset %s has no replicas", name)
	}

	if statefulSet.Spec.Template.Annotations == nil {
		statefulSet.Spec.Template.Annotations = make(map[string]string)
	}
	statefulSet.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update statefulSet %s: %v", name, err)
	}
	return nil
}

