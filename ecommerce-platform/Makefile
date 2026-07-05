.PHONY: help setup up down logs test docker-build docker-push k8s-deploy proto swagger clean

SERVICES := user-service product-service order-service payment-service notification-service api-gateway
IMAGE_REGISTRY ?= YOUR_AWS_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/ecommerce
IMAGE_TAG      ?= latest
ENV            ?= dev
SVC            ?= user-service

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-25s %s\n", $$1, $$2}'

setup: ## One-command local setup
	chmod +x scripts/setup-local.sh && ./scripts/setup-local.sh

up: ## Start all services
	docker compose up -d

down: ## Stop all services
	docker compose down -v

logs: ## Tail logs (usage: make logs SVC=user-service)
	docker compose logs -f $(SVC)

test: ## Run unit tests for all services
	@for svc in $(SERVICES); do \
		echo "Testing $$svc..."; \
		cd services/$$svc && go test ./... -v -race -count=1 2>&1 | tail -5 && cd ../..; \
	done

docker-build: ## Build all Docker images
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		docker build -t $(IMAGE_REGISTRY)/$$svc:$(IMAGE_TAG) services/$$svc; \
	done

docker-push: ## Push images to ECR
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(IMAGE_REGISTRY)
	@for svc in $(SERVICES); do docker push $(IMAGE_REGISTRY)/$$svc:$(IMAGE_TAG); done

docker-scan: ## Trivy scan all images
	@for svc in $(SERVICES); do trivy image --severity HIGH,CRITICAL $(IMAGE_REGISTRY)/$$svc:$(IMAGE_TAG); done

k8s-deploy: ## Deploy all Kubernetes manifests
	kubectl apply -f infrastructure/kubernetes/namespaces/
	kubectl apply -f infrastructure/kubernetes/storage/
	kubectl apply -f infrastructure/kubernetes/statefulsets/
	kubectl apply -f infrastructure/kubernetes/rbac/
	kubectl apply -f infrastructure/kubernetes/security/
	kubectl apply -f infrastructure/kubernetes/pdb/
	kubectl apply -f infrastructure/kubernetes/services/
	kubectl apply -f infrastructure/kubernetes/ingress/
	kubectl apply -f infrastructure/kubernetes/monitoring/

k8s-status: ## Show pod status
	kubectl get pods -n ecommerce -o wide

k8s-logs: ## Stream logs (usage: make k8s-logs SVC=user-service)
	kubectl logs -n ecommerce -l app=$(SVC) -f --tail=100

k8s-restart: ## Rolling restart (usage: make k8s-restart SVC=user-service)
	kubectl rollout restart deployment/$(SVC) -n ecommerce

tf-init: ## Terraform init (usage: make tf-init ENV=dev)
	cd infrastructure/terraform && terraform init -backend-config="key=ecommerce/$(ENV)/terraform.tfstate"

tf-plan: ## Terraform plan
	cd infrastructure/terraform && terraform plan -var-file=environments/$(ENV)/terraform.tfvars

tf-apply: ## Terraform apply
	cd infrastructure/terraform && terraform apply -var-file=environments/$(ENV)/terraform.tfvars -auto-approve

proto: ## Generate gRPC code
	chmod +x scripts/generate-proto.sh && ./scripts/generate-proto.sh

clean: ## Remove build artifacts
	@for svc in $(SERVICES); do rm -f services/$$svc/server services/$$svc/coverage.out; done
