SERVICES := user-service product-service order-service payment-service notification-service api-gateway
ECR ?= YOUR_ECR_URL
TAG ?= latest

up:
	docker compose up -d

down:
	docker compose down -v

test:
	@for s in $(SERVICES); do cd services/$$s && go test ./... -race && cd ../..; done

build:
	@for s in $(SERVICES); do docker build -t $(ECR)/$$s:$(TAG) services/$$s; done

push:
	@for s in $(SERVICES); do docker push $(ECR)/$$s:$(TAG); done

k8s:
	kubectl apply -f infrastructure/kubernetes/namespaces/
	kubectl apply -f infrastructure/kubernetes/storage/
	kubectl apply -f infrastructure/kubernetes/statefulsets/
	kubectl apply -f infrastructure/kubernetes/rbac/
	kubectl apply -f infrastructure/kubernetes/security/
	kubectl apply -f infrastructure/kubernetes/services/
	kubectl apply -f infrastructure/kubernetes/ingress/

status:
	kubectl get pods -n ecommerce -o wide
