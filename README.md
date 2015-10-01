# slack-cortana
slack RTM-API frontend for golang

### how to run
1. edit settings.json.sample and rename to settings.json
```
{
	"token": "your-slack-RTM-token",
	# channel which receives bot outgoing messages
	"main_channel": "general",
	# slack-cortana uses --net="host", so if port is collide, change this
	"bind_port": 8008,
	"docker": {
		# docker server address which child container runs
		"server_address": "localhost",
		# mount $DOCKER_CERT_PATH to this path
		"cert_path": "/docker-certs",
		"containers": {
			# declare containers which is delegated all RTM event from slack.
			# key name (eg. patcher) is used for determining which template files are used for format response from containers
			"patcher": {
				# config and host_config can receive same setting as
				# https://godoc.org/github.com/fsouza/go-dockerclient#Config, HostConfig
				"config": {
					"Image": "your/image",
					"Env": ["YOUR_ENV=your_value"]
				}
				"host_config": {
					...
				}
			}
		}
	}
}
```

2. build docker image 
```
docker build -t your_account/cortana .
```

3. run in docker
```
./run.sh your_account/cortana
```

### requirement for container which is delegated RTM event
below we call these containers "sub-brain". sub-brains are launched by cortana with environment variable "CORTANA_ADDR". 
CORTANA_ADDR's format is address:port. such as "172.17.42.1:8008". 

each sub-brain should establish persistent connection with cortana using this address.
then cortana sends all RTM event received through this connection, with text protocol (seperated by '\n').

sub-brain processes such event if it should be processed, and also return some response to cortana by writing json string with \n.
json string format is like following: 
```
{
	"Kind": "payload_type",
	"Payload": {
		"key1": "value1",
		"key2": "value2",
		...
	}
}
```

if no templates are given, cortana sends these string to the channel specified in "main_channel" configuration.
otherwise you can format these response by putting "template" under templates/ directory.
for example, suppose container which name is "patcher" (see above configuration example), then you can put templates/patcher.json to format response from "patcher" sub-brain.
patcher.json looks like this:
```
{
	# value is go's text/template format.
	"payload_type": "hi, key1's value is {{.key1}}, key2's value is {{.value2}}.",
}
```
when cortana receives above payload from patcher sub-brain, then cortana format message like this.
```
hi, key1's value is value1, key2's value is value2
```


