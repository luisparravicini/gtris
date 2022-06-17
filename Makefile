BINARY_NAME=bin_gtris
WEB_BINARY_NAME=gtris
WEB_DST_PATH=dist

build:
	go build -o ${BINARY_NAME} main.go

build_web: clean_web
	mkdir -p ${WEB_DST_PATH}

	GOOS=js GOARCH=wasm go build -o ${WEB_DST_PATH}/${WEB_BINARY_NAME}.wasm main.go
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js ${WEB_DST_PATH}
	cp web/index.html ${WEB_DST_PATH}

	cd ${WEB_DST_PATH} && zip -r ${WEB_BINARY_NAME}.zip * && cd -

run: build
	./${BINARY_NAME}

run_web: build_web
	cd ${WEB_DST_PATH} && python -m http.server

clean_web:
	rm -rf ${WEB_DST_PATH}

clean: clean_web
	go clean
	rm -f ${BINARY_NAME}
