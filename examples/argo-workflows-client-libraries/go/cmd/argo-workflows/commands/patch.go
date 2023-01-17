package commands

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	i "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
)

func PatchWorkflowFromObject(workflowObject wfv1.Workflow, namespace string, client i.WorkflowInterface, ctx context.Context) {

	// patching a workflow to argo from an wfv1 Object, if its json then no need
	marshalledWorkflow := wfv1.MustMarshallJSON(workflowObject)
	x := []byte(marshalledWorkflow)

	//insert your object to be updated/patched
	createdWf, err := client.Patch(ctx, "hello-world-thxj6", types.MergePatchType, x, metav1.PatchOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Workflow %s patched\n", createdWf.Name)

}
