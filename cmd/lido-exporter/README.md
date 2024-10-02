# Lido Exporter

Lido Exporter is a service that exports data from the Lido CSM smart contracts as Prometheus metrics.

## Features

- Collects metrics for penalties, exit requests, node operator info, keys info, bond info, and rewards info
- Supports Holesky network (Mainnet support can be easily added)
- Exports metrics in Prometheus format

## Usage

### Running with Docker

1. Build the Docker image:
   ```
   docker build -t lido-exporter .
   ```

2. Run the Docker container:
   ```
   docker run -d -p 8080:8080 -e LIDO_EXPORTER_NODE_OPERATOR_ID=<your_node_operator_id> -e LIDO_EXPORTER_NETWORK=<network_name> lido-exporter
   ```

### Running as a CLI Application

1. Build the application:
   ```
   go build -o lido-exporter cmd/lido-exporter/main.go
   ```

2. Run the application:
   ```
   ./lido-exporter --node-operator-id <your_node_operator_id> --network <network_name>
   ```

## Configuration

The service can be configured using the following methods (in order of precedence):

1. Environment variables
2. Command-line flags

Available settings:

- `LIDO_EXPORTER_NODE_OPERATOR_ID` (required): Node Operator ID
- `LIDO_EXPORTER_REWARD_ADDRESS` (optional): Reward address of Node Operator. It is used to calculate Node Operator ID if not set
- `LIDO_EXPORTER_NETWORK`: Network name (default: "holesky")
- `LIDO_EXPORTER_RPC_ENDPOINTS`: Comma-separated list of Ethereum RPC endpoints
- `LIDO_EXPORTER_WS_ENDPOINTS`: Comma-separated list of Ethereum WebSocket endpoints
- `LIDO_EXPORTER_PORT`: Port to listen on (default: "8080")
- `LIDO_EXPORTER_SCRAPE_TIME`: Scrape interval (default: 30s)
- `LIDO_EXPORTER_LOG_LEVEL`: Log level (default: "info")