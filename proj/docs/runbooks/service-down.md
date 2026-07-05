# Runbook: Service Down
## Steps
```bash
kubectl get pods -n ecommerce -l app=<service>
kubectl logs -n ecommerce <pod> --previous
kubectl rollout undo deployment/<service> -n ecommerce
```
