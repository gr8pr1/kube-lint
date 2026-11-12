package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfigPath := clientcmd.RecommendedHomeFile

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading kubeconfig: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "building clientset: %v\n", err)
		os.Exit(1)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "listing pods: %v\n", err)
		os.Exit(1)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAMESPACE\tNAME\tPHASE\tRESTARTS\tAGE\tNODE")
	for _, pod := range pods.Items {
		nodeName := pod.Spec.NodeName
		if nodeName == "" {
			nodeName = "<none>"
		}

		var restarts int32
		// Init containers are excluded to match kubectl's RESTARTS column.
		for _, containerStatus := range pod.Status.ContainerStatuses {
			restarts += containerStatus.RestartCount
		}

		age := duration.HumanDuration(time.Since(pod.CreationTimestamp.Time))

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", pod.Namespace, pod.Name, pod.Status.Phase, restarts, age, nodeName)
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing output: %v\n", err)
		os.Exit(1)
	}
}
