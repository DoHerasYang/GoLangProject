- `kubectl get deployments --all-namespaces`
- `kubectl get pods`
- `kubectl delete deployment <deployment-name> --cascade=background`

```shell
kubectl apply -f mysql-deployment.yaml & kubectl apply -f mysql-service.yaml & kubectl apply -f mysql-pvc.yaml  & kubectl apply -f mysql-pv.yaml

kubectl apply -f webook-deployment.yaml & kubectl apply -f webook-service.yaml
```

---