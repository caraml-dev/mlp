# Machine Learning Platform

Machine Learning Platform (MLP) is a unified set of products for developing and operating the machine learning systems at the various stages of machine learning life cycle. The typical ML life cycle can be viewed through the following nine stages:

![Machine learning life cycle](./docs/diagrams/machine_learning_life_cycle.drawio.svg)

MLP Products are systems and services that are specifically built to solve one or multiple stages of the machine learning life cycle's problems. Currently, we have published the following MLP products:

* [**Feast**](https://github.com/caraml-dev/caraml-store) - For managing and serving machine learning features.
* [**Merlin**](https://github.com/caraml-dev/merlin) - For deploying, serving, and monitoring machine learning models.
* [**Turing**](https://github.com/caraml-dev/turing) - For designing, deploying, and evaluating machine learning experiments.

## Architecture overview

![Architecture overview](./docs/diagrams/architecture_overview.drawio.svg)

The MLP Server provides REST API used across MLP Products. It exposes a shared concepts such as [ML Project](./docs/concepts.md#ml-project). This repository also hosts [Go](./api/pkg) and React ([@caraml-dev/ui-lib](./ui/packages/lib)) libraries used to build a common MLP functionailty.

## Getting started

### Prerequisites

1. [Google Oauth credential](https://developers.google.com/identity/protocols/oauth2/javascript-implicit-flow)

    MLP uses Google Sign-in to authenticate the user to access the API and UI. After you get the client ID, specify it into `REACT_APP_OAUTH_CLIENT_ID` in `.env.development` file.

### From Docker Compose

If you already have [Docker](https://docs.docker.com/get-docker/) installed, you can spin up MLP and its dependencies by running:

```shell script
docker-compose up
```

MLP will now be reachable at <http://localhost:8080>.

### From source

To build and run MLP from the source code, you need to have [Go](https://golang.org/doc/install), [Node.js](https://nodejs.org/), and [Yarn](https://yarnpkg.com/) installed. 
You will also need a running Postgresql database, Keto, and Vault servers. 
MLP uses Docker to make the task of setting up the dependencies a little easier. You can run `make local-env` to starting up all those dependencies.

```shell script
make local-env
make run
```

OR

```shell script
# `make` will execute `make local-env` and `make run`
make
```

## Documentation

Go to the [docs](/docs) folder for the full documentation and guides.

## React UI development

For more information on building, running, and developing the UI app and library, see the UI's [README.md](ui/README.md).
