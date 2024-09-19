# U-Service Authentication

Firebase authentication proxy. Handle user-related operations.

## Requirements

- Git: Version control system
  - macOS:
    ```bash
    brew install git
    ```
  - Ubuntu:
    ```bash
    sudo apt install git-all
    ```
  - Windows: Try [Git bash](https://git-scm.com/downloads)
- [Go](https://go.dev/doc/install): The main development language
- (Optional, recommended) [direnv](https://direnv.net/docs/installation.html): environment variable manager
- [Docker](https://www.docker.com/products/docker-desktop/): Run the application locally
- (Optional, recommended) Make: Automated scripts for local development
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
- [Mockery](https://github.com/vektra/mockery): Generates mocks for Go interfaces. Requires Go.
  ```bash
  go install github.com/vektra/mockery/v2@v2.43.2
  ```
- (Optional) [Node.js](https://nodejs.org/en/download/package-manager): Can help install some dependencies. Chose LTS (stable version).
- [Firebase CLI](https://firebase.google.com/docs/cli?hl=fr): used to interact with our Firebase Stack.
  - If you installed Node.js, you can use:
    ```bash
    npm install -g firebase-tools
    ```

## Installation

Ensure you have [SSH configured on GitHub](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)
for your machine.

Make sure you are using git in SSH mode.

```bash
git config --global url.ssh://git@github.com/.insteadOf https://github.com/
```

Ensure the `GOPRIVATE` variable is set in your local terminal:

```bash
# You can skip this step if you have direnv configured.
export GOPRIVATE=github.com/in-rich/*
```

Install the project dependencies:

```bash
go mod download
```

Make sure you have access to the [firebase development project](https://console.firebase.google.com/u/0/project/inrich-dev/overview)
(requires a Google account).

Then, log in to the firebase CLI (select the correct project):

```bash
firebase login
```

Create a file named `firebase-key.json`, and paste the content of the 
[firebase private access key](https://console.cloud.google.com/security/secret-manager/secret/firebase-ley/versions?project=inrich-dev).

Make sure Docker is running, and available as a command:

```bash
docker ps -a
```

✅ Congrats, you're ready to go!

Check your environment:

```bash
firebase use
# Active Project: default (inrich-dev)
go version
# go version go1.23rc1 darwin/amd64
docker -v
# Docker version 24.0.7, build afdd53b
make -v
# GNU Make 3.81
echo $GOPRIVATE
# github.com/in-rich/*
```

## Usage

Run the server:

```bash
make run
```

Run tests:

```bash
make test
```

## For Windows Users

We recommend using a bash terminal emulator. One such example is [Git bash](https://git-scm.com/downloads).

You may also use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install). 

In both cases, make sure the dependencies you install are available under your bash
environment. This is automatic for Git Bash, but might require a separate setup
for WSL.
