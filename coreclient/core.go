package coreclient

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var clientSet *client.CoreV1Client

func GetNodePods(node string) ([]v1.Pod, error) {
	if clientSet == nil {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientSet = client.NewForConfigOrDie(config)
	}

	pods, err := clientSet.Pods("").List(metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + node,
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}
