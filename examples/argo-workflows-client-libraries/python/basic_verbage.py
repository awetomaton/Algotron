from argo_workflows.apis import WorkflowServiceApi
from argo_workflows.models import IoArgoprojWorkflowV1alpha1WorkflowList
from hera.task import Task
from hera.workflow import Workflow
from hera.workflow_service import WorkflowService
from hera.config import Config
from hera.client import Client
import base64
import errno
import os
from typing import Optional
from time import sleep

from kubernetes import client, config


def get_sa_token(service_account: str, namespace: str = "default", config_file: Optional[str] = None):
    """Get ServiceAccount token using kubernetes config.

     Parameters
    ----------
    service_account: str
        The service account to authenticate from.
    namespace: str = 'default'
        The K8S namespace the workflow service submits workflows to. This defaults to the `default` namespace.
    config_file: Optional[str] = None
        The path to k8s configuration file.

     Raises
    ------
    FileNotFoundError
        When the config_file can not be found.
    """
    if config_file is not None and not os.path.isfile(config_file):
        raise FileNotFoundError(
            errno.ENOENT, os.strerror(errno.ENOENT), config_file)

    config.load_kube_config(config_file=config_file)
    v1 = client.CoreV1Api()
    secret_name = v1.read_namespaced_service_account(
        service_account, namespace).secrets[0].name
    sec = v1.read_namespaced_secret(secret_name, namespace).data
    return base64.b64decode(sec["token"]).decode()


def hello():
    print("Hello, Hera!")


# replace namespace with your namespace name
# replace user with username in the quotes
_namespace = "argo-workflows"
_host = "127.0.0.1:8080"
_token = get_sa_token("user", namespace=_namespace)
apiclient = Client(Config(host=_host,
                   verify_ssl=False), token=_token).api_client

# print(_token)

wsa = WorkflowServiceApi(api_client=apiclient)

# host is your argo server address
# http://argo-workflows.pd.k8s.altamiracorp.com

with Workflow(
    "k8s-sa", service=WorkflowService(host=_host, token=_token, namespace=_namespace)
) as w:
    Task("t", hello)

with Workflow(
    "k8s-sb", service=WorkflowService(host=_host, token=_token, namespace=_namespace)
) as ws:
    Task("t", hello)

w.create()
ws.create()

sleep(15)

print(wsa.get_workflow(_namespace, "k8s-sb").metadata)

print(wsa.list_workflows(_namespace, async_req=False))

wsa.delete_workflow(_namespace, "k8s-sa")
wsa.delete_workflow(_namespace, "k8s-sb")
