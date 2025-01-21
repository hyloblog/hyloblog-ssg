GO = go
BIN = ${CURDIR}/bin
SOURCES := $(shell find $(CURDIR) -name '*.go')

SSG = $(BIN)/hylodoc-ssg

$(SSG): $(BIN) $(SOURCES)
	@printf 'GO\t$@\n'
	@$(GO) build -o $@

$(BIN):
	@mkdir -p $(BIN)

clean:
	@rm -rf $(BIN)
