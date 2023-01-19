package main

import (
	"argo/cmd/argo-workflows/commands"
	"context"
	"flag"
	"fmt"
	"os/user"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
)

var createWorkflow = wfv1.Workflow{
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

var patchWorkflow = wfv1.Workflow{
	Spec: wfv1.WorkflowSpec{
		Templates: []wfv1.Template{
			{
				Container: &corev1.Container{
					Command: []string{"cowsay", "hello world Im patched"},
				},
			},
		},
	},
}

func main() {
	// get current user to determine home directory
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}

	// get kubeconfig file location
	kubeconfig := flag.String("kubeconfig", filepath.Join(usr.HomeDir, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()
	fmt.Printf("Kube Config %s Used\n", *kubeconfig)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// set the namespace you want to use
	namespace := "default"

	// create the workflow client
	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows(namespace)
	ctx := context.Background()

	// run the different commands from the /commands folder

	//commands.CreateWorkflowFromObject(createWorkflow, namespace, wfClient, ctx)
	//commands.CreateWorkflowFromFile("cmd/argo-workflows/commands/assets/yaml/coinflip.yaml", namespace, wfClient, ctx)

	//commands.DeleteWorkFlow(wfClient, "hello-world-hzm85", ctx)

	//commands.ListWorkFlows(wfClient, ctx)

	commands.PatchWorkflowFromObject(patchWorkflow, namespace, wfClient, ctx)

}
