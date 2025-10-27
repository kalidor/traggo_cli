cover:
	mkdir cover
	go test ./tests/... -coverpkg ./... -coverprofile=cover/c.out
	go tool cover -html=cover/c.out -o cover/coverage.html
