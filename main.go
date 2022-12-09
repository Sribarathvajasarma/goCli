package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	ingress2 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateApp(clientset *kubernetes.Clientset, name string, image string, portNo string, domain string) {
	port, err6 := strconv.Atoi(portNo)
	if err6 != nil {
		fmt.Println("Error while type convertion")
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
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
									ContainerPort: int32(port),
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentsClient := clientset.AppsV1().Deployments("dev")

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	CreateService(clientset, name, port, domain)
}

func GetApps(clientset *kubernetes.Clientset) {

	deployments, err4 := clientset.AppsV1().Deployments("dev").List(context.Background(), metav1.ListOptions{})
	if err4 != nil {
		fmt.Println("Error in creating deployments")
	}
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}
}

func GetServices(clientset *kubernetes.Clientset) {

	services, err4 := clientset.CoreV1().Services("dev").List(context.Background(), metav1.ListOptions{})
	if err4 != nil {
		fmt.Println("Error while getting service")
	}
	for _, service := range services.Items {
		fmt.Println(service.Name)
	}
}

func CreateService(clientset *kubernetes.Clientset, name string, port int, domain string) {
	Service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-service",
			Namespace: "dev",
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol:   apiv1.ProtocolTCP,
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
		},
	}

	fmt.Println("Creating service...")
	result, err := clientset.CoreV1().Services("dev").Create(context.TODO(), Service, metav1.CreateOptions{})

	if err != nil {
		fmt.Println("Error while getting service")
	}

	var service_name string = result.GetObjectMeta().GetName()
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	CreateIngress(clientset, service_name, domain, port)
}

func GetIngresses(clientset *kubernetes.Clientset) {
	ingressList, err := clientset.NetworkingV1().Ingresses("dev").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error while getting ingresses")
	}
	ingressCtrls := ingressList.Items
	if len(ingressCtrls) > 0 {
		for _, ingress := range ingressCtrls {
			fmt.Printf("ingress %s exists in namespace %s\n", ingress.Name, ingress.Namespace)
		}
	} else {
		fmt.Println("no ingress found")
	}
}

func CreateIngress(clientset *kubernetes.Clientset, service string, domain string, port int) {
	var str string = "Prefix"
	ingressService := &ingress2.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service + "-ingress",
			Namespace: "dev",
		},
		Spec: ingress2.IngressSpec{
			Rules: []ingress2.IngressRule{
				{
					Host: domain,
					IngressRuleValue: ingress2.IngressRuleValue{
						HTTP: &ingress2.HTTPIngressRuleValue{
							Paths: []ingress2.HTTPIngressPath{
								{
									Path:     "/",
									PathType: (*ingress2.PathType)(&str),
									Backend: ingress2.IngressBackend{
										Service: &ingress2.IngressServiceBackend{
											Name: service,
											Port: ingress2.ServiceBackendPort{
												Number: int32(port),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			// TLS: []ingress2.IngressTLS{
			// 	{
			// 		Hosts: []string{
			// 			"demo.mlopshub.com",
			// 		},
			// 		SecretName: "hello-app-tls",
			// 	},
			// },
		},
	}

	fmt.Println("Creating Ingress...")
	result, err := clientset.NetworkingV1().Ingresses("dev").Create(context.TODO(), ingressService, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func main() {

	Name := flag.String("name", "", "The name of the app")
	containerImage := flag.String("image", "", "Name of the image")
	port := flag.String("port", "", "Port number")
	domain := flag.String("domain", "", "Domain name")

	//DNS := flag.String("dns", "", "Domain name")

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

	flag.Parse()

	if *Name != "" && *containerImage != "" && *port != "" && *domain != "" {
		// PortNo, err6 := strconv.Atoi(*Port)
		// if err6 != nil {
		// 	fmt.Println("Error while type convertion")
		// }
		CreateApp(clientset, *Name, *containerImage, *port, *domain)
		// hosts, err := txeh.NewHostsDefault()
		// if err != nil {
		// 	panic(err)
		// }

		// hosts.AddHost("127.0.0.1", *domain)
		// hfData := hosts.RenderHostsFile()

		// if you like to see what the outcome will
		// look like
		// fmt.Println(hfData)

		// hosts.Save()

	} else {
		//CreateApp(clientset, "demo", "gcr.io/google-samples/hello-app:2.0", 80, "example.com")
		//CreateIngress(clientset, "demo-service", "example.com", 80)
		//CreateService(clientset, "demo", 80, "example.com")
		GetApps(clientset)
		//GetIngresses(clientset)
		//GetServices(clientset)
		//CreateService(clientset)

	}

	// CreateApp("demo2", "nginx:1.12", int32(80))

	// GetApps()

	//createService()

}
