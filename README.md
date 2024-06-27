
# node-affinity-guard

**node-affinity-guard** is a Kubernetes utility designed to ensure that a specific StatefulSet or Deployment is running on the same node as a given floating IP (VIP).

It manages node labels for node affinity rules to maintain the desired state within the cluster.

## Features

- Monitors the presence of a specified VIP in the cluster.
- Ensures that a specified StatefulSet or Deployment is running on the same node as the VIP.
- Manages node labels to facilitate the correct placement of workloads using node affinity rules.
- Configurable logging levels for easy debugging and monitoring.

## Prerequisites

- Kubernetes cluster
- Appropriate permissions to interact with the Kubernetes API

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/conclusionlogic/node-affinity-guard.git
   cd node-affinity-guard
   ```

2. Build the binary:

   ```bash
   go build -o node-affinity-guard main.go
   ```

   or with Docker:

   ```bash
   docker build .
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

# Using the node-affinity-guard Helm Chart

This repository contains the Helm chart for deploying the **node-affinity-guard**. Follow the steps below to add the Helm repository and deploy the chart.

## Helm Chart prerequisites

- Helm installed
- Kubernetes cluster

## Adding the Helm Repository

1. Add the Helm repository:

    ```bash
    helm repo add node-affinity-guard https://conclusionlogic.github.io/node-affinity-guard/
    ```

2. Update the Helm repository to ensure you have the latest version of the chart:

    ```bash
    helm repo update
    ```

## Deploying the Chart

1. Install the **node-affinity-guard** chart with a release name of your choice (e.g., `my-release`):

    ```bash
    helm install my-release node-affinity-guard/node-affinity-guard
    ```

2. To customize the deployment, you can pass custom values via the `--set` flag or provide a `values.yaml` file. For example:

    ```bash
    helm install my-release node-affinity-guard/node-affinity-guard --set hostname=node-1,checkInterval=1m,ipAddress=192.168.1.100,nodeLabelKey=vip-active,nodeLabelValueActive=true,nodeLabelValueInactive=false,logLevel=info
    ```

## Upgrading the Chart

To upgrade the chart to a new version:

   ```bash
   helm upgrade my-release node-affinity-guard/node-affinity-guard
   ```

## Uninstalling the Chart

To uninstall the deployed release:

   ```bash
   helm uninstall my-release
   ```

## Values

The following table lists the configurable parameters of the **node-affinity-guard** chart and their default values.

| Parameter                | Description                                         | Default            |
|--------------------------|-----------------------------------------------------|--------------------|
| `hostname`               | The hostname of the node to monitor                 | `""`               |
| `checkInterval`          | The interval at which to check the VIP              | `"1m"`             |
| `ipAddress`              | The floating IP address to monitor                  | `""`               |
| `nodeLabelKey`           | The key of the node label to manage                 | `"vip-active"`     |
| `nodeLabelValueActive`   | The value of the node label when VIP is present     | `"true"`           |
| `nodeLabelValueInactive` | The value of the node label when VIP is not present | `"false"`          |
| `logLevel`               | The level of logging detail                         | `"info"`           |

## Example

To deploy the **node-affinity-guard** with specific values:

   ```bash
   helm install my-release node-affinity-guard/node-affinity-guard --set hostname=node-1,checkInterval=1m,ipAddress=192.168.1.100,nodeLabelKey=vip-active,nodeLabelValueActive=true,nodeLabelValueInactive=false,logLevel=info
   ```

For more details, refer to the [Helm documentation](https://helm.sh/docs/).

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
