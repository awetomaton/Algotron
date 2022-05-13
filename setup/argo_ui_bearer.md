How to retrieve bearer token

```
SECRET=$(kubectl get sa choochoo -n workflows -o=jsonpath='{.secrets[0].name}')

ARGO_TOKEN="Bearer $(kubectl get secret $SECRET -n workflows -o=jsonpath='{.data.token}' | base64 -d)"

echo $ARGO_TOKEN
```