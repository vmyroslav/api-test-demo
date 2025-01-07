# API Testing with Hoverfly Demo

This project demonstrates how to effectively use Hoverfly for API testing in various non-trivial scenarios. It includes examples of capturing, simulating, and managing API interactions using Hoverfly, along with different client generation approaches.

## Features

- Comprehensive Hoverfly integration for API testing
- Support for multiple testing modes (capture, simulate, spy)
- Automated test execution in both local and containerized environments
- Client generation using both oapi-codegen and OpenAPI Generator
- Custom simulation post-processing capabilities
- Docker Compose setup for consistent development environment
- Task-based workflow management using Taskfile

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or later
- Task (taskfile) installed (`go install github.com/go-task/task/v3/cmd/task@latest`)

## Project Structure

```
.
├── api/                    # API specifications
├── client/                # Generated API clients
│   ├── oapi/             # oapi-codegen generated client
│   └── openapi/          # OpenAPI Generator client
├── docker/               # Dockerfile definitions
├── testdata/            # Test data and Hoverfly simulations
│   └── hoverfly/        # Hoverfly simulation files
├── tools/               # Project tools
│   └── postprocessor/   # Simulation post-processor
├── docker-compose.yml   # Docker services configuration
└── Taskfile.yml        # Task definitions
```

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/vmyroslav/api-test-demo.git
cd api-test-demo
```

2. Build the containers:
```bash
task build:container
```

## Testing Modes

### Local Testing

Run tests on your local machine:

```bash
# Run tests against real API
task test:local

# Run tests with Hoverfly capture mode
task test:local:capture

# Run tests with Hoverfly simulate mode
task test:local:simulate
```

## Post-processing Simulations

The project includes a custom post-processor for Hoverfly simulations. To process a simulation:

```bash
task hoverfly:process-simulation SIMULATION_FILE=<path> PROCESSOR=default
```

Available processors:
- `default`: Standard simulation processor
- `null`: No-op processor

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

[Add your license information here]

## Acknowledgments

- [Hoverfly](https://hoverfly.io/) - API virtualization tool
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI client generator for Go
- [OpenAPI Generator](https://openapi-generator.tech/) - OpenAPI client generator