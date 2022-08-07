.PHONY: build test clean

GO=go

build:
	@$(GO) build -o cachee

test:
	@$(GO) test -v .

clean:
	@rm cachee
