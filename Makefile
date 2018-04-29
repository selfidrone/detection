build_base:
	docker build --no-cache -t nicholasjackson/go-opencv:1.10.1 -f Dockerfile.base .

build_linux:
	docker build -t detectionbuild -f Dockerfile.build .
	docker run -it --rm -v $(shell pwd):/go/src/github.com/selfidrone/detection detectionbuild go build -o detection.linux main.go

test_build:
	docker run -it --rm -v $(shell pwd):/go/src/github.com/selfidrone/detection detectionbuild /bin/sh

build_docker:
	docker build -t quay.io/selfidrone/detection:latest .

run_docker:
	docker run --rm -it -p 9999:9999 quay.io/selfidrone/detection:latest
