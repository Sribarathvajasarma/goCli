package main

import (
	"context"
	"flag"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateApp(clientset *kubernetes.Clientset, name string, image string, port int32) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: port,
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func GetApps(clientset *kubernetes.Clientset) {

	deployments, err4 := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
	if err4 != nil {
		fmt.Println("Error in creating deployments")
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}
}

// func createService() {
// 	kubeconfig := flag.String("kubeconfig", "C:\\Users\\WSO2\\.kube\\config", "location to your cube config files")
// 	config, err1 := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err1 != nil {
// 		//handle error
// 		fmt.Println("Error in creating config")
// 	}

// 	clientset, err2 := kubernetes.NewForConfig(config)
// 	if err2 != nil {
// 		fmt.Println("Error in creating clientset")
// 	}

// 	clientset.CoreV1().Services("kube-system").Create(&apiv1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "Test service",
// 			Namespace: "kube-system",
// 			Labels: map[string]string{
// 				"k8s-app": "kube-controller-manager",
// 			},
// 		},
// 		Spec: apiv1.ServiceSpec{
// 			Ports:     nil,
// 			Selector:  nil,
// 			ClusterIP: "",
// 		},
// 	})

// 	// service := apiv1.Service{
// 	// 	ObjectMeta: metav1.ObjectMeta{
// 	// 		Name:      "myservice",
// 	// 		Namespace: "default",
// 	// 		Labels: map[string]string{
// 	// 			"app": "myapp",
// 	// 		},
// 	// 	},
// 	// 	Spec: apiv1.ServiceSpec{
// 	// 		Ports:     nil,
// 	// 		Selector:  nil,
// 	// 		ClusterIP: "",
// 	// 	},
// 	// }

// 	// serviceClient := clientset.AppsV1().Services(apiv1.NamespaceDefault)

// 	// // Create Service
// 	// fmt.Println("Creating service...")
// 	// result, err := serviceClient.Create(service)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// }

func main() {

	jobName := flag.String("jobname", "", "The name of the job")
	containerImage := flag.String("image", "", "Name of the container image")

	flag.Parse()

	kubeconfig := flag.String("kubeconfig", "C:\\Users\\WSO2\\.kube\\config", "location to your cube config files")
	config, err1 := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err1 != nil {
		//handle error
		fmt.Println("Error in creating config")
	}

	clientset, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		fmt.Println("Error in creating clientset")
	}

	if *jobName != "" && *containerImage != "" {
		CreateApp(clientset, *jobName, *containerImage, 80)

	} else {
		GetApps(clientset)

	}

	// CreateApp("demo2", "nginx:1.12", int32(80))

	// GetApps()

	//createService()

	// kubeconfig := flag.String("kubeconfig", "C:\\Users\\WSO2\\.kube\\config", "location to your cube config files")
	// config, err1 := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err1 != nil {
	// 	//handle error
	// 	fmt.Println("Error in creating config")
	// }

	// clientset, err2 := kubernetes.NewForConfig(config)
	// if err2 != nil {
	// 	fmt.Println("Error in creating clientset")
	// }

	// ctx := context.Background()
	// pods, err3 := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	// if err3 != nil {
	// 	fmt.Println("Error in creating pods")
	// }

	// for _, pod := range pods.Items {
	// 	fmt.Printf("%s", pod.Name)
	// }

	// deployments, err4 := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	// if err4 != nil {
	// 	fmt.Println("Error in creating deployments")
	// }
	// for _, deployment := range deployments.Items {
	// 	fmt.Printf("%s", deployment.Name)
	// }

	//deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	// Create Deployment
	// fmt.Println("Creating deployment...")
	// result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}
