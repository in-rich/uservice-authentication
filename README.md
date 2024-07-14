# U-Service Authentication

Firebase authentication proxy.

## Requirements

- [Go](https://go.dev/doc/install)
- Make:
  - macOS: 
    ```bash
    brew install make
    ```
  - Ubuntu:
    ```bash
    sudo apt-get install make
    ```
  - Windows: Install [chocolatey](https://chocolatey.org/install) (from a PowerShell with admin privileges), then run:
    ```bash
    choco install make
    ```

## Installation

Install dependencies:

```bash
go mod download
```

## Usage

Run the server:

```bash
make run
```

## Test

Install [gotestsum](https://github.com/gotestyourself/gotestsum):

```bash
go install gotest.tools/gotestsum@latest
```

Install [mockery](https://vektra.github.io/mockery/latest/installation/):

```bash
go install github.com/vektra/mockery/v2@v2.43.2
```

Run tests:

```bash
make test
```
