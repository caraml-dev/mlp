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
      Authorization: `Bearer ${authCtx.state.accessToken}`
    };
  }
  return fetchWithTimeout(url, optionsWithAuth, options.timeout).then(
    handleHttpStatus(authCtx)
  );
};

export default async (url, authCtx, options = {}) => {
  return fetchData(url, authCtx, options)
    .then(response =>
      parseJson(response).then(result => {
        var headers = {};
        for (var pair of response.headers.entries()) {
          headers[pair[0]] = pair[1];
        }
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
