#!/bin/bash -eux

export buildArch="-X gitlab.com/fenrirunbound/pipeline-queue/internal.buildArch"

systems=(darwin linux windows)

for goos in ${systems[@]}; do
    GOOS=$goos go build -ldflags "$FLAGS $buildArch=${goos}-amd64" \
        -o ./target/${goos}/pipeline-queue
done
echo $VERSION > ./target/version

