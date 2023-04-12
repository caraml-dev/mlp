# @caraml-dev/ui-lib Library

A library of common React components used by the MLP solutions.

## Install

```shell script
yarn add @caraml-dev/ui-lib
```

## Packages

### `auth`

Contains React context (`AuthContext`) and necessary components (`AuthContextProvider` and `PrivateRoute`) for implementing views that require a user to be authenticated to access them.

### `components`

Generic UI components used by MLP apps (such as header, breadcrumbs, project dropdown etc.).

### `hooks`

Custom React hooks:

- [`useApi`](./src/hooks/useApi.js) - generic hook for interacting with REST API.

**Example:**

```js
const DummyButton = () => {
  const [response, sendRequest] = useApi(
    "https://run.mocky.io/v3/7533b0dd-df1b-4f3d-9068-a6ca9a448b8d",
    {
      method: "POST",
      body: JSON.stringify({ hello: "world" }),
      headers: { "Content-Type": "application/json" }
    }, // request options
    {}, // optional, AuthContext, if the API requires authentication
    {}, // initial value of the response.data
    false // whether or not to send the request to the API immediately
  );

  // Send request explicitly
  const onClick = sendRequest;

  // Log response payload when request succeeded
  useEffect(() => {
    const { data, isLoading, isLoaded, error } = response;

    console.log(data);
  }, [response]);

  return <button onClick={onClick}>Hello World!</button>;
};
```

---

- [`useMlpApi`](./src/hooks/useMlpApi.js) - custom React hook to talk to `mlp-api`. It utilizes `useApi` hook under the hook, but pre-populates it with the `AuthContext` and modifies `endpoint` value to prefix it with the root URL where MLP-api is accessible.

**Example:**

```js
const [response, fetch] = useMlpApi(
  `/v1/projects/${projectId}/environments`,
  {}, // request options
  [], // initial value of the response.data
  true // whether or not to send the request to the API immediately
);
```

---

- [`useToggle`](src/hooks/useToggle.js) - custom React hook for defining a boolean value that can only be switched on and off. To be used in pop-overs, modals etc, where it can represent whether to show or hide a component.

**Example:**

```js
const [isShowing, toggle] = useToggle(
  true // initialState – optional, default false
);

useEffect(() => {
  if (isShowing) {
    toggle();
  }
  console.log(isShowing);
}, [isShowing]);
```

**Output:**

```js
true; // initialState
false; // calling `toggle()` switched the state
```

---

### `providers`

Context providers that supply config/data to children components:

- [`MlpApiContextProvider`](./src/providers/api)
- [`ApplicationsContextProvider`](./src/providers/application)
- [`ProjectsContextProvider`](./src/providers/project)

---

### `utils`

Misc utils.

## Available Scripts

### Dev Server

```shell script
yarn start
```

### Production Bundle

```shell script
yarn build
```

### Run Lint

```shell script
yarn lint
```

To let the linter to try fixing observed issues, run:

```shell script
yarn lint:fix
```

## Link Library Locally

It can be handy, to link this library locally, when you are working on the application, that has `@caraml-dev/ui-lib` as a dependency. For doing it, run:

```shell script
yarn run link
```

This will link `@caraml-dev/ui-lib` as well as `react` and `react-dom` locally, so it can be used in your application. Then run following commands from your project's directory:

```shell script
cd <your app project>

yarn link @caraml-dev/ui-lib
yarn link react
yarn link react-dom
```

When you no longer want to have a local link of `@caraml-dev/ui-lib` and want to resolve the library from the npm registry, run:

```shell script
cd </path/to/mlp-ui/packages/lib>

yarn run unlink
```

and then:

```shell script
cd <your app project>

yarn unlink @caraml-dev/ui-lib
yarn unlink react
yarn unlink react-dom
```
