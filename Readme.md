# Use the CLI

## Create the mqtt server
```bash
make mosquitto-up
```

## Start up the up
Open another terminal and run the app with the simple config
```bash
 go run ./cmd/calibrator/main.go --path test/input/simple_config.yml
```

## Check that the calibration is working by mimic the sensor behavior using the cli.

Publish to the sensor topic
```bash
go run ./cmd/cli/main.go --t 17.0
```

Publish to the thermostat topic
```bash
go run ./cmd/cli/main.go --c 0.0 --th 16.33 --topic topic/thermostat --isth true
```
