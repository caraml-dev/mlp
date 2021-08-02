import React, { useContext, useEffect } from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiLoadingContent,
  EuiPage,
  EuiPageBody,
  EuiSpacer
} from "@elastic/eui";
import { ApplicationsContext, CurrentProjectContext } from "@gojek/mlp-ui";
import { Instances } from "./components/Instances";
// import { InstancesSummary } from "./components/InstancesSummary";
import { ProjectSummary } from "./components/ProjectSummary";
import { Resources } from "./components/Resources";
import { useMerlinApi } from "../../hooks/useMerlinApi";
import { useTuringApi } from "../../hooks/useTuringApi";

import "./components/ListGroup.scss";
import { useFeastCoreApi } from "../../hooks/useFeastCoreApi";
import { ComingSoonPanel } from "./components/ComingSoonPanel";

const Project = () => {
  const { apps } = useContext(ApplicationsContext);
  const { projectId, project } = useContext(CurrentProjectContext);

  const [{ data: entities }, fetchEntities] = useFeastCoreApi(
    `/feast.core.CoreService/ListEntities`,
    { method: "POST" },
    undefined,
    false
  );
  const [{ data: featureTables }, fetchFeatureTables] = useFeastCoreApi(
    `/feast.core.CoreService/ListFeatureTables`,
    { method: "POST" },
    undefined,
    false
  );
  const [
    { data: feastIngestionJobs },
    fetchFeastIngestionJobs
  ] = useFeastCoreApi(
    `/feast_spark.api.JobService/ListJobs`,
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
    if (project) {
      fetchEntities({
        body: JSON.stringify({ filter: { project: project.name } })
      });
      fetchFeatureTables({
        body: JSON.stringify({ filter: { project: project.name } })
      });
      fetchFeastIngestionJobs({
        body: JSON.stringify({
          include_terminated: true,
          project: project.name.replace(/-/g, "_")
        })
      });
      fetchModels();
      fetchRouters();
    }
  }, [
    project,
    fetchEntities,
    fetchFeatureTables,
    fetchFeastIngestionJobs,
    fetchModels,
    fetchRouters
  ]);

  return (
    <EuiPage>
      <EuiPageBody>
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
                {/* <InstancesSummary project={project} models={models} /> */}
                <ComingSoonPanel title="Monthly Cost" iconType="visPie" />
              </EuiFlexItem>
              <EuiFlexItem grow={1}></EuiFlexItem>
            </EuiFlexGroup>

            <EuiSpacer size="l" />

            <EuiFlexGroup>
              <EuiFlexItem grow={6}>
                <Instances
                  project={project}
                  feastIngestionJobs={feastIngestionJobs}
                  models={models}
                  routers={routers}
                />
              </EuiFlexItem>
              <EuiFlexItem grow={4}></EuiFlexItem>
            </EuiFlexGroup>

            <EuiSpacer size="l" />

            <EuiFlexGroup>
              <EuiFlexItem grow={5}>
                <ComingSoonPanel
                  title="Health Monitoring"
                  iconType="monitoringApp"
                />
              </EuiFlexItem>
              <EuiFlexItem grow={4}>
                <ComingSoonPanel title="Error Summary" iconType="bug" />
              </EuiFlexItem>
              <EuiFlexItem grow={1}></EuiFlexItem>
            </EuiFlexGroup>
          </>
        ) : (
          <EuiLoadingContent lines={5} />
        )}
      </EuiPageBody>
    </EuiPage>
  );
};

export default Project;
