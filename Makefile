build: .FORCE
	env GOOS=linux GOARCH=arm go build -o api.scainimatteo.dev .

deploy:
	@make build
	scp api.scainimatteo.dev pi:/home/pi/Projects/api.scainimatteo.dev/api.scainimatteo.dev.new
	scp config.json pi:/home/pi/Projects/api.scainimatteo.dev
	ssh pi "mv /home/pi/Projects/api.scainimatteo.dev/api.scainimatteo.dev.new /home/pi/Projects/api.scainimatteo.dev/api.scainimatteo.dev; sudo systemctl restart api.scainimatteo.dev"

.FORCE: