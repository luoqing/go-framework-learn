#!/bin/bash

#/***************************************************************
# name: build.sh
# date: 2021-07-18
# author:
# desc:
#****************************************************************/
go test -v -test.run TestGroupCache -server localhost:8089 &
go test -v -test.run TestGroupCache -server localhost:8090 &
go test -v -test.run TestGroupCache -server localhost:8091 &
go test -v -test.run TestGroupCache -server localhost:8092 &
go test -v -test.run TestGroupCache -server localhost:8088 &
