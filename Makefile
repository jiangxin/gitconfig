TARGET=gitconfig

define message
	@echo "### $(1)"
endef

all: $(TARGET)

gitconfig: $(shell find . -name '*.go')
	$(call message,Build $@)
	go build -o $@ cmd/gitconfig/main.go

test: $(TARGET) golint
	$(call message,Testing gitconfig using golint for coding style)
	@golint
	$(call message,Testing gitconfig for unit tests)
	@go test

golint:
	@if ! type golint >/dev/null 2>&1; then \
		go get golang.org/x/lint/golint; \
	fi

clean:
	rm -f $(TARGET)

.PHONY: test clean golint
