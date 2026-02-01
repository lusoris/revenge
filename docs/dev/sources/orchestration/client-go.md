# k8s.io/client-go

> Source: https://pkg.go.dev/k8s.io/client-go
> Fetched: 2026-02-01T11:51:01.994475+00:00
> Content-Hash: 1a738e4bdb7baf20
> Type: html

---

### Overview ¶

  * Key Packages
  * Connecting to the API
  * Interacting with API Objects
  * Handling API Errors
  * Building Controllers



Package clientgo is the official Go client for the Kubernetes API. It provides a standard set of clients and tools for building applications, controllers, and operators that communicate with a Kubernetes cluster. 

#### Key Packages ¶

  * kubernetes: Contains the typed Clientset for interacting with built-in, versioned API objects (e.g., Pods, Deployments). This is the most common starting point. 

  * dynamic: Provides a dynamic client that can perform operations on any Kubernetes object, including Custom Resources (CRs). It is essential for building controllers that work with CRDs. 

  * discovery: Used to discover the API groups, versions, and resources supported by a Kubernetes cluster. 

  * tools/cache: The foundation of the controller pattern. This package provides efficient caching and synchronization mechanisms (Informers and Listers) for building controllers. 

  * tools/clientcmd: Provides methods for loading client configuration from kubeconfig files. This is essential for out-of-cluster applications and CLI tools. 

  * rest: Provides a lower-level RESTClient that manages the details of communicating with the Kubernetes API server. It is useful for advanced use cases that require fine-grained control over requests, such as working with non-standard REST verbs. 




#### Connecting to the API ¶

There are two primary ways to configure a client to connect to the API server: 

  1. In-Cluster Configuration: For applications running inside a Kubernetes pod, the `rest.InClusterConfig()` function provides a straightforward way to configure the client. It automatically uses the pod's service account for authentication and is the recommended approach for controllers and operators. 

  2. Out-of-Cluster Configuration: For local development or command-line tools, the `clientcmd` package is used to load configuration from a kubeconfig file. 




The `rest.Config` object allows for fine-grained control over client-side performance and reliability. Key settings include: 

  * QPS: The maximum number of queries per second to the API server.
  * Burst: The maximum number of queries that can be issued in a single burst.
  * Timeout: The timeout for individual requests.



#### Interacting with API Objects ¶

Once configured, a client can be used to interact with objects in the cluster. 

  * The Typed Clientset (`kubernetes` package) provides a strongly typed interface for working with built-in Kubernetes objects. 

  * The Dynamic Client (`dynamic` package) can work with any object, including Custom Resources, using `unstructured.Unstructured` types. 

  * For Custom Resources (CRDs), the `k8s.io/code-generator` repository contains the tools to generate typed clients, informers, and listers. The `sample-controller` is the canonical example of this pattern. 

  * Server-Side Apply is a patching strategy that allows multiple actors to share management of an object by tracking field ownership. This prevents actors from inadvertently overwriting each other's changes and provides a mechanism for resolving conflicts. The `applyconfigurations` package provides the necessary tools for this declarative approach. 




#### Handling API Errors ¶

Robust error handling is essential when interacting with the API. The `k8s.io/apimachinery/pkg/api/errors` package provides functions to inspect errors and check for common conditions, such as whether a resource was not found or already exists. This allows controllers to implement robust, idempotent reconciliation logic. 

#### Building Controllers ¶

The controller pattern is central to Kubernetes. A controller observes the state of the cluster and works to bring it to the desired state. 

  * The `tools/cache` package provides the building blocks for this pattern. Informers watch the API server and maintain a local cache, Listers provide read-only access to the cache, and Workqueues decouple event detection from processing. 

  * In a high-availability deployment where multiple instances of a controller are running, leader election (`tools/leaderelection`) is used to ensure that only one instance is active at a time. 

  * Client-side feature gates allow for enabling or disabling experimental features in `client-go`. They can be configured via the `rest.Config` object. 


