#!/usr/bin/env bash

docker build -t fint-model --build-arg VERSION=0.$(date +%y%m%d.%H%M) .
