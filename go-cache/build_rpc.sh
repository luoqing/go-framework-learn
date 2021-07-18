#!/bin/bash

#/***************************************************************
# name: build.sh
# date: 2021-07-18
# author:
# desc:
#****************************************************************/
go build -o ./gee-cache

./gee-cache -rpc_server localhost:8088 -api_server localhost:8188 -rpc_sinker file &
./gee-cache -rpc_server localhost:8089 -api_server localhost:8189 -rpc_sinker file &
./gee-cache -rpc_server localhost:8090 -api_server localhost:8190 -rpc_sinker file &
./gee-cache -rpc_server localhost:8091 -api_server localhost:8191 -rpc_sinker file &
./gee-cache -rpc_server localhost:8092 -api_server localhost:8192 -rpc_sinker file &
