package operator

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8sdeploy/pkg/constants"
)

//DeployLogger creates an informer that looks for replicaSet update events
func DeployLogger(clientset *kubernetes.Clientset, release string, deploys []string, today time.Time) error {
	namespaces := len(deploys)
	var wg sync.WaitGroup
	dep := 0

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Apps().V1().ReplicaSets().Informer()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"App", "Namespace", "Status"})

	fmt.Printf("Starting deployment watcher...\n")
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			d := obj.(*v1.ReplicaSet)
			if d.GetCreationTimestamp().After(today) {
				fmt.Printf("Starting deployment in namespace=%s for app=%s at %s \n", d.GetNamespace(), d.GetLabels()["app"], d.GetCreationTimestamp())
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			d := newObj.(*v1.ReplicaSet)
			replicas, _ := strconv.Atoi(d.GetAnnotations()["deployment.kubernetes.io/desired-replicas"])
			r := int32(replicas)
			if d.GetCreationTimestamp().After(today) && release == d.GetLabels()["release"] {
				fmt.Printf("Waiting for deployment %s rollout to finish: %d of %s updated replicas are available...\n", d.GetLabels()["release"], d.Status.ReadyReplicas, d.GetAnnotations()["deployment.kubernetes.io/desired-replicas"])
				// check if ready replicas are equal to the total replicas, and if the na
				if d.Status.AvailableReplicas == d.Status.ReadyReplicas && d.Status.ReadyReplicas == r && find(deploys, d.GetNamespace()) {
					fmt.Printf("Successful Deployment of %s on %s\n", d.GetLabels()["app"], d.GetNamespace())
					dep++
					t.AppendRow([]interface{}{d.GetLabels()["app"], d.GetNamespace(), "Success"})
					deploys = remove(deploys, d.GetNamespace())
				}
			}
			// lets finish up if total success messages equals namespaces entered
			if dep == namespaces {
				fmt.Println("All deployments finished, sutting down watcher gracefully")
				t.SortBy([]table.SortBy{
					{Name: "Namespace", Mode: table.Asc},
				})
				t.Render()
				os.Exit(0)
			}
		}})

	stopper := make(chan struct{})
	go func() {
		informer.Run(stopper)
		defer close(stopper)
		wg.Wait()
	}()

	// stop the Deploywatcher if it hits the timeout, count the ones that did not get deployed as "failed"
	select {
	case <-stopper:
	case <-time.After(constants.WatcherTimeout):
		fmt.Println("Application timed out ..exiting")
		if len(deploys) > 0 {
			for _, x := range deploys {
				t.AppendRow([]interface{}{"", x, "Failed"})
			}
			t.Render()
		}
		os.Exit(0)
	}
	return nil
}

// remove an item in a slice
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// check if a string exists in a slice
func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
