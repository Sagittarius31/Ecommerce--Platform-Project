# Runbook: Service Down

**Severity:** P1 | **Alert:** ServiceDown

## Immediate Steps (0-2 minutes)
```bash
kubectl get pods -n ecommerce -l app=<service-name> -o wide
kubectl logs -n ecommerce <pod-name> --previous
kubectl describe pod -n ecommerce <pod-name>
```

## Common Causes
| Symptom | Fix |
|---|---|
| CrashLoopBackOff | Read --previous logs for panic |
| OOMKilled | Increase memory limits |
| ImagePullBackOff | Check ECR URL and IAM |
| Pending | Check node capacity |

## Rollback
```bash
kubectl rollout undo deployment/<service-name> -n ecommerce
# OR GitOps rollback:
git revert HEAD && git push
```
