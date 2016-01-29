#!/usr/bin/env bash
set -e

VERSION=$1
PROJECT=$2
QUAYUSERNAME=$3

if [ -z "${VERSION}" ] || [ -z "${PROJECT}" ] || [ -z "${QUAYUSERNAME}" ]; then
	echo "Usage: ${0} VERSION PROJECT QUAYUSERNAME" >> /dev/stderr
	exit 255
fi

if ! command -v docker >/dev/null; then
    echo "cannot find docker"
    exit 1
fi

LXY_ROOT=$(dirname "${BASH_SOURCE}")/..

echo $LXY_ROOT

pushd ${LXY_ROOT} >/dev/null
	#echo Building lxy binary...
	#./scripts/build-binary ${VERSION}
	echo Building docker image...
	LXYDIR=$LXY_ROOT ./scripts/build-docker ${VERSION}
	docker tag quay.io/${QUAYUSERNAME}/lxy:${VERSION} gcr.io/${PROJECT}/lxy:${VERSION}
	gcloud docker push gcr.io/${PROJECT}/lxy:${VERSION}
popd >/dev/null
