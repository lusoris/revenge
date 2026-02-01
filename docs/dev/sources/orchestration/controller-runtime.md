# sigs.k8s.io/controller-runtime

> Source: https://pkg.go.dev/sigs.k8s.io/controller-runtime
> Fetched: 2026-02-01T11:51:04.266157+00:00
> Content-Hash: 28589d4bc1b4c8c4
> Type: html

---

### Overview ¶

  * Getting Started
  * Organization
  * Managers
  * Controllers
  * Reconcilers
  * Clients and Caches
  * Schemes
  * Webhooks
  * Logging and Metrics
  * Testing



Package controllerruntime provides tools to construct Kubernetes-style controllers that manipulate both Kubernetes CRDs and aggregated/built-in Kubernetes APIs. 

It defines easy helpers for the common use cases when building CRDs, built on top of customizable layers of abstraction. Common cases should be easy, and uncommon cases should be possible. In general, controller-runtime tries to guide users towards Kubernetes controller best-practices. 

#### Getting Started ¶

The main entrypoint for controller-runtime is this root package, which contains all of the common types needed to get started building controllers: 
    
    
    import (
    	ctrl "sigs.k8s.io/controller-runtime"
    )
    

The examples in this package walk through a basic controller setup. The kubebuilder book (<https://book.kubebuilder.io>) has some more in-depth walkthroughs. 

controller-runtime favors structs with sane defaults over constructors, so it's fairly common to see structs being used directly in controller-runtime. 

#### Organization ¶

A brief-ish walkthrough of the layout of this library can be found below. Each package contains more information about how to use it. 

Frequently asked questions about using controller-runtime and designing controllers can be found at <https://github.com/kubernetes-sigs/controller-runtime/blob/main/FAQ.md>. 

#### Managers ¶

Every controller and webhook is ultimately run by a Manager (pkg/manager). A manager is responsible for running controllers and webhooks, and setting up common dependencies, like shared caches and clients, as well as managing leader election (pkg/leaderelection). Managers are generally configured to gracefully shut down controllers on pod termination by wiring up a signal handler (pkg/manager/signals). 

#### Controllers ¶

Controllers (pkg/controller) use events (pkg/event) to eventually trigger reconcile requests. They may be constructed manually, but are often constructed with a Builder (pkg/builder), which eases the wiring of event sources (pkg/source), like Kubernetes API object changes, to event handlers (pkg/handler), like "enqueue a reconcile request for the object owner". Predicates (pkg/predicate) can be used to filter which events actually trigger reconciles. There are pre-written utilities for the common cases, and interfaces and helpers for advanced cases. 

#### Reconcilers ¶

Controller logic is implemented in terms of Reconcilers (pkg/reconcile). A Reconciler implements a function which takes a reconcile Request containing the name and namespace of the object to reconcile, reconciles the object, and returns a Response or an error indicating whether to requeue for a second round of processing. 

#### Clients and Caches ¶

Reconcilers use Clients (pkg/client) to access API objects. The default client provided by the manager reads from a local shared cache (pkg/cache) and writes directly to the API server, but clients can be constructed that only talk to the API server, without a cache. The Cache will auto-populate with watched objects, as well as when other structured objects are requested. The default split client does not promise to invalidate the cache during writes (nor does it promise sequential create/get coherence), and code should not assume a get immediately following a create/update will return the updated resource. Caches may also have indexes, which can be created via a FieldIndexer (pkg/client) obtained from the manager. Indexes can be used to quickly and easily look up all objects with certain fields set. Reconcilers may retrieve event recorders (pkg/recorder) to emit events using the manager. 

#### Schemes ¶

Clients, Caches, and many other things in Kubernetes use Schemes (pkg/scheme) to associate Go types to Kubernetes API Kinds (Group-Version-Kinds, to be specific). 

#### Webhooks ¶

Similarly, webhooks (pkg/webhook/admission) may be implemented directly, but are often constructed using a builder (pkg/webhook/admission/builder). They are run via a server (pkg/webhook) which is managed by a Manager. 

#### Logging and Metrics ¶

Logging (pkg/log) in controller-runtime is done via structured logs, using a log set of interfaces called logr (<https://pkg.go.dev/github.com/go-logr/logr>). While controller-runtime provides easy setup for using Zap (<https://go.uber.org/zap>, pkg/log/zap), you can provide any implementation of logr as the base logger for controller-runtime. 

Metrics (pkg/metrics) provided by controller-runtime are registered into a controller-runtime-specific Prometheus metrics registry. The manager can serve these by an HTTP endpoint, and additional metrics may be registered to this Registry as normal. 

#### Testing ¶

You can easily build integration and unit tests for your controllers and webhooks using the test Environment (pkg/envtest). This will automatically stand up a copy of etcd and kube-apiserver, and provide the correct options to connect to the API server. It's designed to work well with the Ginkgo testing framework, but should work with any testing setup. 

Example ¶

This example creates a simple application Controller that is configured for ReplicaSets and Pods. 

* Create a new application for ReplicaSets that manages Pods owned by the ReplicaSet and calls into ReplicaSetReconciler. 

* Start the application. 
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    
    	appsv1 "k8s.io/api/apps/v1"
    	corev1 "k8s.io/api/core/v1"
    
    	ctrl "sigs.k8s.io/controller-runtime"
    	"sigs.k8s.io/controller-runtime/pkg/client"
    
    	// since we invoke tests with -ginkgo.junit-report we need to import ginkgo.
    	_ "github.com/onsi/ginkgo/v2"
    )
    
    func main() {
    	log := ctrl.Log.WithName("builder-examples")
    
    	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
    	if err != nil {
    		log.Error(err, "could not create manager")
    		os.Exit(1)
    	}
    
    	err = ctrl.
    		NewControllerManagedBy(manager). // Create the Controller
    		For(&appsv1.ReplicaSet{}).       // ReplicaSet is the Application API
    		Owns(&corev1.Pod{}).             // ReplicaSet owns Pods created by it
    		Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
    	if err != nil {
    		log.Error(err, "could not create controller")
    		os.Exit(1)
    	}
    
    	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
    		log.Error(err, "could not start manager")
    		os.Exit(1)
    	}
    }
    
    // ReplicaSetReconciler is a simple Controller example implementation.
    type ReplicaSetReconciler struct {
    	client.Client
    }
    
    // Implement the business logic:
    // This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
    // to a ReplicaSet.
    //
    // * Read the ReplicaSet
    // * Read the Pods
    // * Set a Label on the ReplicaSet with the Pod count.
    func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    
    	rs := &appsv1.ReplicaSet{}
    	err := a.Get(ctx, req.NamespacedName, rs)
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	pods := &corev1.PodList{}
    	err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
    	err = a.Update(ctx, rs)
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	return ctrl.Result{}, nil
    }
    