Example (HandlingAPIErrors) ¶
    
    
    package main
    
    import (
    	"fmt"
    
    	"k8s.io/apimachinery/pkg/api/errors"
    	"k8s.io/apimachinery/pkg/runtime/schema"
    )
    
    func main() {
    	// This example demonstrates how to handle common API errors.
    	// Create a NotFound error to simulate a failed API request.
    	err := errors.NewNotFound(schema.GroupResource{Group: "v1", Resource: "pods"}, "my-pod")
    
    	// The errors.IsNotFound() function checks if an error is a NotFound error.
    	if errors.IsNotFound(err) {
    		fmt.Println("Pod not found")
    	}
    
    	// The errors package provides functions for other common API errors, such as:
    	// - errors.IsAlreadyExists(err)
    	// - errors.IsConflict(err)
    	// - errors.IsServerTimeout(err)
    
    }
    
    
    
    Output:
    
    Pod not found
    

Share Format Run

Example (InClusterConfiguration) ¶
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    
    	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    	"k8s.io/client-go/kubernetes"
    	"k8s.io/client-go/rest"
    )
    
    func main() {
    	// This example demonstrates how to create a clientset for an application
    	// running inside a Kubernetes cluster. It uses the pod's service account
    	// for authentication.
    
    	// rest.InClusterConfig() returns a configuration object that can be used to
    	// create a clientset. It is the recommended way to configure a client for
    	// in-cluster applications.
    	config, err := rest.InClusterConfig()
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Create the clientset.
    	clientset, err := kubernetes.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Use the clientset to interact with the API.
    
    	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	fmt.Printf("There are %d pods in the default namespace\n", len(pods.Items))
    }
    

Share Format Run

Example (LeaderElection) ¶
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"log"
    	"os"
    	"os/signal"
    	"path/filepath"
    	"syscall"
    	"time"
    
    	"k8s.io/client-go/kubernetes"
    	"k8s.io/client-go/tools/clientcmd"
    	"k8s.io/client-go/tools/leaderelection"
    	"k8s.io/client-go/tools/leaderelection/resourcelock"
    )
    
    func main() {
    	// This example demonstrates the leader election pattern. A controller running
    	// multiple replicas uses leader election to ensure that only one replica is
    	// active at a time.
    
    	// Configure the client (out-of-cluster for this example).
    	var kubeconfig string
    	if home := os.Getenv("HOME"); home != "" {
    		kubeconfig = filepath.Join(home, ".kube", "config")
    	}
    	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    	if err != nil {
    		fmt.Printf("Error building kubeconfig: %s\n", err.Error())
    		return
    	}
    	clientset, err := kubernetes.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error building clientset: %s\n", err.Error())
    		return
    	}
    
    	// The unique name of the controller writing to the Lease object.
    	id := "my-controller"
    
    	// The namespace and name of the Lease object.
    	leaseNamespace := "default"
    	leaseName := "my-controller-lease"
    
    	// Create a context that can be cancelled to stop the leader election.
    	ctx, cancel := context.WithCancel(context.Background())
    	defer cancel()
    
    	// Set up a signal handler to cancel the context on termination.
    	sigCh := make(chan os.Signal, 1)
    	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    	go func() {
    		<-sigCh
    		log.Println("Received termination signal, shutting down...")
    		cancel()
    	}()
    
    	// Create the lock object.
    	lock, err := resourcelock.New(resourcelock.LeasesResourceLock,
    		leaseNamespace,
    		leaseName,
    		clientset.CoreV1(),
    		clientset.CoordinationV1(),
    		resourcelock.ResourceLockConfig{
    			Identity: id,
    		})
    	if err != nil {
    		fmt.Printf("Error creating lock: %v\n", err)
    		return
    	}
    
    	// Create the leader elector.
    	elector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
    		Lock:          lock,
    		LeaseDuration: 15 * time.Second,
    		RenewDeadline: 10 * time.Second,
    		RetryPeriod:   2 * time.Second,
    		Callbacks: leaderelection.LeaderCallbacks{
    			OnStartedLeading: func(ctx context.Context) {
    				// This function is called when the controller becomes the leader.
    				// You would start your controller's main logic here.
    				log.Println("Became leader, starting controller.")
    				// This is a simple placeholder for the controller's work.
    				for {
    					select {
    					case <-time.After(5 * time.Second):
    						log.Println("Doing controller work...")
    					case <-ctx.Done():
    						log.Println("Controller stopped.")
    						return
    					}
    				}
    			},
    			OnStoppedLeading: func() {
    				// This function is called when the controller loses leadership.
    				// You should stop any active work and gracefully shut down.
    				log.Printf("Lost leadership, shutting down.")
    			},
    			OnNewLeader: func(identity string) {
    				// This function is called when a new leader is elected.
    				if identity != id {
    					log.Printf("New leader elected: %s", identity)
    				}
    			},
    		},
    	})
    	if err != nil {
    		fmt.Printf("Error creating leader elector: %v\n", err)
    		return
    	}
    
    	// Start the leader election loop.
    	elector.Run(ctx)
    }
    

