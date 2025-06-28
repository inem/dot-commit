include make-*.mk

logs:
	tail .git/commit-msg-debug.log

hook:
	cat .git/hooks/commit-msg

install:
	go mod init dotdotdot && go get github.com/sashabaranov/go-openai

install-hook:
	curl -fsSL https://raw.githubusercontent.com/inem/dotdotdot/main/install.sh | sh

build:
	go build -o release/commit-msg-go commit-msg.go

release:
	git tag v0.0.$(ARGS)
	git push origin v0.0.$(ARGS)

open:
	open https://github.com/inem/dotdotdot

uncommit:
	git reset --soft HEAD^
