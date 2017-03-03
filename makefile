sv:
	docker build -t $(IMAGE) .

restart: sv
	-docker kill cortana
	-docker rm cortana
	docker run -d --net="host" -p 8008:8008 --name cortana $(IMAGE)

bin: factory
	docker run --rm -ti -v 	$(BBOTPATH):/server umegaya/brickbot-factory bash -c "cd /server && go build -o brickbot"

image:
	docker build -t $(IMAGE) .

.PHONY: factory
factory:
	docker build -t umegaya/brickbot-factory factory

