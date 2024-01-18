docker:
	sudo DOCKER_BUILDKIT=1 docker build -t account_server .

frontend:
	cd src/frontend && yarnpkg run build && cd -

backend:
	cd src && GOBIN=${PWD}/bin go install knilson.org/accounts/cmd/{query,import,update,learn,server} && cd -
