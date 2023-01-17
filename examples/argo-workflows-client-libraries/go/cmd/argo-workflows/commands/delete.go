package commands

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	i "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
)

// Deletes a Workflow by Name
func DeleteWorkFlow(client i.WorkflowInterface, workflowname string, ctx context.Context) {

	deleteWf := client.Delete(ctx, workflowname, metav1.DeleteOptions{})

	hasError := deleteWf

	if hasError != nil {
		fmt.Printf("Workflow Error: %v\n", hasError)
	} else {
		fmt.Printf("Workflow %v Deleted\n", workflowname)
	}

}
