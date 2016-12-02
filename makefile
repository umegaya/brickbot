sv:
	docker build -t $(IMAGE) .

restart: sv
	-docker kill cortana
	-docker rm cortana
	docker run -d --net="host" -p 8008:8008 --name cortana $(IMAGE)