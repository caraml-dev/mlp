import React, { useContext, useEffect, useState } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiLoadingContent,
  EuiSpacer,
  EuiPageTemplate
} from "@elastic/eui";
import { ApplicationsContext, CurrentProjectContext } from "@gojek/mlp-ui";
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
  const { projectId, project } = useContext(CurrentProjectContext);

  const [projectName, setProjectName] = useState("");
  useEffect(() => {
    if (project) {
      setProjectName(project.name);
    }
  }, [project]);

  const [{ data: entities }, fetchEntities] = useFeastCoreApi(
    `/entities?project=${projectName}`,
    { method: "GET" },
    undefined,
    false
  );
  const [{ data: featureTables }, fetchFeatureTables] = useFeastCoreApi(
    `/tables?project=${projectName}`,
    { method: "GET" },
    undefined,
    false
  );
  const [
    { data: feastStreamIngestionJobs },
    fetchFeastStreamIngestionJobs
  ] = useFeastCoreApi(
    `/jobs/ingestion/stream`,
    { method: "POST" },
    undefined,
    false
  );
  const [
    { data: feastBatchIngestionJobs },
    fetchFeastBatchIngestionJobs
  ] = useFeastCoreApi(
    `/jobs/ingestion/batch`,
    { method: "POST" },
    undefined,
    false
  );
  const [{ data: models }, fetchModels] = useMerlinApi(
    `/projects/${projectId}/models`,
    { method: "GET" },
    undefined,
    false
  );
  const [{ data: routers }, fetchRouters] = useTuringApi(
    `/projects/${projectId}/routers`,
    { method: "GET" },
    undefined,
    false
  );

  useEffect(() => {
    if (projectName) {
      fetchEntities();
      fetchFeatureTables();
      fetchFeastBatchIngestionJobs({
        body: JSON.stringify({
          include_terminated: true,
          project: projectName.replace(/-/g, "_")
        })
      });
      fetchFeastStreamIngestionJobs({
        body: JSON.stringify({
          include_terminated: true,
          project: projectName.replace(/-/g, "_")
        })
      });
      fetchModels();
      fetchRouters();
    }
  }, [
    projectName,
    fetchEntities,
    fetchFeatureTables,
    fetchFeastStreamIngestionJobs,
    fetchFeastBatchIngestionJobs,
    fetchModels,
    fetchRouters
  ]);

  return (
    <EuiPageTemplate panelled={false} restrictWidth="90%">
      <EuiPageTemplate.Section>
        {apps && project ? (
          <>
            <EuiFlexGroup>
              <EuiFlexItem grow={3}>
                <ProjectSummary project={project} />
              </EuiFlexItem>
              <EuiFlexItem grow={3}>
                <Resources
                  apps={apps}
                  project={project}
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
                  project={project}
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
