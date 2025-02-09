# frp-port-keeper

This is a plugin for the awesome [frp reverse proxy](https://github.com/fatedier/frp).

| :exclamation: This is an early alpha version which needs further refactoring and improvements (see the [TODO](#todo) section) |
| ----------------------------------------------------------------------------------------------------------------------------- |

## What is it for?

The purpose of this plugin is to keep track of `remote_ports` that are being assigned
to frp clients upon initial connection to frp server. With this plugin, you can be
sure that whenever a client connects to an frp server, it would get the same `remote_port`
number that it was allocated initially.

## Requirements

- [Bun](https://bun.sh/) >= 1.1.x
- [frp](https://github.com/fatedier/frp) version >= v0.60.0

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/librepod/frp-port-keeper.git
   cd frp-port-keeper
   ```

2. **Install dependencies:**

   ```bash
   bun install
   ```

3. **Build the project:**

   ```bash
   bun run build
   ```

4. **Start the server:**

   ```bash
   bun run start
   ```

## Configuration

Create a `.env` file in the project root to specify configuration variables:

```env
PORT=8080
ALLOW_PORTS=6000-7000
```

**Environment Variables:**

- `PORT`: The port on which the frp-port-keeper server will run (default `8080`)
- `ALLOW_PORTS`: Comma-separated list of port ranges that clients are allowed to use (e.g., `6000-7000,8000-9000`)

## Usage

1. Register plugin in frps.yaml like this:

   ```yaml
   httpPlugins:
     - name: frp-port-keeper
       addr: '127.0.0.1:8080'
       path: /port-registrations
       ops:
         - NewProxy
   ```

2. Run the frp-port-keeper plugin (preferably via a systemd service) and make
   sure that it works fine (hit the `GET /ping` endpoint).

3. Run the frp server.

### Usage as Systemd service

You may want to delegate the control of frp-port-keeper to Systemd just like
you probably did with the frps service. There are sample Systemd unit files in
the `systemd` folder. Just tweak them to your liking and copy to the `/etc/systemd/system/`
folder.

## Development

To run the server in development mode with hot reloading:

```bash
bun run start:watch
```

**Scripts:**

- `bun run build`: Compiles the project into a single executable binary
- `bun run start`: Runs the project
- `npm run test`: Runs the test suite using mocha (TBD)

## API Documentation

### Endpoint details

This handler is used to allocate ports for the proxy requests storing the mapping of
`user` param specified in frpc.yaml and a free port available.

#### Request

The handler expects a JSON payload with the following structure:

```json
{
  "version": "0.1.0",
  "op": "NewProxy",
  "content": {
    "user": {
      "user": "myiphone"
    },
    "proxy_name": "myiphone.my_wireguard_proxy",
    "proxy_type": "udp"
  }
}
```

#### Response

The response body will be a JSON with the following structure:

```json
{
  "unchange": false,
  "content": {
    "user": {
      "user": "myiphone"
    },
    "proxy_name": "myiphone.my_wireguard_proxy",
    "proxy_type": "udp",
    "remote_port": 12345
  }
}
```

## TODO

- [ ] Implement Redis for persistent storage
- [ ] Add unit tests
- [ ] Improve error handling and input validation
- [ ] Implement proper logging system
- [ ] Pass `allow_ports` param via cli
- [x] Cross compile for other platforms (currently supports only amd64)
- [ ] Build docker images for other platforms (currently supports only linux/amd64 and linux/arm64)
- [ ] Update systemd files and instructions to run the plugin as systemd service
