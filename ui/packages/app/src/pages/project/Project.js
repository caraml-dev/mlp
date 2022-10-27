import React, { useContext } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiLoadingContent,
  EuiSpacer,
  EuiPageTemplate
} from "@elastic/eui";
import { ApplicationsContext, ProjectsContext } from "@gojek/mlp-ui";
import { Instances } from "./components/Instances";
import { ProjectSummary } from "./components/ProjectSummary";
import { Resources } from "./components/Resources";
import { useMerlinApi } from "../../hooks/useMerlinApi";
import { useTuringApi } from "../../hooks/useTuringApi";
import { useFeastCoreApi } from "../../hooks/useFeastCoreApi";
import { ComingSoonPanel } from "./components/ComingSoonPanel";

import imageCharts from "../../images/charts.svg";
import "./components/ListGroup.scss";

const Project = () => {
  const { apps } = useContext(ApplicationsContext);
  const { currentProject } = useContext(ProjectsContext);

  const [{ data: entities }] = useFeastCoreApi(
    `/entities?project=${currentProject?.name}`,
    { method: "GET" },
    undefined,
    !!currentProject
  );
  const [{ data: featureTables }] = useFeastCoreApi(
    `/tables?project=${currentProject?.name}`,
    { method: "GET" },
    undefined,
    !!currentProject
  );
  const [{ data: feastStreamIngestionJobs }] = useFeastCoreApi(
    `/jobs/ingestion/stream`,
    {
      method: "POST",
      body: JSON.stringify({
        include_terminated: true,
        project: (currentProject?.name || "").replace(/-/g, "_")
      })
    },
    undefined,
    !!currentProject
  );
  const [{ data: feastBatchIngestionJobs }] = useFeastCoreApi(
    `/jobs/ingestion/batch`,
    {
      method: "POST",
      body: JSON.stringify({
        include_terminated: true,
        project: (currentProject?.name || "").replace(/-/g, "_")
      })
    },
    undefined,
    !!currentProject
  );
  const [{ data: models }] = useMerlinApi(
    `/projects/${currentProject?.id}/models`,
    { method: "GET" },
    undefined,
    !!currentProject
  );
  const [{ data: routers }] = useTuringApi(
    `/projects/${currentProject?.id}/routers`,
    { method: "GET" },
    undefined,
    !!currentProject
  );

  return (
    <EuiPageTemplate panelled={false} restrictWidth="90%">
      <EuiPageTemplate.Section>
        {apps && !!currentProject ? (
          <>
            <EuiFlexGroup>
              <EuiFlexItem grow={3}>
                <ProjectSummary project={currentProject} />
              </EuiFlexItem>
              <EuiFlexItem grow={3}>
                <Resources
                  apps={apps}
                  project={currentProject}
                  entities={entities}
                  featureTables={featureTables}
                  models={models}
                  routers={routers}
                />
              </EuiFlexItem>
              <EuiFlexItem grow={3}>
                <ComingSoonPanel
                  title="Monthly Cost"
                  layout="vertical"
                  image={imageCharts}
                />
              </EuiFlexItem>
            </EuiFlexGroup>

            <EuiSpacer size="l" />

            <EuiFlexGroup>
              <EuiFlexItem grow={true}>
                <Instances
                  project={currentProject}
                  feastStreamIngestionJobs={feastStreamIngestionJobs}
                  feastBatchIngestionJobs={feastBatchIngestionJobs}
                  models={models}
                  routers={routers}
                />
              </EuiFlexItem>
            </EuiFlexGroup>

            <EuiSpacer size="l" />

            <EuiFlexGroup>
              <EuiFlexItem grow={4}>
                <ComingSoonPanel
                  title="Health Monitoring"
                  iconType="monitoringApp"
                />
              </EuiFlexItem>
              <EuiFlexItem grow={4}>
                <ComingSoonPanel title="Error Summary" iconType="bug" />
              </EuiFlexItem>
            </EuiFlexGroup>
          </>
        ) : (
          <EuiLoadingContent lines={5} />
        )}
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};

export default Project;
