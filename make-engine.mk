GIT_REPO = $(shell git remote get-url origin | sed -E 's/.*[\/:]([^\/]+\/[^\/]+)\.git$$/\1/' | sed 's/\.git$$//')
GIT_URL = $(shell git remote get-url origin)
GITHUB_USER = $(shell \
	if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then \
		gh api user --jq '.login' 2>/dev/null; \
	elif git config --get user.email >/dev/null 2>&1; then \
		curl -s "https://api.github.com/search/users?q=$$(git config --get user.email)+in:email" | grep '"login"' | head -1 | cut -d'"' -f4 2>/dev/null; \
	else \
		echo "inem"; \
	fi)
EMAIL = $(shell git config --get user.email)
GITLAB_USER = $(shell \
	if command -v glab >/dev/null 2>&1 && glab auth status >/dev/null 2>&1; then \
		glab api user --jq '.username' 2>/dev/null; \
	else \
		echo "inem"; \
	fi)
REGISTRY = $(shell \
	if echo "$(GIT_URL)" | grep -q "github.com"; then \
		echo "ghcr.io/$(GIT_REPO)"; \
	elif echo "$(GIT_URL)" | grep -q "gitlab.com"; then \
		echo "registry.gitlab.com/$(GIT_REPO)"; \
	else \
		echo "registry.example.com/$(GIT_REPO)"; \
	fi)
IMAGE_NAME = $(shell echo "$(GIT_REPO)" | sed 's/\//-/g' | sed 's/\./-/g')
IMAGE_TAG = latest

USER = "$(shell id -u):$(shell id -g)"
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

ARGS            = $(filter-out $@,$(MAKECMDGOALS))   # allows: make foo bar
LAST_MISSING    = .git/.last_missing_target

.DEFAULT:
	@if [ "$@" = "from" ] || [ "$@" = "it" ] || [ "$@" = "help" ]; then \
		: ; \
	else \
		echo "$@" > $(LAST_MISSING); \
		echo "Target '$@' not found. Try: make it"; \
		exit 1; \
	fi

it:
	@if [ -f $(LAST_MISSING) ]; then \
		target=$$(cat $(LAST_MISSING) | tr -d '\n'); \
		curl -fsSL "instll.sh/inem/makefiles" | bash -s "make $$target"; \
		rm -f $(LAST_MISSING); \
	else \
		echo "No info about last failed command"; \
	fi

it!:
	@if [ -f $(LAST_MISSING) ]; then \
		target=$$(cat $(LAST_MISSING) | tr -d '\n'); \
		EXECUTE=1 bash -c 'curl -fsSL "instll.sh/inem/makefiles" | bash -s "make '$$target'"'; \
		rm -f $(LAST_MISSING); \
	else \
		echo "No info about last failed command"; \
	fi

info:
	@echo "EMAIL: $(EMAIL)"
	@echo "GITHUB_USER: $(GITHUB_USER)"
	@echo "GITLAB_USER: $(GITLAB_USER)"
	@echo
	@echo "GIT_REPO: $(GIT_REPO)"
	@echo "GIT_URL: $(GIT_URL)"
	@echo "BRANCH: $(BRANCH)"
	@echo
	@echo "REGISTRY: $(REGISTRY)"
	@echo "IMAGE_NAME: $(IMAGE_NAME)"
