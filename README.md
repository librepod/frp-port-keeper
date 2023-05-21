# frp-port-keeper
This is a plugin for the awesome [frp reverse proxy](https://github.com/fatedier/frp). 

## What is it for?
The purpose of this plugin is to keep track of `remote_ports` that are being assigned
to frp clients upon initial connection to frp server. With this plugin, you can be
sure that whenever a client connects to an frp server, it would get the same `remote_port`
number that it was allocated initially.

## Implementation
frp-port-keeper is a simple server that exposes a `POST /port-registrations` endpoint
that processes the *NewProxy* payload from the frp server. It is utilizing a simple
key/value store to track ports and correspinding users. Port allocation data
persists in json files under the `gokv` folder (The `gokv` folder is created in the
same directory where the frp-port-keeper executable is executed from).

### Endpoint details
This handler is used to allocate ports for the proxy requests storing the mapping of
`user` param specified in frpc.ini and a free port available.

#### Request
The hendler expects a JSON payload with the following structure:
```json
{
	"version": "0.1.0",
	"op": "NewProxy",
	"content": {
		"user": {
			"user": "myiphone",
		},
	"proxy_name": "myiphone.my_wireguard_proxy",
	"proxy_type": "udp"
	}
}
```
The corresponding frpc.ini config that generates this kind of payload should have 
the following mandatory parameters specified:
```ini
[common]
server_addr = <your_server_address>
server_port = 7000
user = myiphone

[my_wireguard_proxy]
type = udp
local_ip = 127.0.0.1
local_port = 51820
# remote_port may be omitted since the actual remote_port value will be assigned by the plugin
# remote_port = 1000
```

#### Response
The response body will be a JSON with the following structure:

```json
{
	"unchange": false,
	"content": {
		"user": {
			"user": "myiphone",
		},
		"proxy_name": "myiphone.my_wireguard_proxy",
		"proxy_type": "udp",
		"remote_port": 12345
	}
}
```

If the request is not valid due to missing mandatory fields in frpc.ini config, the
response will be of status 400 with the following content:
```json
{
	"error":   "VALIDATEERR",
	"message": "Invalid inputs. Please check your frpc.ini config",
}
```

In case if there is an internal error or no more free ports left to allocate,
the status would be 200 with the following content:

```json
{
	"reject": true,
	"reject_reason": "<reject reason>"
}
```

## Requirements

frp version >= v0.48.0

It is possible that the plugin works for older version even though it has not been tested.

## How to run


## TODO
[ ] Add unit tests  
[ ] Add proper error handling in case if payload is not as expected  
[ ] Cross compile for other platforms (currently supports only amd64)  
[ ] Refactor by improving modules/folder structure following golang best practices  

