#!/bin/sh

download_fe() {
	if [ ! -e frontend ]; then
		echo "FE doesn't exist, cloning";
		git clone https://github.com/ikeohachidi/chapi-fe;
		mv chapi-fe frontend;
		npm --prefix frontend install;
	fi
}

build_fe() {
	if [ ! -e frontend/dist ]; then
		echo "Starting FE build process";
		npm --prefix frontend install;
		npm --prefix frontend run build;
	fi
}

download_fe;
build_fe;