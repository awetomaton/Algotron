package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/argoproj/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
)

// Create Any WorkFlow SPEC here
var helloWorldWorkflow = wfv1.Workflow{
	ObjectMeta: metav1.ObjectMeta{
		GenerateName: "hello-world-",
	},
	Spec: wfv1.WorkflowSpec{
		Entrypoint: "whalesay",
		Templates: []wfv1.Template{
			{
				Name: "whalesay",
				Container: &corev1.Container{
					Image:   "docker/whalesay:latest",
					Command: []string{"cowsay", "hello world"},
				},
			},
		},
	},
}

func main() {
	// get current user to determine home directory
	usr, err := user.Current()
	checkErr(err)

	// get kubeconfig file location
	kubeconfig := flag.String("kubeconfig", filepath.Join(usr.HomeDir, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	//alternate setup if running inside a cluster

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	checkErr(err)
	namespace := "default"

	// create the workflow client
	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(namespace)

	// submit the hello world workflow
	ctx := context.Background()
	createdWf, err := wfClient.Create(ctx, &helloWorldWorkflow, metav1.CreateOptions{})
	checkErr(err)
	fmt.Printf("Workflow %s submitted\n", createdWf.Name)

	// sample sending a workflow to argo from a yaml file
	// create the workflow client
	wfClient2 := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(namespace)
	ctx2 := context.Background()
	file, err := os.ReadFile("coinflip.yaml")
	checkErr(err)
	strFile := string(file)
	coinFlip := wfv1.MustUnmarshalWorkflow(strFile)
	createdCoinFlip, err := wfClient2.Create(ctx2, coinFlip, metav1.CreateOptions{})
	checkErr(err)

	fmt.Printf("Workflow %s submitted\n", createdCoinFlip.Name)

	// wait for the hello world workflow to complete
	fieldSelector := fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", createdWf.Name))
	watchIf, err := wfClient.Watch(ctx, metav1.ListOptions{FieldSelector: fieldSelector.String(), TimeoutSeconds: pointer.Int64(180)})
	errors.CheckError(err)
	defer watchIf.Stop()
	for next := range watchIf.ResultChan() {
		wf, ok := next.Object.(*wfv1.Workflow)
		if !ok {
			continue
		}
		if !wf.Status.FinishedAt.IsZero() {
			fmt.Printf("Workflow %s %s at %v. Message: %s.\n", wf.Name, wf.Status.Phase, wf.Status.FinishedAt, wf.Status.Message)
			break
		}
	}

	// wait for the coinflip workflow to complete
	selector := fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", createdCoinFlip.Name))

	watch, err := wfClient.Watch(ctx, metav1.ListOptions{FieldSelector: selector.String(), TimeoutSeconds: pointer.Int64(180)})
	errors.CheckError(err)
	defer watch.Stop()
	for next := range watch.ResultChan() {
		wf, ok := next.Object.(*wfv1.Workflow)
		if !ok {
			continue
		}
		if !wf.Status.FinishedAt.IsZero() {
			fmt.Printf("Workflow %s %s at %v. Message: %s.\n", wf.Name, wf.Status.Phase, wf.Status.FinishedAt, wf.Status.Message)
			break
		}
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
