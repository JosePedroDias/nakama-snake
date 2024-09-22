#!/bin/bash

go build -buildmode=plugin -trimpath -o ./modules/snake.so
cp modules/snake.so ~/Personal/nakama/data

