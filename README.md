# AOC 2022

### How to use

Test: `docker run --rm $(docker build -q .)`

Run: `docker run --rm $(docker build -q .) go run ./01`

New: `./createDay.sh 02`

### Dev environment

Open in VSCode, enable (Remote - Containers)[https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers], and re-open the project in the dev container. Then you'll have golang installed/isolated (and you won't need to keep building and running the container).

### Tech

- golang
- vscode devcontainers
- docker