Share Format Run

Example (CustomHandler) ¶

This example creates a simple application Controller that is configured for ExampleCRDWithConfigMapRef CRD. Any change in the configMap referenced in this Custom Resource will cause the re-reconcile of the parent ExampleCRDWithConfigMapRef due to the implementation of the .Watches method of "sigs.k8s.io/controller-runtime/pkg/builder".Builder. 
    
    
    package main
    
    import (
    	"context"
    	"encoding/json"
    	"os"
    
    	corev1 "k8s.io/api/core/v1"
    
    	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    	"k8s.io/apimachinery/pkg/runtime"
    	"k8s.io/apimachinery/pkg/types"
    
    	ctrl "sigs.k8s.io/controller-runtime"
    	"sigs.k8s.io/controller-runtime/pkg/client"
    	"sigs.k8s.io/controller-runtime/pkg/handler"
    	"sigs.k8s.io/controller-runtime/pkg/reconcile"
    
    	// since we invoke tests with -ginkgo.junit-report we need to import ginkgo.
    	_ "github.com/onsi/ginkgo/v2"
    )
    
    type ExampleCRDWithConfigMapRef struct {
    	metav1.TypeMeta   `json:",inline"`
    	metav1.ObjectMeta `json:"metadata,omitempty"`
    	ConfigMapRef      corev1.LocalObjectReference `json:"configMapRef"`
    }
    
    func deepCopyObject(arg any) runtime.Object {
    
    	argBytes, err := json.Marshal(arg)
    	if err != nil {
    		panic(err)
    	}
    	out := &ExampleCRDWithConfigMapRefList{}
    	if err := json.Unmarshal(argBytes, out); err != nil {
    		panic(err)
    	}
    	return out
    }
    
    // DeepCopyObject implements client.Object.
    func (in *ExampleCRDWithConfigMapRef) DeepCopyObject() runtime.Object {
    	return deepCopyObject(in)
    }
    
    type ExampleCRDWithConfigMapRefList struct {
    	metav1.TypeMeta `json:",inline"`
    	metav1.ListMeta `json:"metadata,omitempty"`
    	Items           []ExampleCRDWithConfigMapRef `json:"items"`
    }
    
    // DeepCopyObject implements client.ObjectList.
    func (in *ExampleCRDWithConfigMapRefList) DeepCopyObject() runtime.Object {
    	return deepCopyObject(in)
    }
    
    func main() {
    	log := ctrl.Log.WithName("builder-examples")
    
    	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
    	if err != nil {
    		log.Error(err, "could not create manager")
    		os.Exit(1)
    	}
    
    	err = ctrl.
    		NewControllerManagedBy(manager).
    		For(&ExampleCRDWithConfigMapRef{}).
    		Watches(&corev1.ConfigMap{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, cm client.Object) []ctrl.Request {
    			// map a change from referenced configMap to ExampleCRDWithConfigMapRef, which causes its re-reconcile
    			crList := &ExampleCRDWithConfigMapRefList{}
    			if err := manager.GetClient().List(ctx, crList); err != nil {
    				manager.GetLogger().Error(err, "while listing ExampleCRDWithConfigMapRefs")
    				return nil
    			}
    
    			reqs := make([]ctrl.Request, 0, len(crList.Items))
    			for _, item := range crList.Items {
    				if item.ConfigMapRef.Name == cm.GetName() {
    					reqs = append(reqs, ctrl.Request{
    						NamespacedName: types.NamespacedName{
    							Namespace: item.GetNamespace(),
    							Name:      item.GetName(),
    						},
    					})
    				}
    			}
    
    			return reqs
    		})).
    		Complete(reconcile.Func(func(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
    			// Your business logic to implement the API by creating, updating, deleting objects goes here.
    			return reconcile.Result{}, nil
    		}))
    	if err != nil {
    		log.Error(err, "could not create controller")
    		os.Exit(1)
    	}
    
    	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
    		log.Error(err, "could not start manager")
    		os.Exit(1)
    	}
    }
    

Share Format Run

Example (UpdateLeaderElectionDurations) ¶

This example creates a simple application Controller that is configured for ReplicaSets and Pods. This application controller will be running leader election with the provided configuration in the manager options. If leader election configuration is not provided, controller runs leader election with default values. Default values taken from: <https://github.com/kubernetes/component-base/blob/master/config/v1alpha1/defaults.go> * defaultLeaseDuration = 15 * time.Second * defaultRenewDeadline = 10 * time.Second * defaultRetryPeriod = 2 * time.Second 

* Create a new application for ReplicaSets that manages Pods owned by the ReplicaSet and calls into ReplicaSetReconciler. 

* Start the application. 
    
    
    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    	"time"
    
    	appsv1 "k8s.io/api/apps/v1"
    	corev1 "k8s.io/api/core/v1"
    
    	ctrl "sigs.k8s.io/controller-runtime"
    	"sigs.k8s.io/controller-runtime/pkg/client"
    
    	// since we invoke tests with -ginkgo.junit-report we need to import ginkgo.
    	_ "github.com/onsi/ginkgo/v2"
    )
    
    func main() {
    	log := ctrl.Log.WithName("builder-examples")
    	leaseDuration := 100 * time.Second
    	renewDeadline := 80 * time.Second
    	retryPeriod := 20 * time.Second
    	manager, err := ctrl.NewManager(
    		ctrl.GetConfigOrDie(),
    		ctrl.Options{
    			LeaseDuration: &leaseDuration,
    			RenewDeadline: &renewDeadline,
    			RetryPeriod:   &retryPeriod,
    		})
    	if err != nil {
    		log.Error(err, "could not create manager")
    		os.Exit(1)
    	}
    
    	err = ctrl.
    		NewControllerManagedBy(manager). // Create the Controller
    		For(&appsv1.ReplicaSet{}).       // ReplicaSet is the Application API
    		Owns(&corev1.Pod{}).             // ReplicaSet owns Pods created by it
    		Complete(&ReplicaSetReconciler{Client: manager.GetClient()})
    	if err != nil {
    		log.Error(err, "could not create controller")
    		os.Exit(1)
    	}
    
    	if err := manager.Start(ctrl.SetupSignalHandler()); err != nil {
    		log.Error(err, "could not start manager")
    		os.Exit(1)
    	}
    }
    
    // ReplicaSetReconciler is a simple Controller example implementation.
    type ReplicaSetReconciler struct {
    	client.Client
    }
    
    // Implement the business logic:
    // This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
    // to a ReplicaSet.
    //
    // * Read the ReplicaSet
    // * Read the Pods
    // * Set a Label on the ReplicaSet with the Pod count.
    func (a *ReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    
    	rs := &appsv1.ReplicaSet{}
    	err := a.Get(ctx, req.NamespacedName, rs)
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	pods := &corev1.PodList{}
    	err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
    	err = a.Update(ctx, rs)
    	if err != nil {
    		return ctrl.Result{}, err
    	}
    
    	return ctrl.Result{}, nil
    }
    

Share Format Run

### Index ¶

  * Variables
  * func NewWebhookManagedBy[T runtime.Object](mgr manager.Manager, obj T) *builder.WebhookBuilder[T]
  * type Builder
  * type GroupResource
  * type GroupVersion
  * type Manager
  * type ObjectMeta
  * type Options
  * type Request
  * type Result
  * type SchemeBuilder
  * type TypeMeta



### Examples ¶

  * Package
  * Package (CustomHandler)
  * Package (UpdateLeaderElectionDurations)



### Constants ¶

This section is empty.

### Variables ¶

[View Source](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L72)
    
    
    var (
    	// RegisterFlags registers flag variables to the given FlagSet if not already registered.
    	// It uses the default command line FlagSet, if none is provided. Currently, it only registers the kubeconfig flag.
    	RegisterFlags = [config](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config).[RegisterFlags](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config#RegisterFlags)
    
    	// GetConfigOrDie creates a *rest.Config for talking to a Kubernetes apiserver.
    	// If --kubeconfig is set, will use the kubeconfig file at that location.  Otherwise will assume running
    	// in cluster and use the cluster provided kubeconfig.
    	//
    	// The returned `*rest.Config` has client-side ratelimting disabled as we can rely on API priority and
    	// fairness. Set its QPS to a value equal or bigger than 0 to re-enable it.
    	//
    	// Will log an error and exit if there is an error creating the rest.Config.
    	GetConfigOrDie = [config](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config).[GetConfigOrDie](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config#GetConfigOrDie)
    
    	// GetConfig creates a *rest.Config for talking to a Kubernetes apiserver.
    	// If --kubeconfig is set, will use the kubeconfig file at that location.  Otherwise will assume running
    	// in cluster and use the cluster provided kubeconfig.
    	//
    	// The returned `*rest.Config` has client-side ratelimting disabled as we can rely on API priority and
    	// fairness. Set its QPS to a value equal or bigger than 0 to re-enable it.
    	//
    	// Config precedence
    	//
    	// * --kubeconfig flag pointing at a file
    	//
    	// * KUBECONFIG environment variable pointing at a file
    	//
    	// * In-cluster config if running in cluster
    	//
    	// * $HOME/.kube/config if exists.
    	GetConfig = [config](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config).[GetConfig](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/client/config#GetConfig)
    
    	// NewControllerManagedBy returns a new controller builder that will be started by the provided Manager.
    	NewControllerManagedBy = [builder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder).[ControllerManagedBy](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder#ControllerManagedBy)
    
    	// NewManager returns a new Manager for creating Controllers.
    	// Note that if ContentType in the given config is not set, "application/vnd.kubernetes.protobuf"
    	// will be used for all built-in resources of Kubernetes, and "application/json" is for other types
    	// including all CRD resources.
    	NewManager = [manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager).[New](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager#New)
    
    	// CreateOrPatch creates or patches the given object obj in the Kubernetes
    	// cluster. The object's desired state should be reconciled with the existing
    	// state using the passed in ReconcileFn. obj must be a struct pointer so that
    	// obj can be patched with the content returned by the Server.
    	//
    	// It returns the executed operation and an error.
    	CreateOrPatch = [controllerutil](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil).[CreateOrPatch](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil#CreateOrPatch)
    
    	// CreateOrUpdate creates or updates the given object obj in the Kubernetes
    	// cluster. The object's desired state should be reconciled with the existing
    	// state using the passed in ReconcileFn. obj must be a struct pointer so that
    	// obj can be updated with the content returned by the Server.
    	//
    	// It returns the executed operation and an error.
    	CreateOrUpdate = [controllerutil](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil).[CreateOrUpdate](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil#CreateOrUpdate)
    
    	// SetControllerReference sets owner as a Controller OwnerReference on owned.
    	// This is used for garbage collection of the owned object and for
    	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
    	// Since only one OwnerReference can be a controller, it returns an error if
    	// there is another OwnerReference with Controller flag set.
    	SetControllerReference = [controllerutil](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil).[SetControllerReference](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/controller/controllerutil#SetControllerReference)
    
    	// SetupSignalHandler registers for SIGTERM and SIGINT. A context is returned
    	// which is canceled on one of these signals. If a second signal is caught, the program
    	// is terminated with exit code 1.
    	SetupSignalHandler = [signals](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager/signals).[SetupSignalHandler](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager/signals#SetupSignalHandler)
    
    	// Log is the base logger used by controller-runtime.  It delegates
    	// to another logr.Logger.  You *must* call SetLogger to
    	// get any actual logging.
    	Log = [log](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log).[Log](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log#Log)
    
    	// LoggerFrom returns a logger with predefined values from a context.Context.
    	// The logger, when used with controllers, can be expected to contain basic information about the object
    	// that's being reconciled like:
    	// - `reconciler group` and `reconciler kind` coming from the For(...) object passed in when building a controller.
    	// - `name` and `namespace` from the reconciliation request.
    	//
    	// This is meant to be used with the context supplied in a struct that satisfies the Reconciler interface.
    	LoggerFrom = [log](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log).[FromContext](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log#FromContext)
    
    	// LoggerInto takes a context and sets the logger as one of its keys.
    	//
    	// This is meant to be used in reconcilers to enrich the logger within a context with additional values.
    	LoggerInto = [log](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log).[IntoContext](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log#IntoContext)
    
    	// SetLogger sets a concrete logging implementation for all deferred Loggers.
    	SetLogger = [log](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log).[SetLogger](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/log#SetLogger)
    )

### Functions ¶

####  func [NewWebhookManagedBy](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L166) ¶ added in v0.2.0
    
    
    func NewWebhookManagedBy[T [runtime](/k8s.io/apimachinery/pkg/runtime).[Object](/k8s.io/apimachinery/pkg/runtime#Object)](mgr [manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager).[Manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager#Manager), obj T) *[builder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder).[WebhookBuilder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder#WebhookBuilder)[T]

NewWebhookManagedBy returns a new webhook builder for the provided type T. 

### Types ¶

####  type [Builder](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L34) ¶
    
    
    type Builder = [builder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder).[Builder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/builder#Builder)

Builder builds an Application ControllerManagedBy (e.g. Operator) and returns a manager.Manager to start it. 

####  type [GroupResource](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L59) ¶
    
    
    type GroupResource = [schema](/k8s.io/apimachinery/pkg/runtime/schema).[GroupResource](/k8s.io/apimachinery/pkg/runtime/schema#GroupResource)

GroupResource specifies a Group and a Resource, but does not force a version. This is useful for identifying concepts during lookup stages without having partially valid types. 

####  type [GroupVersion](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L55) ¶
    
    
    type GroupVersion = [schema](/k8s.io/apimachinery/pkg/runtime/schema).[GroupVersion](/k8s.io/apimachinery/pkg/runtime/schema#GroupVersion)

GroupVersion contains the "group" and the "version", which uniquely identifies the API. 

####  type [Manager](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L46) ¶
    
    
    type Manager = [manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager).[Manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager#Manager)

Manager initializes shared dependencies such as Caches and Clients, and provides them to Runnables. A Manager is required to create Controllers. 

####  type [ObjectMeta](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L70) ¶
    
    
    type ObjectMeta = [metav1](/k8s.io/apimachinery/pkg/apis/meta/v1).[ObjectMeta](/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta)

ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create. 

####  type [Options](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L49) ¶
    
    
    type Options = [manager](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager).[Options](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/manager#Options)

Options are the arguments for creating a new Manager. 

####  type [Request](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L39) ¶
    
    
    type Request = [reconcile](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile).[Request](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile#Request)

Request contains the information necessary to reconcile a Kubernetes object. This includes the information to uniquely identify the object - its Name and Namespace. It does NOT contain information about any specific Event or the object contents itself. 

####  type [Result](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L42) ¶
    
    
    type Result = [reconcile](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile).[Result](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile#Result)

Result contains the result of a Reconciler invocation. 

####  type [SchemeBuilder](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L52) ¶
    
    
    type SchemeBuilder = [scheme](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/scheme).[Builder](/sigs.k8s.io/controller-runtime@v0.23.1/pkg/scheme#Builder)

SchemeBuilder builds a new Scheme for mapping go types to Kubernetes GroupVersionKinds. 

####  type [TypeMeta](https://github.com/kubernetes-sigs/controller-runtime/blob/v0.23.1/alias.go#L66) ¶
    
    
    type TypeMeta = [metav1](/k8s.io/apimachinery/pkg/apis/meta/v1).[TypeMeta](/k8s.io/apimachinery/pkg/apis/meta/v1#TypeMeta)

TypeMeta describes an individual object in an API response or request with strings representing the type of the object and its API schema version. Structures that are versioned or persisted should inline TypeMeta. 

+k8s:deepcopy-gen=false 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
