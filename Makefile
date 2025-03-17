SHELL=/bin/bash
SERVICE=accounts

.PHONY: up
up: docker
	/n/config/compose up -d ${SERVICE}

.PHONY: docker
docker:
	docker build . -t rcbilson/${SERVICE}

.PHONY: frontend
frontend:
	cd frontend && yarnpkg run build && cd -

.PHONY: backend
backend:
	cd backend && GOBIN=${PWD}/bin go install knilson.org/accounts/cmd/{query,import,update,learn,server,explain}

.PHONY: upgrade-frontend
upgrade-frontend:
	cd frontend && yarn upgrade --latest

.PHONY: upgrade-backend
upgrade-backend:
	cd backend && go get go@latest && go get -u ./...

.PHONY: upgrade
upgrade: upgrade-frontend upgrade-backend
