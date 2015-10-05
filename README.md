# brickbot
slack RTM-API frontend for modularizing bot functionality by using multiple container


### motivation
I used to build many bot functionality for the project using slack. some of functionalities are project specific, and others are common for every project.

often such a specific and non-specific functionality is mixed in one source tree, and it leads to bad software architecture, which breaks reusability and testability. don't you write bot for new project by copy-and-paste some part of codes from bot of old project, or rewrite from scratch? 

ok. you are good programmer so you can well modularize your bot code and propertly reuse them. but how hard setup dependency for running them? like storage (say, mysql, redis), many gems or node modules. don't you have experience to fix your bot code because of conflict of dependency of bot functionalities (say, functionality A requires version X of gem G, B requires version Y)

yes. now we have a linux container technology to solve such problem. but if we pack all bot functionality in one container, same problem (mixup code, conflict dependency) as above will happen again. even if we seperate containers for each functionality, it is boring and wasteful to setup bot connection to slack for each container. 

brickbot solves these problem by:
- modularizing bot functionality by container each of them represent one functionality (called module container in this doc)
- root container launches module container and delegate RTM event to them, so that complex bot functionality can be achieved by combine simple module containers.

that provides
- easier reusability of single bot functionality
- unified management of all bot feature, by simple configuration

have some interest? then proceed to "how to run".


### how to run
1. edit settings.json.sample and rename to settings.json (#... is comment)
```
{
	"token": "your-slack-RTM-token",
	# channel which receives bot outgoing messages
	"main_channel": "general",
	# brickbot uses --net="host", so if port is collide, change this
	"bind_port": 8008,
	"docker": {
		# docker server address which child container runs
		"server_address": "localhost",
		# this must contains certs for docker server specified with "server_address" configuration
		"cert_path": "/docker-certs",
		# declare module container which gives necessary functionality to cortana
		"module_containers": {
			# declare containers which is delegated all RTM event from slack.
			# key name (eg. patcher) is used for determining which template files are used for format response from containers
			"patcher": {
				# config and host_config can receive same setting as
				# https://godoc.org/github.com/fsouza/go-dockerclient#Config, HostConfig
				"config": {
					"Image": "your/image",
					# Env often uses for giving customized configuration to module
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
docker build -t your_account/brickbot .
```

3. run in docker
```
./run.sh your_account/brickbot
```


### requirement for module container which is delegated RTM event
module containers are launched by brickbot with environment variable "BRICKBOT_ADDR". 
BRICKBOT_ADDR's format is address:port. such as "172.17.42.1:8008". 

each module container should establish persistent connection with brickbot using this address.
then brickbot sends all RTM event received through this connection, with text protocol (seperated by '\n').

module container processes such event if it should be processed, and return some JSON response to brickbot by writing json string with \n.
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

if no templates are given, brickbot sends these string to the channel specified in "main_channel" configuration.
otherwise you can format these response by putting "template" under templates/ directory.
for example, suppose container which name is "patcher" (see above configuration example), then you can put templates/patcher.json to format response from "patcher" module.
patcher.json looks like this:
```
{
	# value is go's text/template format.
	"payload_type": "hi, key1's value is {{.key1}}, key2's value is {{.value2}}.",
}
```
when brickbot receives above payload from patcher module container, then brickbot format message like this.
```
hi, key1's value is value1, key2's value is value2
```


### TODO
- messaging functionality between module container
```
{ "Kind": "foo", "To": "dest_module_container", "Payload": { "X": 1, "Y": 2 }}
```
- torelance to failure by launching multiple brickbot and stores RTMEvent to some reliable storage (cockroachDB?)
