TARGET=gitconfig

all: $(TARGET)

gitconfig: $(shell find . -name '*.go')
	go build -o $@ cmd/gitconfig/main.go

test:
	golint && go test ./...

clean:
	rm -f $(TARGET)