Share Format Run

Example (OutOfClusterConfiguration) ¶
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    	"path/filepath"
    
    	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    	"k8s.io/client-go/kubernetes"
    	"k8s.io/client-go/tools/clientcmd"
    )
    
    func main() {
    	// This example demonstrates how to create a clientset for an application
    	// running outside of a Kubernetes cluster, using a kubeconfig file. This is
    	// the standard approach for local development and command-line tools.
    
    	// The default location for the kubeconfig file is in the user's home directory.
    	var kubeconfig string
    	if home := os.Getenv("HOME"); home != "" {
    		kubeconfig = filepath.Join(home, ".kube", "config")
    	}
    
    	// Create the client configuration from the kubeconfig file.
    	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Configure client-side rate limiting.
    	config.QPS = 50
    	config.Burst = 100
    
    	// A clientset contains clients for all the API groups and versions supported
    	// by the cluster.
    	clientset, err := kubernetes.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Use the clientset to interact with the API.
    
    	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	fmt.Printf("There are %d pods in the default namespace\n", len(pods.Items))
    }
    

Share Format Run

Example (ServerSideApply) ¶
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    	"path/filepath"
    
    	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    
    	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
    
    	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
    
    	metav1ac "k8s.io/client-go/applyconfigurations/meta/v1"
    	"k8s.io/client-go/kubernetes"
    	"k8s.io/client-go/tools/clientcmd"
    )
    
    func main() {
    	// This example demonstrates how to use Server-Side Apply to declaratively
    	// manage a Deployment object. Server-Side Apply is the recommended approach
    	// for controllers and operators to manage objects.
    
    	// Configure the client (out-of-cluster for this example).
    	var kubeconfig string
    	if home := os.Getenv("HOME"); home != "" {
    		kubeconfig = filepath.Join(home, ".kube", "config")
    	}
    	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    	clientset, err := kubernetes.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Define the desired state of the Deployment using the applyconfigurations package.
    	// This provides a typed, structured way to build the patch.
    	deploymentName := "my-app"
    	replicas := int32(2)
    	image := "nginx:1.14.2"
    
    	// The FieldManager is a required field that identifies the controller managing
    	// this object's state.
    	fieldManager := "my-controller"
    
    	// Build the apply configuration.
    	deploymentApplyConfig := appsv1ac.Deployment(deploymentName, "default").
    		WithSpec(appsv1ac.DeploymentSpec().
    			WithReplicas(replicas).
    			WithSelector(metav1ac.LabelSelector().WithMatchLabels(map[string]string{"app": "my-app"})).
    			WithTemplate(corev1ac.PodTemplateSpec().
    				WithLabels(map[string]string{"app": "my-app"}).
    				WithSpec(corev1ac.PodSpec().
    					WithContainers(corev1ac.Container().
    						WithName("nginx").
    						WithImage(image),
    					),
    				),
    			),
    		)
    
    	// Perform the Server-Side Apply patch. The PatchType must be types.ApplyPatchType.
    	// The context, name, apply configuration, and patch options are required.
    	result, err := clientset.AppsV1().Deployments("default").Apply(
    		context.TODO(),
    		deploymentApplyConfig,
    		metav1.ApplyOptions{FieldManager: fieldManager},
    	)
    
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	fmt.Printf("Deployment %q applied successfully.\n", result.Name)
    }
    

Share Format Run

