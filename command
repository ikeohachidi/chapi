#!/bin/bash

usage() {
cat <<HERE
usage:
	--build-fe				installs dependencies and builds frontend

	--dev-fe				run the server in dev mode

	--build-be				builds the backend

	--build					builds entire app and stores binary in bin folder

	--docker-build			builds the docker image for chapi

	--docker-run			runs the docker image for chapi, NOTE: runs --docker-build first
HERE
}

build_fe() {
	echo "Starting FE build process";

	# npm processes mess up on newer version of node
	export NODE_OPTIONS=--openssl-legacy-provider;

	npm --prefix frontend install;
	npm --prefix frontend run build;
}

dev_fe() {
	# npm processes mess up on newer version of node
	export NODE_OPTIONS=--openssl-legacy-provider;

	npm --prefix frontend run serve;
}

build_be() {
	echo "Starting BE build process";
	go build -o bin/chapi;
}

build_project() {
	echo "Building project";
	build_fe;
	build_be;
}

docker_build() {
	docker image rm -f chapi;
	docker image build -t chapi .;
}

docker_run() {
	docker_build;
	docker container rm -f chapi;

	docker run --detach \
	--env=PORT=${PORT} \
	--env=LOCAL_FRONTEND=${LOCAL_FRONTEND} \
	--env=PSQL_PASS=${PSQL_PASS} \
	--env=PG_CRYPT_KEY=${PG_CRYPT_KEY} \
	--env=SESSION_KEY=${SESSION_KEY} \
	--env=ENV=${ENV} \
	--env=CHAPI_SERVER_URL=${CHAPI_SERVER_URL} \
	--env=CHAPI_GITHUB_SECRET=${CHAPI_GITHUB_SECRET} \
	--env=CHAPI_GITHUB_ID=${CHAPI_GITHUB_ID} \
	--env=OAUTH_STATE=${OAUTH_STATE} \
	--publish=5000:5000 \
	--network=host \
	--restart=unless-stopped \
	--name=chapi \
	chapi;

	docker system prune -f;
}

case $1 in
	"--build-fe") build_fe;;
	"--build-be") build_be;;
	"--dev-fe") dev_fe;;
	"--build") build_project;;
	"--docker-build") docker_build;;
	"--docker-run") docker_run;;
	*) usage;;
esac