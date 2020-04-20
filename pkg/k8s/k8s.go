package k8s

import (
	"errors"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	ErrNoClusterConnectivity = errors.New("not able to connect to Kubernetes Cluster")
)

type KubeAPI struct {
	Config clientcmd.ClientConfig
	Clientset *kubernetes.Clientset
}

func NewKubeAPI(kubeConfig KubeConfig) (*KubeAPI, error) {
	config, err := BuildConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	kubeapi := &KubeAPI{
		Config: config,
		Clientset: clientset,
	}

	// HACK: Checking the connectivity of cluster
	_, err = kubeapi.GetNamespaces()
	if err != nil {
		return nil, ErrNoClusterConnectivity
	}

	return kubeapi, nil
}

func BuildConfig(kubeConfig KubeConfig) (clientcmd.ClientConfig, error) {
	var (
		// config *rest.Config
		clientConfig clientcmd.ClientConfig
		err    error
	)

	if kubeConfig.Type == "MANIFEST" {
		// Building config from Manifest YAML File Content
		// config, err = clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig.Manifest))
		// if err != nil {
		// 	return nil, err
		// }
		clientConfig, err = clientcmd.NewClientConfigFromBytes([]byte(kubeConfig.Manifest))
		if err != nil {
			return nil, err
		}
	}
	return clientConfig, err
}

func (kubeapi *KubeAPI) GetNamespaces() ([]v1.Namespace, error) {
	namespaceList, err := kubeapi.Clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return namespaceList.Items, nil
}

func (kubeapi *KubeAPI) GetNodes() ([]v1.Node, error) {
	nodeList, err := kubeapi.Clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

func (kubeapi *KubeAPI) GetPods(namespace string) ([]v1.Pod, error) {
	podList, err := kubeapi.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Pods: ")
	// for _, pod := range podList.Items {
	// 	fmt.Println("\t", pod.GetName())
	// }
	return podList.Items, nil
}

func (kubeapi *KubeAPI) GetServices(namespace string) ([]v1.Service, error) {
	serviceList, err := kubeapi.Clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Services: ")
	// for _, service := range serviceList.Items {
	// 	fmt.Println("\t", service.GetName())
	// }
	return serviceList.Items, nil
}

func (kubeapi *KubeAPI) GetContainers(namespace string, podName string) ([]v1.Container, error) {
	pod, err := kubeapi.Clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// fmt.Println("Containers: ")
	// for _, container := range pod.Spec.Containers {
	// 	fmt.Println("\t", container.Name)
	// }
	return pod.Spec.Containers, nil
}

func (kubeapi *KubeAPI) GetContainerLogs(namespace string, podName string, containerName string, out io.Writer) error {
	// tailLines := int64(100)
	podLogOptions := v1.PodLogOptions{
		Container: containerName,
		// TailLines: &tailLines,
	}

	// fmt.Println("Logs: ")
	logRequest := kubeapi.Clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOptions)

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

func (kubeapi *KubeAPI) GetDeployments(namespace v1.Namespace) ([]v1beta1.Deployment, error) {
	deploymentList, err := kubeapi.Clientset.ExtensionsV1beta1().Deployments(namespace.GetName()).List(metav1.ListOptions{})
	// AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Deployments: ")
	for _, deployment := range deploymentList.Items {
		fmt.Println("\t", deployment.GetName())
	}
	return deploymentList.Items, nil
}

func (kubeapi *KubeAPI) DryRun() {
	fmt.Println("Dry Run")

	// Namespaces
	namespaces, _ := kubeapi.GetNamespaces()
	fmt.Println("Namespaces: ")
	for _, namespace := range namespaces {
		fmt.Println("\t", namespace.GetName())
	}

	// Pods
	for _, namespace := range namespaces {
		fmt.Println("\t", namespace.GetName())

		pods, _ := kubeapi.GetPods(namespace.GetName())
		fmt.Println("Pods: ")
		for _, pod := range pods {
			fmt.Println("\t", pod.GetName())
		}
	}

	// Services
	for _, namespace := range namespaces {
		fmt.Println("\t", namespace.GetName())

		services, _ := kubeapi.GetServices(namespace.GetName())
		fmt.Println("Services: ")
		for _, service := range services {
			fmt.Println("\t", service.GetName())
		}
	}

	// k8s.GetContainers(clientset, "kube-system", "kube-apiserver-kind-control-plane")
	// k8s.GetContainers(clientset, "kube-system", "kube-controller-manager-kind-control-plane")
	// k8s.GetContainers(clientset, "kube-system", "kube-scheduler-kind-control-plane")

	// err := k8s.GetContainerLogs(clientset, "kube-system", "kube-apiserver-kind-control-plane", "kube-apiserver", os.Stdout)

	// for _, node := range nodeList.Items {
	// 	fmt.Println("\t", node.GetName())
	// }
	// for _, condition := range node.Status.Conditions {
	// 	if condition.Reason == "KubeletReady" {
	// 		if condition.Status == "True" {
	// 			nodeStatus = "Ready"
	// 		} else if condition.Reason == "False" {
	// 			nodeStatus = "NotReady"
	// 		} else {
	// 			nodeStatus = "Unknown"
	// 		}
	// 	}
	// }
}
