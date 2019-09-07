@echo off

go get github.com/seamia/tools
pushd .
cd %GOPATH%/src/github.com/seamia/tools/assets/cmd/assets
go install
popd 
assets.exe -output ${GOPATH}/src/github.com/seamia/tools/assets/example/assets.go -package main -root ${GOPATH}/src/github.com/seamia/tools/assets/example/assets -src files.list -header header.txt
