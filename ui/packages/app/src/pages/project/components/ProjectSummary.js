import React from "react";
import { EuiIcon, EuiText } from "@elastic/eui";
import { Panel } from "./Panel";
import { Link } from "react-router-dom";

const moment = require("moment");

export const ProjectSummary = ({ project }) => {
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
    }
  ];

  const actions = (
    <Link to={`settings`}>
      <EuiText size="s">
        <EuiIcon type="arrowRight" /> Go to Project Settings
      </EuiText>
    </Link>
  );

  return (
    <Panel
      title="Project Info"
      items={items}
      actions={actions}
      iconType="machineLearningApp"
    />
  );
};
