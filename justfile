set shell := ["bash", "-cu"]

BUILD_DIR := "build"
EXECUTABLE := "project_name"

clean:
	rm -rf {{BUILD_DIR}}

build:
	mkdir -p {{BUILD_DIR}}
	cd {{BUILD_DIR}} && cmake -G Ninja .. && ninja

rebuild:
	just clean
	just build

run args="":
	./{{BUILD_DIR}}/{{EXECUTABLE}} {{ args }}

test:
	just build
	cd {{BUILD_DIR}} && ctest --output-on-failure
