check:
	@echo "-> go fmt"
	@test "$(shell gofmt -l . | wc -l)" -eq "0"
	@echo "ok"
	@echo "-> goimports"
	@test "$(shell goimports -l . | wc -l)" -eq "0"
	@echo "ok"
	@echo "-> gometalinter"
	@gometalinter --vendor ./...
	@echo "ok"

lint:
	go fmt ./...
	goimports -w .

test:
	go test -v -race -covermode=atomic -coverprofile=profile.cov
