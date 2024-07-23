# Go Demo API

## Getting Started

You can clone this repository by running the following command:

```bash
git clone https://github.com/aaronzi/go-demo-api.git
```

After you cloned the repository, open the folder in vscode.
It will automatically detect the devscontainer configuration and ask you to reopen the folder in the container.
This also inmcludes the database container.

From the debug menu, you can run the API by selecting the `Launch Server` configuration.

## Using the Docker Image
You can use the Docker Image from either Docker Hub or GitHub Container Registry.
You can find the Docker Image on Docker Hub [here](https://hub.docker.com/r/aaronzi/go-demo-api) and on GitHub Container Registry on the following link [here](https://github.com/aaronzi/go-demo-api/packages).

In order to run the Docker Image, you should use it together with mysql:latest running in the same docker network or in the same docker-compose file.

## Rebuild the Swagger Documentation

To build the Swagger documentation, run the following command:

```bash
swag init -g cmd/movie-api/main.go --parseDependency --parseInternal -o docs
```

This will generate the `docs` directory with the Swagger documentation.

> **Note:** You should run this command every time you make changes to the API.
