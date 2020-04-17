package k8s

import (
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildConfig(kubeConfig KubeConfig) (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)

	if kubeConfig.Type == "MANIFEST" {
		// Building config from Manifest YAML File Content
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig.Manifest))
	}
	return config, err
}

func GetClientset(kubeConfig KubeConfig) (*kubernetes.Clientset, error) {
	config, err := BuildConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// api := clientset.CoreV1()

	return clientset, nil
}

func SearchNamespaces(clientset *kubernetes.Clientset) ([]v1.Namespace, error) {
	namespaceList, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Namespaces: ")
	// for _, pod := range namespaceList.Items {
	// 	fmt.Println("\t", pod.GetName())
	// }
	return namespaceList.Items, nil
}

func SearchPods(clientset *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	podList, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Pods: ")
	// for _, pod := range podList.Items {
	// 	fmt.Println("\t", pod.GetName())
	// }
	return podList.Items, nil
}

func SearchServices(clientset *kubernetes.Clientset, namespace string) ([]v1.Service, error) {
	serviceList, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Services: ")
	// for _, service := range serviceList.Items {
	// 	fmt.Println("\t", service.GetName())
	// }

	return serviceList.Items, nil
}

func GetContainers(clientset *kubernetes.Clientset, namespace string, podName string) ([]v1.Container, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Containers: ")
	for _, container := range pod.Spec.Containers {
		fmt.Println("\t", container.Name)
	}
	return pod.Spec.Containers, nil
}

func GetContainerLogs(clientset *kubernetes.Clientset, namespace string, podName string, containerName string, out io.Writer) error {
	tailLines := int64(100)
	podLogOptions := v1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLines,
	}

	fmt.Println("Logs: ")

	logRequest := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOptions)

	readCloser, err := logRequest.Stream()
	if readCloser != nil {
		defer readCloser.Close()
	}
	if err != nil {
		return err
	}

	_, err = io.Copy(out, readCloser)
	return err
}

func SearchDeployments(clientset *kubernetes.Clientset, namespace v1.Namespace) ([]v1beta1.Deployment, error) {
	deploymentList, err := clientset.ExtensionsV1beta1().Deployments(namespace.GetName()).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployments: ")
	for _, deployment := range deploymentList.Items {
		fmt.Println("\t", deployment.GetName())
	}
	return deploymentList.Items, nil
}
