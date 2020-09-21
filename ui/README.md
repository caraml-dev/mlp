# MLP UI

## Packages

We are using [Yarn workspace](https://classic.yarnpkg.com/en/docs/workspaces) to set up UI's package architecture that consists of:

- [mlp-ui](./packages/app) – React app (CRA) with common UI functionality.
- [@gojek/mlp-ui](./packages/lib/README.md) – library of shared UI components.

## Development environment

To work with the UI code, you will need to have the following tools installed:

- The [Node.js](https://nodejs.org/) JavaScript runtime.
- The [Yarn](https://yarnpkg.com/) package manager.

## Installing npm dependencies

To fetch all the dependencies and set a local development environment, run:

```shell script
yarn
```

## Running a local development server

You can start a development server for the MLP UI outside of a running MLP server by running:

```shell script
yarn start
```

This will open a browser window with the React app running on <http://localhost:3001/>. The page will reload if you make edits to the source code in both packages. You will also see any lint errors in the console.

Due to a `"proxy": "http://localhost:8080/v1"` setting in the mlp-ui's `package.json` file, any API requests from the React UI are proxied to `localhost` on port `8080` by the development server. This allows you to run a normal MLP server to handle API requests, while iterating separately on the UI.

```
[browser] ----> [localhost:3001 (dev server)] --(proxy API requests)--> [localhost:8080/v1 (MLP API)]
```

## Linting

We use [lint-staged](https://github.com/okonet/lint-staged) for the linter. To detect and automatically fix lint errors against staged git files, run:

```shell script
yarn lint
```

This is also available via the `lint-ui` target in the main MLP `Makefile`.

## Building the app for production

To build a production-optimized version of the React app to a `build` subdirectory, run:

```shell script
yarn app build
```

**NOTE:** You will likely not need to do this directly. Instead, this is taken care of by the `build` target in the main MLP `Makefile` when building the full binary.

## Integration into MLP

To build a MLP binary that includes a compiled-in version of the production build of the React app, change to the root of the repository and run:

```shell script
make build
```

This installs npm dependencies via Yarn, builds a production build of the React app, and then finally compiles in all web assets into the MLP binary.

## Running package specific commands

It is possible to run a command for the specific package (app/lib) from the root directory:

```shell script
yarn app <command>
```

or

```shell script
yarn lib <command>
```

For example, it's possible to build a library by running:

```shell script
yarn lib build
```
