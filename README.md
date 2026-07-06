# Cloud-Native E-Commerce Microservices Platform

[![CI](https://github.com/Sagittarius31/Ecommerce--Platform-Project/actions/workflows/ci.yml/badge.svg)](https://github.com/Sagittarius31/Ecommerce--Platform-Project/actions)
[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://go.dev/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.30-326CE5.svg)](https://kubernetes.io/)
[![Terraform](https://img.shields.io/badge/Terraform-1.7-7B42BC.svg)](https://www.terraform.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A production-grade cloud-native e-commerce platform built with **6 Go microservices**, deployed on **AWS EKS** using **Kubernetes**, provisioned by **Terraform**, and delivered via **ArgoCD GitOps**. Includes full observability, security hardening, and automated CI/CD.

---

## Architecture

```
Client (Browser / Mobile)
         │
         ▼
  NGINX Ingress (TLS · Rate Limiting)
         │
         ▼
    API Gateway (:8080)
    JWT Validation · CORS · Reverse Proxy
         │
   ┌─────┼──────┬──────────┬──────────────┐
   ▼     ▼      ▼          ▼              ▼
 User  Product  Order    Payment      Notification
 Svc    Svc     Svc       Svc           Svc
:8081  :8082   :8083     :8084          :8085
  │      │       │  ──gRPC──►  │              │
  │      │       │             │              │
  ▼      ▼       ▼      ▼      ▼              ▼
users  products orders redis payments   [consumer]
 _db    _db    _db  [cache]  _db
                       │
                       ▼
                  RabbitMQ 3.13
              (topic exchange · DLQ)
              order.placed ──────────► Payment Svc
              order.placed ──────────► Notification Svc
              payment.succeeded ──────► Notification Svc
```

---

## Tech Stack

| Category | Technology |
|---|---|
| **Language** | Go 1.22 |
| **HTTP Framework** | Gin |
| **Architecture** | Clean Architecture (domain / service / handler / repository) |
| **Databases** | PostgreSQL 16 (per-service), Redis 7 (cache-aside) |
| **Messaging** | RabbitMQ 3.13 (topic exchange, dead-letter queues) |
| **Inter-service** | gRPC + Protocol Buffers |
| **Payments** | Stripe API (PaymentIntent + webhook verification) |
| **Containers** | Docker (multi-stage builds, distroless images ~15MB) |
| **Orchestration** | Kubernetes 1.30 on AWS EKS |
| **IaC** | Terraform 1.7 (VPC, EKS, RDS, ECR, ElastiCache) |
| **CI/CD** | GitHub Actions (path-filtered matrix builds, Trivy scanning) |
| **GitOps** | ArgoCD (automated sync, drift detection, self-healing) |
| **Monitoring** | Prometheus + Grafana + Loki + Alertmanager |
| **Security** | RBAC, Network Policies, IRSA, External Secrets, cert-manager |

---

## Project Structure

```
ecommerce-platform/
├── services/
│   ├── user-service/          # Auth, JWT, profiles
│   ├── product-service/       # Catalog, search, Redis cache
│   ├── order-service/         # Checkout, RabbitMQ publisher, gRPC client
│   ├── payment-service/       # Stripe integration, webhooks
│   ├── notification-service/  # RabbitMQ consumer, SMTP email
│   └── api-gateway/           # Reverse proxy, rate limiting, CORS
│
├── infrastructure/
│   ├── kubernetes/            # All K8s manifests
│   │   ├── namespaces/
│   │   ├── statefulsets/      # PostgreSQL x4, Redis, RabbitMQ
│   │   ├── services/          # Deployments + HPA + Services x6
│   │   ├── rbac/              # ServiceAccounts, Roles
│   │   ├── security/          # NetworkPolicies
│   │   ├── ingress/           # NGINX Ingress + TLS
│   │   ├── monitoring/        # Prometheus alerts
│   │   └── argocd/            # ArgoCD Application manifest
│   │
│   ├── terraform/             # AWS infrastructure as code
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── modules/
│   │       ├── vpc/           # VPC, subnets, NAT Gateway
│   │       ├── eks/           # EKS cluster, node groups, IRSA
│   │       ├── rds/           # PostgreSQL RDS Multi-AZ
│   │       └── ecr/           # ECR repositories x6
│   │
│   └── helm/                  # Helm chart for all services
│
├── proto/                     # gRPC proto definitions
├── scripts/                   # Setup, seed, proto generation
├── docs/                      # ADRs, runbooks
├── .github/workflows/         # CI/CD pipelines
├── docker-compose.yml         # Full local stack
├── Makefile                   # Common commands
└── .env.example               # Environment variable template
```

---

## Quick Start (Local)

### Prerequisites

- Go 1.22+
- Docker Desktop
- Git

### Run locally with Docker Compose

```bash
# Clone the repository
git clone https://github.com/Sagittarius31/Ecommerce--Platform-Project.git
cd ecommerce-platform

# Copy environment variables
cp .env.example .env

# Start the full stack (all 6 services + databases)
docker compose up -d

# Check all services are healthy
docker compose ps
```

### Test the API

```bash
# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123!",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login and get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "Password123!"}'

# Get your profile (use token from login response)
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"

# List products
curl http://localhost:8080/api/v1/products
```

### Health checks

```bash
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Product Service
curl http://localhost:8083/health  # Order Service
curl http://localhost:8084/health  # Payment Service
curl http://localhost:8085/health  # Notification Service
```

---

## Service Endpoints

### API Gateway (public entry point — :8080)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | /api/v1/auth/register | No | Register new user |
| POST | /api/v1/auth/login | No | Login, returns JWT |
| GET | /api/v1/products | No | List products |
| GET | /api/v1/products/:id | No | Get product by ID |
| GET | /api/v1/users/me | Yes | Get current user profile |
| PUT | /api/v1/users/me | Yes | Update profile |
| POST | /api/v1/orders | Yes | Create order |
| GET | /api/v1/orders/:id | Yes | Get order by ID |
| POST | /api/v1/payments | Yes | Create payment intent |
| POST | /webhooks/stripe | No | Stripe webhook (signature verified) |

---

## Kubernetes Deployment

### Prerequisites
- kubectl configured for your cluster
- AWS CLI configured
- Helm 3.x

### Deploy to Kubernetes

```bash
# Apply manifests in correct order
kubectl apply -f infrastructure/kubernetes/namespaces/
kubectl apply -f infrastructure/kubernetes/storage/
kubectl apply -f infrastructure/kubernetes/statefulsets/

# Wait for databases to be ready
kubectl wait pod -l app=postgres-user -n ecommerce \
  --for=condition=ready --timeout=120s

kubectl apply -f infrastructure/kubernetes/rbac/
kubectl apply -f infrastructure/kubernetes/security/
kubectl apply -f infrastructure/kubernetes/pdb/
kubectl apply -f infrastructure/kubernetes/services/
kubectl apply -f infrastructure/kubernetes/ingress/

# Or use make
make k8s
```

### Check status

```bash
kubectl get pods -n ecommerce
kubectl get hpa -n ecommerce
kubectl get ingress -n ecommerce
```

---

## AWS Infrastructure (Terraform)

```bash
# Configure AWS credentials
aws configure

# Create S3 bucket for state (replace with unique name)
aws s3 mb s3://YOUR-TERRAFORM-STATE-BUCKET --region us-east-1

# Update bucket name in infrastructure/terraform/main.tf

# Copy and fill in your variables
cp infrastructure/terraform/environments/dev/terraform.tfvars.example \
   infrastructure/terraform/environments/dev/terraform.tfvars

# Set sensitive values as env vars (never in tfvars)
export TF_VAR_db_password="YourStrongPassword123!"

# Initialize and apply (creates VPC, EKS, RDS, ECR)
cd infrastructure/terraform
terraform init
terraform plan -var-file=environments/dev/terraform.tfvars
terraform apply -var-file=environments/dev/terraform.tfvars

# Configure kubectl for new EKS cluster
aws eks update-kubeconfig --region us-east-1 --name ecommerce-eks-dev
```

---

## CI/CD Pipeline

### GitHub Actions (CI)

Every push triggers:

1. **Change detection** — only rebuilds services that changed
2. **Unit tests** — `go test ./... -race`
3. **Docker build** — multi-stage, distroless final image
4. **Trivy scan** — blocks deployment on CRITICAL CVEs
5. **ECR push** — tagged with git SHA + latest

### ArgoCD (CD / GitOps)

- Git is the single source of truth
- ArgoCD watches the `infrastructure/kubernetes/` directory
- Every merge to `main` automatically syncs to EKS
- Drift detection reverts any manual `kubectl` changes

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd \
  -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Apply the application manifest
kubectl apply -f infrastructure/kubernetes/argocd/applications/ecommerce.yaml
```

---

## Monitoring

### Install monitoring stack

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Prometheus + Grafana + Alertmanager
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace \
  --set grafana.adminPassword=YourGrafanaPassword

# Loki log aggregation
helm install loki grafana/loki-stack \
  --namespace monitoring --set grafana.enabled=false
```

### Access dashboards

```bash
# Grafana
kubectl port-forward svc/monitoring-grafana 3000:80 -n monitoring
# Open http://localhost:3000 — admin / YourGrafanaPassword

# Prometheus
kubectl port-forward svc/prometheus-operated 9090:9090 -n monitoring
# Open http://localhost:9090
```

### Key metrics

```promql
# Request rate per service
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Error rate percentage
rate(http_requests_total{status=~"5.."}[5m])
  / rate(http_requests_total[5m]) * 100
```

---

## Security

| Layer | Control | Purpose |
|---|---|---|
| Container | Distroless image, non-root user | No shell, minimal attack surface |
| Container | Trivy CVE scanning | Blocks CRITICAL vulnerabilities |
| Kubernetes | RBAC (least-privilege) | Pods only access what they need |
| Kubernetes | Default-deny NetworkPolicy | Prevents lateral movement |
| Kubernetes | PodDisruptionBudget | Zero-downtime node drains |
| AWS | IRSA (no long-lived credentials) | Pods assume IAM roles via OIDC |
| AWS | External Secrets Operator | Secrets pulled from Secrets Manager |
| Application | JWT (HMAC-SHA256) | Stateless, scalable authentication |
| Application | bcrypt cost=12 | Brute-force resistant password hashing |
| Application | Parameterized SQL | SQL injection prevention |
| Application | Stripe webhook signature | Prevents fake payment events |
| TLS | cert-manager + Let's Encrypt | Encrypted external traffic |

---

## Makefile Commands

```bash
make up           # Start all services with Docker Compose
make down         # Stop all services
make test         # Run unit tests for all services
make build        # Build all Docker images
make push         # Push images to ECR
make k8s          # Deploy to Kubernetes
make status       # Show pod status
```

---

## Architecture Decisions

| ADR | Decision | Reason |
|---|---|---|
| [ADR-001](docs/architecture/ADR-001-microservices.md) | Microservices | Independent scaling, fault isolation, team autonomy |
| [ADR-002](docs/architecture/ADR-002-database-per-service.md) | Database per service | Schema independence, no cross-service coupling |

---

## Key Design Patterns

- **Clean Architecture** — domain → service → handler → repository (dependencies point inward)
- **Cache-aside** — Redis caches products; miss falls through to PostgreSQL
- **Database-per-service** — each service owns its PostgreSQL database; no shared schemas
- **Event-driven** — RabbitMQ decouples Order, Payment, and Notification services
- **GitOps** — ArgoCD reconciles cluster state to Git continuously
- **IRSA** — pods use AWS IAM roles via service account tokens, no static credentials

---

## What This Demonstrates

This project was built to demonstrate production-grade engineering skills:

- **Go microservices** with Clean Architecture
- **Kubernetes** — HPA, StatefulSets, RBAC, NetworkPolicies, PDB, Ingress
- **Terraform** — modular AWS infrastructure as code
- **GitOps** — ArgoCD with drift detection and self-healing
- **CI/CD** — path-filtered pipelines with security scanning
- **Observability** — the four golden signals (latency, traffic, errors, saturation)
- **Security** — defence-in-depth across every layer

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Author

**Pramod Patil**
- GitHub: [@Sagittarius31](https://github.com/Sagittarius31)
- LinkedIn: [linkedin.com/in/pramod-patil-907003211](https://linkedin.com/in/pramod-patil-907003211)
