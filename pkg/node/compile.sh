#!/bin/bash
protoc node.proto --go_out=plugins=grpc:.
