#!/bin/bash

GOFLAGS=-mod=mod go generate ./...
kratos run