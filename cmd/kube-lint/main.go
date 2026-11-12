package main

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigPath := clientcmd.RecommendedHomeFile

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating config: %v\n", err)
		os.Exit(1)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "building clientset: %v\n", err)
		os.Exit(1)
	}

	pods, err := clientSet.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "listing pods: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Number of pods: ", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("Name: %s\n", pod.Name)
	}
}
