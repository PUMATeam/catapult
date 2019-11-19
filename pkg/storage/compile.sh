#!/bin/bash
protoc storage.proto --go_out=plugins=grpc:.
