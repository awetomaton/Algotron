helm repo add argo https://argoproj.github.io/argo-helm
helm repo add bitnami https://charts.bitnami.com/bitnami

kubectl create namespace workflows

helm install argo-workflows argo/argo-workflows -n workflows --set server.extraArgs={--auth-mode=server},server.serviceType="NodePort",server.serviceNodePort=30010
helm install argo-events argo/argo-events -n workflows
helm install argo-artifacts bitnami/minio -n workflows --set fullnameOverride=argo-artifacts --set auth.rootUser=admin,auth.rootPassword=adminpass --set service.type=NodePort,service.nodePorts.api=30011,service.nodePorts.console=30012

# https://github.com/argoproj/argo-workflows/blob/master/docs/service-accounts.md - grant admin privileges
# Manually patch workflow server with info in rancher-desktop-artifact-repository-configmap.yaml