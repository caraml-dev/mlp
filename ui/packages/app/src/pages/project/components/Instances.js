import React, { useContext } from "react";
import { EuiFlexGroup, EuiFlexItem, EuiText, EuiTitle } from "@elastic/eui";
import { Panel } from "./Panel";
import { FeastJobsTable } from "./FeastJobsTable";
import { MerlinDeploymentsTable } from "./MerlinDeploymentsTable";
import { TuringRoutersTable } from "./TuringRoutersTable";
import { ApplicationsContext } from "@gojek/mlp-ui";

const Title = ({ title, href }) => {
  return (
    <EuiFlexGroup alignItems="center" gutterSize="m">
      <EuiFlexItem grow={false}>
        <EuiTitle size="xxs">
          <h4>{title}</h4>
        </EuiTitle>
      </EuiFlexItem>
      <EuiFlexItem>
        <EuiText size="xs">
          <a href={href}>View all</a>
        </EuiText>
      </EuiFlexItem>
    </EuiFlexGroup>
  );
};

export const Instances = ({ project, models, routers }) => {
  const { apps } = useContext(ApplicationsContext);

  const items = [
    ...(apps.some(a => a.name === "Feast")
      ? [
          {
            title: (
              <Title
                title="Features Ingestion"
                href={`${
                  apps.find(a => a.name === "Feast").homepage
                }/jobs/stream`}
              />
            ),
            description: (
              <FeastJobsTable
                project={project}
                homepage={apps.find(a => a.name === "Feast").homepage}
              />
            )
          }
        ]
      : []),
    ...(apps.some(a => a.name === "Merlin")
      ? [
          {
            title: (
              <Title
                title="Merlin Deployments"
                href={`${
                  apps.find(a => a.name === "Merlin").homepage
                }/projects/${project.id}/models`}
              />
            ),
            description: (
              <MerlinDeploymentsTable
                project={project}
                models={models}
                homepage={apps.find(a => a.name === "Merlin").homepage}
              />
            )
          }
        ]
      : []),
    ...(apps.some(a => a.name === "Turing")
      ? [
          {
            title: (
              <Title
                title="Turing Routers"
                href={`${
                  apps.find(a => a.name === "Turing").homepage
                }/projects/${project.id}/routers`}
              />
            ),
            description: (
              <TuringRoutersTable
                project={project}
                routers={routers}
                homepage={apps.find(a => a.name === "Turing").homepage}
              />
            )
          }
        ]
      : [])
  ];

  return (
    <Panel title="Instances" items={items} type="row" iconType="compute" />
  );
};
