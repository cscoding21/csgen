cov:
	go vet ./...&& \
	go test -coverprofile=c.out ./... && \
	go tool cover -html=c.out;

qual:
	go fmt -x  ./... && \
	golint . && \
	go vet ./... && \
	go test ./...;
