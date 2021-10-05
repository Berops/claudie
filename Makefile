.PHONY: gen contentbox scheduler builder terraformer wireguardian kubeEleven test dockerUp dockerDown dockerBuild docker

#Generate all .proto files
gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:.

contextbox:
	go run services/context-box/server/server.go

scheduler:
	go run services/scheduler/scheduler.go

builder:
	go run services/builder/builder.go

terraformer:
	go run services/terraformer/server/server.go services/terraformer/server/terraform.go

wireguardian:
	go run services/wireguardian/server/server.go

kubeEleven:
	go run services/kube-eleven/server/server.go

# -timeout 0 will disable default timeout
test:
	go test -v ./services/testing-framework/... -timeout 0

# Run all services in docker containers via docker-compose on a local machine
docker:
	DOCKER_BUILDKIT=1 docker-compose --env-file ./K8s-dev-cluster/.env up -d

dockerUp:
	docker-compose --env-file ./K8s-dev-cluster/.env up

dockerDown:
	docker-compose --env-file ./K8s-dev-cluster/.env down

dockerBuild:
	docker-compose --env-file ./K8s-dev-cluster/.env build --parallel