Example (UsingDynamicClient) ¶
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    	"path/filepath"
    
    	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    	"k8s.io/apimachinery/pkg/runtime/schema"
    	"k8s.io/client-go/dynamic"
    	"k8s.io/client-go/tools/clientcmd"
    )
    
    func main() {
    	// This example demonstrates how to create and use a dynamic client to work
    	// with Custom Resources or other objects without needing their Go type
    	// definitions.
    
    	// Configure the client (out-of-cluster for this example).
    	var kubeconfig string
    	if home := os.Getenv("HOME"); home != "" {
    		kubeconfig = filepath.Join(home, ".kube", "config")
    	}
    	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Create a dynamic client.
    	dynamicClient, err := dynamic.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Define the GroupVersionResource for the object you want to access.
    	// For Pods, this is {Group: "", Version: "v1", Resource: "pods"}.
    	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
    
    	// Use the dynamic client to list all pods in the "default" namespace.
    	// The result is an UnstructuredList.
    	list, err := dynamicClient.Resource(gvr).Namespace("default").List(context.TODO(), metav1.ListOptions{})
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Iterate over the list and print the name of each pod.
    	for _, item := range list.Items {
    		name, found, err := unstructured.NestedString(item.Object, "metadata", "name")
    		if err != nil || !found {
    			fmt.Printf("Could not find name for pod: %v\n", err)
    			continue
    		}
    		fmt.Printf("Pod Name: %s\n", name)
    	}
    }
    

Share Format Run

Example (UsingInformers) ¶
    
    
    package main
    
    import (
    	"fmt"
    	"log"
    	"os"
    	"os/signal"
    	"path/filepath"
    	"syscall"
    	"time"
    
    	"k8s.io/client-go/informers"
    	"k8s.io/client-go/kubernetes"
    	"k8s.io/client-go/tools/cache"
    	"k8s.io/client-go/tools/clientcmd"
    )
    
    func main() {
    	// This example demonstrates the basic pattern for using an informer to watch
    	// for changes to Pods. This is a conceptual example; a real controller would
    	// have more robust logic and a workqueue.
    
    	// Configure the client (out-of-cluster for this example).
    	var kubeconfig string
    	if home := os.Getenv("HOME"); home != "" {
    		kubeconfig = filepath.Join(home, ".kube", "config")
    	}
    	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    	clientset, err := kubernetes.NewForConfig(config)
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// A SharedInformerFactory provides a shared cache for multiple informers,
    	// which reduces memory and network overhead.
    	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
    	podInformer := factory.Core().V1().Pods().Informer()
    
    	_, err = podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    		AddFunc: func(obj interface{}) {
    			key, err := cache.MetaNamespaceKeyFunc(obj)
    			if err == nil {
    				log.Printf("Pod ADDED: %s", key)
    			}
    		},
    		UpdateFunc: func(oldObj, newObj interface{}) {
    			key, err := cache.MetaNamespaceKeyFunc(newObj)
    			if err == nil {
    				log.Printf("Pod UPDATED: %s", key)
    			}
    		},
    		DeleteFunc: func(obj interface{}) {
    			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
    			if err == nil {
    				log.Printf("Pod DELETED: %s", key)
    			}
    		},
    	})
    	if err != nil {
    		fmt.Printf("Error encountered: %v\n", err)
    		return
    	}
    
    	// Graceful shutdown requires a two-channel pattern.
    	//
    	// The first channel, `sigCh`, is used by the `signal` package to send us
    	// OS signals (e.g., Ctrl+C). This channel must be of type `chan os.Signal`.
    	//
    	// The second channel, `stopCh`, is used to tell the informer factory to
    	// stop. The informer factory's `Start` method expects a channel of type
    	// `<-chan struct{}`. It will stop when this channel is closed.
    	//
    	// The goroutine below is the "translator" that connects these two channels.
    	// It waits for a signal on `sigCh`, and when it receives one, it closes
    	// `stopCh`, which in turn tells the informer factory to shut down.
    	sigCh := make(chan os.Signal, 1)
    	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    	stopCh := make(chan struct{})
    	go func() {
    		<-sigCh
    		close(stopCh)
    	}()
    
    	// Start the informer.
    	factory.Start(stopCh)
    
    	// Wait for the initial cache sync.
    	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
    		log.Println("Timed out waiting for caches to sync")
    		return
    	}
    
    	log.Println("Informer has synced. Watching for Pod events...")
    
    	// Wait for the stop signal.
    	<-stopCh
    	log.Println("Shutting down...")
    }
    

Share Format Run
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
