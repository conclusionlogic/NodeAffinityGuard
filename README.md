
# node-affinity-guard

**node-affinity-guard** is a Kubernetes utility designed to ensure that a specific StatefulSet or Deployment is running on the same node as a given floating IP (VIP). It manages node labels and node affinity rules to maintain the desired state within the cluster.

## Features

- Monitors the presence of a specified VIP in the cluster.
- Ensures that a specified StatefulSet or Deployment is running on the same node as the VIP.
- Manages node labels to facilitate the correct placement of workloads using node affinity rules.
- Configurable logging levels for easy debugging and monitoring.

## Prerequisites

- Kubernetes cluster
- Go environment set up
- Appropriate permissions to interact with the Kubernetes API

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/node-affinity-guard.git
   cd node-affinity-guard
   ```

2. Build the binary:

   ```bash
   go build -o node-affinity-guard main.go
   ```

3. Set up the required environment variables:

   ```bash
   export HOSTNAME=<your-node-hostname>
   export CHECK_INTERVAL=<interval-duration>
   export IP_ADDRESS=<floating-ip-address>
   export NODE_LABEL_KEY=<node-label-key>
   export NODE_LABEL_VALUE_ACTIVE=<node-label-active-value>
   export NODE_LABEL_VALUE_INACTIVE=<node-label-inactive-value>
   export LOG_LEVEL=<log-level>  # Options: debug, info, warn, error
   ```

## Usage

Run the `node-affinity-guard` binary:

```bash
./node-affinity-guard
```

The program will start monitoring the specified VIP and managing the node labels to ensure that the StatefulSet or Deployment is running on the same node as the VIP.

## Configuration

- `HOSTNAME`: The hostname of the node to monitor.
- `CHECK_INTERVAL`: The interval at which to check the presence of the VIP (e.g., "30s", "5m").
- `IP_ADDRESS`: The floating IP address to monitor.
- `NODE_LABEL_KEY`: The key of the node label to manage.
- `NODE_LABEL_VALUE_ACTIVE`: The value of the node label when the VIP is present on the node.
- `NODE_LABEL_VALUE_INACTIVE`: The value of the node label when the VIP is not present on the node.
- `LOG_LEVEL`: The level of logging detail (options: debug, info, warn, error).

## Example

```bash
export HOSTNAME="node-1"
export CHECK_INTERVAL="1m"
export IP_ADDRESS="192.168.1.100"
export NODE_LABEL_KEY="vip-active"
export NODE_LABEL_VALUE_ACTIVE="true"
export NODE_LABEL_VALUE_INACTIVE="false"
export LOG_LEVEL="info"

./node-affinity-guard
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
