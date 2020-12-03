#!/bin/bash

go build -buildmode=plugin -o bin/ /home/pleimer/go/src/github.com/infrawatch/sg-core-2/plugins/transport/*
go build cmd/core.go
