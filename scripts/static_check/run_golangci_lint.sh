#!/bin/bash
export GOPROXY="http://172.16.20.30:32888,https://goproxy.cn,direct"
export CGO_ENABLED="1"

pwd

res=0
go mod tidy
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./api/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./cmd/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./dao/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./models/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./pkg/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./router/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./service/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"

golangci-lint run ./job/...
r=$?
res=`expr ${res} + ${r}`
echo "$res"
if [ $res -ne 0 ]; then
    exit 1
fi
