import { useCallback, useEffect, useReducer } from "react";
import fetchJson from "../utils/fetchJson";
import urlJoin from "proper-url-join";
import fetchMockData from "../utils/fetchMockData";
import * as queryString from "query-string";

const isStringifyEqual = (a, b) => JSON.stringify(a) === JSON.stringify(b);

const argumentsReducer = (state, action) => {
  const previous = state[action.name];
  const updated = action.value;
  return !isStringifyEqual(previous, updated)
    ? {
        ...state,
        [action.name]: updated
      }
    : state;
};

const zeroState = data => ({
  data: data,
  isLoading: false,
  isLoaded: false,
  error: null,
  headers: null
});

const dataFetchReducer = (state, action) => {
  switch (action.type) {
    case "FETCH_RESET":
      return zeroState(action.payload);
    case "FETCH_INIT":
      return {
        ...state,
        isLoading: true,
        isLoaded: false,
        error: null,
        headers: null
      };
    case "FETCH_SUCCESS":
      return {
        ...state,
        isLoading: false,
        isLoaded: true,
        data: action.payload,
        headers: action.headers
      };
    case "FETCH_FAILURE":
      return {
        ...state,
        isLoading: false,
        isLoaded: true,
        error: action.error
      };
    default:
      throw new Error();
  }
};

export const useApi = (
  endpoint,
  options,
  authCtx,
  result,
  callImmediately = true
) => {
  const [args, dispatchArgsUpdate] = useReducer(argumentsReducer, {
    result,
    options,
    authCtx
  });

  useEffect(() => {
    dispatchArgsUpdate({ name: "options", value: options });
  }, [options]);

  useEffect(() => {
    dispatchArgsUpdate({ name: "authCtx", value: authCtx });
  }, [authCtx]);

  useEffect(() => {
    dispatchArgsUpdate({ name: "result", value: result });
  }, [result]);

  const [state, dispatch] = useReducer(
    dataFetchReducer,
    args.result,
    zeroState
  );

  const fetchData = useCallback(
    options => {
      let didCancel = false;

      const apiOptions = !!options
        ? {
            ...args.options,
            ...options
          }
        : args.options;

      dispatch({ type: "FETCH_INIT" });
      (apiOptions.useMockData && apiOptions.mock
        ? fetchMockData(apiOptions.mock, apiOptions)
        : fetchJson(
            queryString.stringifyUrl({
              url: urlJoin(apiOptions.baseApiUrl, endpoint),
              // query params supplied via `apiOptions` have a higher priority
              // and will override query param with the same name if it is
              // present in the `endpoint`
              query: apiOptions.query
            }),
            args.authCtx,
            apiOptions
          )
      )
        .then(result => {
          if (!didCancel)
            dispatch({
              type: "FETCH_SUCCESS",
              payload: result.body,
              headers: result.headers
            });
        })
        .catch(error => {
          if (!didCancel) dispatch({ type: "FETCH_FAILURE", error: error });
        });

      return {
        cancel: () => {
          didCancel = true;
        }
      };
    },
    [args, endpoint]
  );

  useEffect(() => {
    dispatch({ type: "FETCH_RESET", payload: args.result });

    if (callImmediately) {
      const call = fetchData();
      return call.cancel;
    }
  }, [args.result, callImmediately, fetchData]);

  return [state, fetchData];
};
