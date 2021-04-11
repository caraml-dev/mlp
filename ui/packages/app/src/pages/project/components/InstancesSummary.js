import React from "react";
import { EuiLink } from "@elastic/eui";
import config from "../../../config";
import { Panel } from "./Panel";
import { MerlinDeployments } from "./MerlinDeployments";

export const InstancesSummary = ({ project, models }) => {
  const items = [
    {
      title: "FeatureSet Ingestion",
      description: "TODO"
    },
    {
      title: "Merlin Deployment",
      description: <MerlinDeployments project={project} models={models} />
    },
    {
      title: "Turing Routers",
      description: "TODO"
    },
    {
      title: "Pipelines",
      description: (
        <EuiLink href={config.CLOCKWORK_HOMEPAGE}>Clockwork Pipelines</EuiLink>
      )
    }
  ];

  return <Panel title="MLP Instances" items={items} />;
};
