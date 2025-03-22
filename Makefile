.PHONY:

build:
		go build -0 ./.bin/bot cmd/bot/main.go

run: build
		./.bin/bot

build-image:
		docker build -t polyk005/pocketer:0.1 .

start-container:
		docker run --env-file .env -p 80:80 polyk005/pocketer:0.1