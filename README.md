# README.md

## Project Setup

This project is written in Go. It uses the Go module for dependency management.

### Prerequisites

- Go (version 1.16 or higher)
- Node.js (version 14 or higher)
- npm (version 6 or higher)
- 
### Installation

1. Clone the repository:
2. 
```bash
git clone https://gitlab.utc.fr/royhucheradorni/ia04.git
```

2. Install Go dependencies:
```bash
go mod download
```

## Running the Program

1. Start the Simulation:

```bash
go run gitlab.utc.fr/royhucheradorni/ia04.git/cmd/simulation
cd ia04
```

If you are using WSL, you may need to use the following command instead:

```bash
GOOS=windows go run gitlab.utc.fr/royhucheradorni/ia04.git/cmd/simulation
```