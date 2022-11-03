mosquitto-up:
	rm -rf /tmp/mosquitto  && \
	mkdir -p /tmp/mosquitto && \
	cp test/mosquitto/config/mosquitto.conf /tmp/mosquitto && \
	docker run --rm -it -p 1883:1883 -p 9001:9001 \
	--name mosquitto \
	--mount  type=bind,source=/tmp/mosquitto,target="/mosquitto/config" \
	eclipse-mosquitto

mosquitto-down:
	rm -rf /tmp/mosquitto

run:
	go run cmd/calibrator/main.go --path test/input/config.yml