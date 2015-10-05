# slack-cortana
slack RTM-API frontend


### motivation
I used to build many bot functionality for the project using slack. some of functionalities are project specific, and others are common for every project.

often such a specific and non-specific functionality is mixed in one source tree, and it leads to bad software architecture, which breaks reusability and testability. don't you write bot for new project by copy-and-paste some part of codes from bot of old project, or rewrite from scratch? 

ok. you are good programmer so you can well modularize your bot code and propertly reuse them. but how hard setup dependency for running them? like storage (say, mysql, redis), many gems or node modules. don't you have experience to fix your bot code because of conflict of dependency of bot functionalities (say, functionality A requires version X of gem G, B requires version Y)

yes. now we have a linux container technology to solve such problem. but if we pack all bot functionality in one container, same problem (mixup code, conflict dependency) as above will happen again. even if we seperate containers for each functionality, it is boring and wasteful to setup bot connection to slack for each container. 

slack-cortana solves these problem by:
- modularizing bot functionality by container which represent one functionality
- root container launches module container and delegate RTM event to them

that provides
- easier reusability of single bot functionality
- unified management of all bot feature, by simple configuration

have some interest? then proceed to "how to run".


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
docker build -t your_account/cortana .
```

3. run in docker
```
./run.sh your_account/cortana
```


### requirement for module container which is delegated RTM event
modules are launched by cortana with environment variable "CORTANA_ADDR". 
CORTANA_ADDR's format is address:port. such as "172.17.42.1:8008". 

each module should establish persistent connection with cortana using this address.
then cortana sends all RTM event received through this connection, with text protocol (seperated by '\n').

module processes such event if it should be processed, and also return some response to cortana by writing json string with \n.
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
for example, suppose container which name is "patcher" (see above configuration example), then you can put templates/patcher.json to format response from "patcher" module.
patcher.json looks like this:
```
{
	# value is go's text/template format.
	"payload_type": "hi, key1's value is {{.key1}}, key2's value is {{.value2}}.",
}
```
when cortana receives above payload from patcher module, then cortana format message like this.
```
hi, key1's value is value1, key2's value is value2
```


