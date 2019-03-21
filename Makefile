TARGET=gitconfig

define message
	@echo "### $(1)"
endef

all: $(TARGET)

gitconfig: $(shell find . -name '*.go')
	go build -o $@ cmd/gitconfig/main.go

test:
	$(call message,Testing goconfig using golint for coding style)
	@golint
	$(call message,Testing goconfig for unit tests)
	@go test ./...

clean:
	rm -f $(TARGET)
