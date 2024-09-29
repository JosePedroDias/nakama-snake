#!/bin/bash

rm -rf modules
docker run --rm -w "/builder" --platform linux/amd64 -v "${PWD}:/builder" heroiclabs/nakama-pluginbuilder:3.22.0 build -buildvcs=false -buildmode=plugin -trimpath -o ./modules/snake.so
cp modules/snake.so ~/Personal/nakama/data
