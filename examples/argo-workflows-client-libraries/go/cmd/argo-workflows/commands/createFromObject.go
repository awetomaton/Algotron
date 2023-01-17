package commands

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/utils/pointer"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	i "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
)

// sending a workflow to argo from an wfv1 Object
func CreateWorkflowFromObject(workflowObject wfv1.Workflow, namespace string, client i.WorkflowInterface, ctx context.Context) {

	createdWf, err := client.Create(ctx, &workflowObject, metav1.CreateOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Workflow %s submitted\n", createdWf.Name)
	fieldSelector := fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", createdWf.Name))
	watchIf, err := client.Watch(ctx, metav1.ListOptions{FieldSelector: fieldSelector.String(), TimeoutSeconds: pointer.Int64(180)})

	if err != nil {
		panic(err.Error())
	}

	defer watchIf.Stop()
	for next := range watchIf.ResultChan() {
		wf, ok := next.Object.(*wfv1.Workflow)
		if !ok {
			continue
		}
		if !wf.Status.FinishedAt.IsZero() {
			fmt.Printf("Workflow %s %s at %v. Duration: %s.\n", wf.Name, wf.Status.Phase, wf.Status.FinishedAt, wf.Status.ResourcesDuration)
			break
		}
	}
}
