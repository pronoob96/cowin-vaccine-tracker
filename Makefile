GO=go
SOURCE=main.go
TARGET=main

help:
	@echo "use \`make <target> where <target> is:\'"
	@echo "	help: display this help"
	@echo "	clean: clean the build directory"
	@echo "	build: generate executable"
	@echo "	run: clean, build and execute code"

build:
	${GO} build -o ${TARGET} ${SOURCE}

run: clean build
	./${TARGET}

clean:
	rm -f ${TARGET} 1>/dev/null 2>&1
