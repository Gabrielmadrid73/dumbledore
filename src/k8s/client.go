package k8s

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var K8sClient *kubernetes.Clientset

func InitK8sClient() error {
	var kubeconfig *string

	if os.Getenv("LOCAL") == "true" {
		home := homedir.HomeDir()
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		K8sClient = client
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		K8sClient = client
	}
	log.Println("Initialized Kubernetes Client")
	return nil
}
