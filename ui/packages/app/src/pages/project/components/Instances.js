import React from "react";
import { EuiFlexGroup, EuiFlexItem, EuiText, EuiTitle } from "@elastic/eui";
import { Panel } from "./Panel";
import { FeastJobsTable } from "./FeastJobsTable";
import { MerlinDeploymentsTable } from "./MerlinDeploymentsTable";
import { TuringRoutersTable } from "./TuringRoutersTable";

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

export const Instances = ({ project, feastIngestionJobs, models, routers }) => {
  const items = [
    {
      title: <Title title="Features Ingestion" href={`/feast/jobs/batch`} />,
      description: (
        <FeastJobsTable
          project={project}
          feastIngestionJobs={feastIngestionJobs}
        />
      )
    },
    {
      title: (
        <Title
          title="Merlin Deployments"
          href={`/merlin/projects/${project.id}/models`}
        />
      ),
      description: <MerlinDeploymentsTable project={project} models={models} />
    },
    {
      title: (
        <Title
          title="Turing Routers"
          href={`/turing/projects/${project.id}/routers`}
        />
      ),
      description: <TuringRoutersTable project={project} routers={routers} />
    }
  ];

  return <Panel title="Instances" items={items} type="row" iconType="compute" />;
};
