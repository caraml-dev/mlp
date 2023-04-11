# Developer Guide

This guide helps you get started developing MLP.

## Dependencies

Make sure you have the following dependencies installed before setting up your developer environment:

- [Git](https://git-scm.com/)
- [Go 1.13](https://golang.org/doc/install) (see [go.mod](../go.mod#L3) for minimum required version)
- [Node.js 12](https://nodejs.org) (LTS version or greater)
- [Yarn](https://yarnpkg.com)
- [Docker](https://docs.docker.com/get-docker/)

### Google Oauth Credential

MLP uses Google Sign-in to authenticate the user to access the API and UI. You can follow this [tutorial](https://developers.google.com/identity/sign-in/web/sign-in#create_authorization_credentials) to create your credential. After you get the client ID, specify it into `OAUTH_CLIENT_ID` in the root's [`.env.development`](../..env.development) and `REACT_APP_OAUTH_CLIENT_ID` in the [`ui/packages/app/.env.development`](../ui/packages/app/.env.development).

## Download MLP

We recommend using the Git command-line interface to download the source code for the MLP project:

1. Create `github.com/caraml-dev` directory in your `$GOPATH` by running: `mkdir -p $GOPATH/src/github.com/caraml-dev` from your terminal.
2. Run `cd $GOPATH/src/github.com/caraml-dev` to go to the newly created caraml-dev directory.
3. Run `git clone git@github.com/caraml-dev/mlp.git mlp` to downloads MLP repository to a new `mlp` directory.
4. Open the `mlp` directory using your favorite code editor.

## Build MLP

MLP server consists of two components; the ui (frontend) and the api (backend).

### UI (Frontend)

Before we can build the UI assets, we need to install the dependencies:

```shell script
make init-dep-ui
```

After the command has finished, we can start building our UI source code:

```shell script
make start-ui
```

Once `make start-ui` has built the assets, it will continue to do so whenever any of the UI files change. This means you don't have to manually build the assets every time you change the code.

Next, we'll build the web server that will serve the frontend assets we just built.

### API (Backend)

Build and run the backend by running:

```shell script
make init-dep-api
make build-api
make local-db
make run
```

These commands download dependencies, compile the Go source code, start Postgres Docker container, and start the web server.

By default, you can access the web server at <http://localhost:8080/>.

## Test MLP

The test suite consists of two types of tests: unit tests and end-to-end tests.

### Run unit tests

```shell script
make test-api
```

### Run end-to-end tests

```shell script
make it-test-api-local
```

## Configure MLP for development

MLP follows the twelve-factor app by storing the config in the environment. By default, we include the `.env.development` in the root directory into Makefile. You can override the default configuration by changing this file. If you run the MLP server without using `make` command, you need to set the environment variables value by yourself.

### Configure Database

MLP requires Postgresql to run. We uses Docker to make the task of setting up Postgresql databases a little easier.

```shell script
make local-db
```

The script runs Postgresql Docker container.

## Build a Docker image

To build a Docker image, run:

```shell script
build-docker
```

The resulting image will be tagges as `gojektech/mlp:dev`.
