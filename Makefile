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