build: .FORCE
	env GOOS=linux GOARCH=arm go build -o api.scainimatteo.dev .

deploy:
	@make build
	scp api.scainimatteo.dev pi:/home/pi/Projects/api.scainimatteo.dev
	scp config.json pi:/home/pi/Projects/api.scainimatteo.dev

.FORCE: