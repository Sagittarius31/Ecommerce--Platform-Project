# Cloud-Native E-Commerce Microservices Platform

Production-grade cloud-native microservices platform built with Go, Kubernetes, Terraform, and GitOps.

## Architecture
```
Client (React) -> CloudFront -> NGINX Ingress -> API Gateway
                                                     |-- User Service    -> PostgreSQL + Redis
                                                     |-- Product Service -> PostgreSQL + Redis
                                                     |-- Order Service   -> PostgreSQL + RabbitMQ
                                                     |-- Payment Service -> PostgreSQL + Stripe
                                                     `-- Notification Service (event-driven)
```

## Tech Stack
| Layer | Technology |
|---|---|
| Backend | Go 1.22, Gin, Clean Architecture |
| Databases | PostgreSQL 16 (per-service), Redis 7 |
| Messaging | RabbitMQ 3.13 |
| Inter-service | gRPC |
| Containers | Docker (multi-stage, distroless) |
| Orchestration | Kubernetes 1.30 (AWS EKS) |
| IaC | Terraform 1.7 |
| CI/CD | GitHub Actions |
| GitOps | ArgoCD |
| Monitoring | Prometheus + Grafana + Loki |
| Security | Trivy, RBAC, Network Policies, External Secrets |

## Quick Start
```bash
git clone https://github.com/YOUR_USERNAME/ecommerce-platform
cd ecommerce-platform
cp .env.example .env
make setup
make up
```

API: http://localhost:8080
RabbitMQ UI: http://localhost:15672 (admin/admin123)
Grafana: http://localhost:3000

## Commands
```bash
make help            # show all commands
make test            # run unit tests
make docker-build    # build all Docker images
make k8s-deploy      # deploy to Kubernetes
make tf-plan ENV=dev # Terraform plan
```

## License
MIT
