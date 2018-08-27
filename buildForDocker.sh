#!/bin/bash
docker run --entrypoint /go/build.sh --mount src="$(pwd)",target=/go,type=bind -it hashicorp/terraform:full
