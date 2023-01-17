package commands

import (
	"context"
	"fmt"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	i "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	gjson "github.com/tidwall/gjson"
)

// List Workflows
func ListWorkFlows(client i.WorkflowInterface, ctx context.Context) {

	//change label selector for basic filtering
	listWf, err := client.List(ctx, metav1.ListOptions{Limit: 0, LabelSelector: "workflows.argoproj.io/phase=Succeeded"})
	if err != nil {
		panic(err.Error())
	}

	cleanJSON := wfv1.MustMarshallJSON(listWf)
	metaData := gjson.Get(cleanJSON, "items.#.metadata")
	metaData.ForEach(func(key, value gjson.Result) bool {

		name := gjson.Get(value.String(), "name")
		statusArray := gjson.Get(value.String(), "labels|@values").Array()

		completed := statusArray[0]
		phase := statusArray[1]

		fmt.Printf("Name is: %v\t Completed: %v\t Phase:%v\n", name, completed, phase)
		return true
	})

}
