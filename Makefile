cov:
	go vet ./...&& \
	go test -coverprofile=c.out ./... && \
	go tool cover -html=c.out;

test:
	go fmt -x  ./... && \
	golint . && \
	go vet ./... && \
	go test ./...;
	
qual:
	go vet . && \
	golint .;