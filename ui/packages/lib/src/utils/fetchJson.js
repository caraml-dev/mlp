import React from "react";
import { addToast } from "../components/Toast";
import { parseJson } from "./parseJson";

const fetchWithTimeout = async (url, options, time = 1800 * 1000) => {
  const controller = new AbortController();
  const config = { ...options, signal: controller.signal };

  setTimeout(() => {
    controller.abort();
  }, time);

  // https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API
  return fetch(url, config).catch(error => {
    if (error.name === "AbortError") {
      throw new Error("Response timed out");
    }
    throw new Error(error.message);
  });
};

const handleHttpStatus = authCtx => async response => {
  if (!response.ok) {
    if (response.status === 403 && authCtx) {
      setTimeout(() => authCtx.onLogout(), 2000);
    }
    const errMessage = await response.json();
    if (errMessage.error) {
      throw new Error(`${response.status}: ${errMessage.error}`);
    } else {
      throw new Error(`${response.status}: ${response.statusText}`);
    }
  }
  return response;
};

const fetchData = async (url, authCtx, options) => {
  let optionsWithAuth = { ...options };
  if (typeof authCtx !== "undefined") {
    optionsWithAuth.headers = {
      ...options.headers,
      Authorization: `Bearer ${authCtx.state.jwt}`
    };
  }
  return fetchWithTimeout(url, optionsWithAuth, options.timeout)
    .then(handleHttpStatus(authCtx))
    .then(response =>
      parseJson(response, !!options.parseBigInt).then(result => {
        const headers = Array.from(response.headers.entries()).reduce(
          (acc, [header, value]) => {
            acc[header] = value;
            return acc;
          },
          {}
        );
        return { body: result, headers: headers };
      })
    )
    .then(result => {
      if (options.addToast) {
        addToast({
          id: "fetch-success",
          title: "Success",
          color: "success",
          iconType: "check"
        });
      }
      return result;
    })
    .catch(error => {
      if (!options.muteError) {
        addToast({
          id: `fetch-error-${error.message}`,
          title: "Oops, there was an error",
          color: "danger",
          iconType: "alert",
          text: <p>{error.message}</p>
        });
      }
      throw error;
    });
};

export default fetchData;
