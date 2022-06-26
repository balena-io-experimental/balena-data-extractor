BINARY_NAME=balena-data-extractor
 
build:
	go build -o ${BINARY_NAME} -ldflags '-w -s' ./cmd/balena-data-extractor/
 
run:
	go build -o ${BINARY_NAME} -ldflags '-w -s' ./cmd/balena-data-extractor/
	./${BINARY_NAME}
