#!/bin/bash

usage() {
cat <<HERE
usage:
	--build-fe				installs dependencies and builds frontend

	--dev-fe				run the server in dev mode

	--build-be				builds the backend

	--build					builds entire app and stores binary in bin folder
HERE
}

build_fe() {
	if [ ! -e frontend/dist ]; then
		echo "Starting FE build process";

		# npm processes mess up on newer version of node
		export NODE_OPTIONS=--openssl-legacy-provider;

		npm --prefix frontend install;
		npm --prefix frontend run build;
	fi
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

case $1 in
	"--build-fe") build_fe;;
	"--build-be") build_be;;
	"--dev-fe") dev_fe;;
	"--build") build_project;;
	*) usage;;
esac