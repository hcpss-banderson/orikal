#!/usr/bin/env bash

archs=(amd64 arm64)
plats=(linux darwin)

for arch in ${archs[@]}; do
  for plat in ${plats[@]}; do
	  env GOOS=${plat} GOARCH=${arch} go build -o orikal_${plat}_${arch}
	done
done
