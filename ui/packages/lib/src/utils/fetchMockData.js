import React from "react";
import { addToast } from "../components/Toast";
import { parseJson } from "./parseJson";

const mockData = async file => {
  return new Response(JSON.stringify(file), {
    status: 200,
    headers: { "Content-Type": "application/json" }
  });
};

export default async (file, options) => {
  let blob = "mock-".concat(
    Math.random()
      .toString(36)
      .substring(8)
  );
  return mockData(file)
    .then(parseJson)
    .then(result => {
      if (options.addToast) {
        addToast({
          id: `fetch-success${blob}`,
          title: "Success",
          color: "success",
          iconType: "check"
        });
      }
      return result;
    })
    .catch(error => {
      addToast({
        id: `fetch-error${blob}`,
        title: "Oops, there was an error",
        color: "danger",
        iconType: "alert",
        text: <p>{error.message}</p>
      });
      throw error;
    });
};
