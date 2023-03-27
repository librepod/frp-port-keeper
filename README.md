# How it works

The server exposes the following endpoint:
- /port-registrations

### POST /port-registrations
Register new machine by its `machineId`. Upon succesfull registration it should
response with a JSON with `remotePort` value, i.e. `{"message": "ok", "remotePort": 12345}`

### GET /port-registrations
Returns all the registered machines with their corresponding port registrations 
