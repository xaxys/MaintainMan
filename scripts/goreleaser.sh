#!/bin/bash
docker run --rm --privileged -e CGO_ENABLED=1 -e GOPROXY=https://goproxy.cn,direct -e GOVERSION=$(go version | awk '{print $3;}') -v /var/run/docker.sock:/var/run/docker.sock -v `pwd`:/go/src/maintainman -v /usr:/sysroot/usr -w /go/src/maintainman xaxy/goreleaser-cross:v1.18.0 release --snapshot --rm-dist