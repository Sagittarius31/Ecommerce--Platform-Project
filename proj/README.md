# Cloud-Native E-Commerce Microservices Platform

Production-grade platform built with Go, Kubernetes, Terraform, and GitOps.

## Architecture
```
Client → NGINX Ingress → API Gateway
                              ├── User Service    → PostgreSQL + Redis
                              ├── Product Service → PostgreSQL + Redis
                              ├── Order Service   → PostgreSQL + RabbitMQ
                              ├── Payment Service → PostgreSQL + Stripe
                              └── Notification Service (event-driven)
```

## Tech Stack
| Layer | Technology |
|---|---|
| Backend | Go 1.22, Gin, Clean Architecture |
| Databases | PostgreSQL 16 (per-service), Redis 7 |
| Messaging | RabbitMQ 3.13, gRPC |
| Containers | Docker (multi-stage, distroless) |
| Orchestration | Kubernetes 1.30 (AWS EKS) |
| IaC | Terraform 1.7 |
| CI/CD | GitHub Actions + ArgoCD |
| Monitoring | Prometheus + Grafana + Loki |

## Quick Start
```bash
git clone https://github.com/YOUR_USERNAME/ecommerce-platform
cd ecommerce-platform
cp .env.example .env
docker compose up -d
```
