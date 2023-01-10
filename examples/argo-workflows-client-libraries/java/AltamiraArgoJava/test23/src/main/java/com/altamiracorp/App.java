package com.altamiracorp;

import io.argoproj.workflow.ApiClient;
import io.argoproj.workflow.ApiException;
import io.argoproj.workflow.Configuration;
import io.argoproj.workflow.auth.*;
import io.argoproj.workflow.models.*;
import io.kubernetes.client.openapi.models.V1Container;
import io.kubernetes.client.openapi.models.V1ObjectMeta;
import io.argoproj.workflow.apis.WorkflowServiceApi;
import java.util.Collections; 
import java.util.ArrayList;
import java.util.List;

public class App 
{
    public static void main( String[] args )
    {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath("http://argo-workflows.pd.k8s.altamiracorp.com");
        
        WorkflowServiceApi apiInstance = new WorkflowServiceApi(defaultClient);

        ApiKeyAuth BearerToken = (ApiKeyAuth) defaultClient.getAuthentication("BearerToken");
        BearerToken.setApiKeyPrefix("Token");
        String namespace = "argo-workflows"; 

        IoArgoprojWorkflowV1alpha1WorkflowCreateRequest req = new IoArgoprojWorkflowV1alpha1WorkflowCreateRequest();
        List<String> commands = new ArrayList<>();
        commands.add("pwd");
        commands.add("ls");
        IoArgoprojWorkflowV1alpha1TTLStrategy ttl_Strategy = new IoArgoprojWorkflowV1alpha1TTLStrategy();
        ttl_Strategy.setSecondsAfterSuccess(600);
        req.setWorkflow(
           new IoArgoprojWorkflowV1alpha1Workflow()
              .metadata(new V1ObjectMeta().generateName("javasdktest-"))
              .spec(
                  new IoArgoprojWorkflowV1alpha1WorkflowSpec()
                      .entrypoint("whalesay")
                      .templates(
                          Collections.singletonList(
                            new IoArgoprojWorkflowV1alpha1Template()
                              .name("whalesay")
                              .container(
                                  new V1Container()
                                      .image("docker/whalesay:latest")
                                      .command(commands)
                              )
                          )
                      )
                  .ttlStrategy(ttl_Strategy)
              )
        );
        try {
          apiInstance.workflowServiceCreateWorkflow(namespace, req);
          System.out.println("\n Workflow Created");
        } catch (ApiException e) {
          System.err.println("Exception when calling WorkflowServiceApi#workflowServiceCreateWorkflow");
          System.err.println("Status code: " + e.getCode());
          System.err.println("Reason: " + e.getResponseBody());
          System.err.println("Response headers: " + e.getResponseHeaders());
          e.printStackTrace();
        } catch (Exception ex) {
          System.out.println("\n *** Exception ***");
          System.out.println(ex);
        }        
    }
}
