GITIGNORE_REPO := https://github.com/github/gitignore.git
GITIGNORE_DIR := gitignore
SHELL := $(shell basename $$SHELL)

.PHONY: all
all: build install_completions

.PHONY: copy_gitignore
copy_gitignore:
	@if [ ! -d "$(GITIGNORE_DIR)" ]; then \
		git clone --depth=1 $(GITIGNORE_REPO) $(GITIGNORE_DIR); \
	fi
	@cd ${GITIGNORE_DIR} && \
	for f in $$(find -H -type f -name '*.gitignore'); do \
        cp -vur "$$f" "../src/gitignore/$$(basename $$f | tr '[A-Z]' '[a-z]')"; \
    done

.PHONY: build
build: copy_gitignore
	go build -ldflags "-s -w" -o gen

.PHONY: install_completions
install_completions: build
	@gen completion $(SHELL) > gen.$(SHELL)
	@cp gen.$(SHELL) $(HOME)/.local/share/completions/gen.$(SHELL)

.PHONY: install
install: install_completions
	@cp gen $(HOME)/bin/gen

.PHONY: clean
clean:
		@rm -fr ${GITIGNORE_DIR} gen.${SHELL} gen src/gitignore/*.gitignore
