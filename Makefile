NAME=Subscription-bot

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

#.SILENT:

define colored
	@echo '${GREEN}$1${RESET}'
endef

## dependencies - fetch all dependencies for scripts
dependencies:
	${call colored,dependensies is running...}
	./scripts/get-dependencies.sh

## lint project
lint:
	${call colored,lint is running...}
	./scripts/linters.sh
.PHONY: lint

## ------------------------------------------------- Common commands: --------------------------------------------------
## Formats the code.
format:
	${call colored,formatting is running...}
	go vet ./...
	go fmt ./...

## Fix-imports order.
fix-imports:
	${call colored,fixing imports...}
	./scripts/fix-imports-order.sh
