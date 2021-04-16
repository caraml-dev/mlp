import React from "react";
import {
  EuiIcon,
  EuiLink,
  EuiText
} from "@elastic/eui";
// import { CollapsibleLabelsPanel } from "@gojek/mlp-ui";
import { Panel } from "./Panel";

const moment = require("moment");

export const ProjectSummary = ({ project, environments }) => {
  const items = [
    {
      title: "Name",
      description: project.name
    },
    {
      title: "Team",
      description: project.team
    },
    {
      title: "Stream",
      description: project.stream
    },
    {
      title: "Created at",
      description: moment(project.created_at).format("DD-MM-YYYY")
    },
    {
      title: "Modified at",
      description: moment(project.updated_at).format("DD-MM-YYYY")
    },
    // {
    //   title: "Labels",
    //   description: <CollapsibleLabelsPanel labels={project.labels} />
    // },
    // {
    //   title: "Administrator",
    //   description: project.administrators.join(", "),
    // },
  ];

  const actions = (
    <EuiLink href={`/projects/${project.id}/settings`}>
      <EuiText size="s"><EuiIcon type="arrowRight" /> Go to Project Settings</EuiText>
    </EuiLink>
  );

  return <Panel title="Project Info" items={items} actions={actions} iconType="machineLearningApp" />;
};
