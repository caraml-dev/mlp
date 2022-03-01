import React from "react";
import { useMlpApi } from "@gojek/mlp-ui";

const Pipeline = () => {
  const [{ data: pipelines, fetchPipelines }] = useMlpApi(`/pipelines`);

  return <h1>Pipeline goes here: {`${pipelines.workflows}`}</h1>;
};

export default Pipeline;
