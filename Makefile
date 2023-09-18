
BINARY=./bin/compare

build:
	go build -o ${BINARY} compare.go types.go main.go

run: build
	${BINARY} -bucket-a=simonetes-bucket-a -bucket-b=simonetes-bucket-b

clean:
	go clean
	rm ${BINARY}
