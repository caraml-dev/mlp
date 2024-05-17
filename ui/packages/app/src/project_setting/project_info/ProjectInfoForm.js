import React, { useContext, useState } from "react";
import {
  EuiFlexItem,
  EuiText,
  EuiFlexGroup,
  EuiForm,
  EuiSpacer,
  EuiTitle
} from "@elastic/eui";
import SubmitProjectInfoForm from "./SubmitProjectInfoForm";
import config from "../../config";
import { isDNS1123Label } from "../../validation/validation";
import { ProjectFormContext } from "../form/context";
import { Labels } from "../form/Labels";
import { Stream } from "../form/Stream";
import { Team } from "../form/Team";

const ProjectInfoForm = ({ originalProject, fetchUpdates }) => {
  const { project, setStream, setTeam, setLabels } = useContext(
    ProjectFormContext
  );

  const [isValidStream, setIsValidStream] = useState(
    isDNS1123Label(project.stream)
  );
  const [isValidTeam, setIsValidTeam] = useState(isDNS1123Label(project.team));

  const [isValidLabels, setIsValidLabels] = useState(
    project.labels.length === 0
      ? true
      : project.labels.reduce((labelsValid, label) => {
          return (
            labelsValid &&
            isDNS1123Label(label.key) &&
            isDNS1123Label(label.value)
          );
        }, true)
  );

  const isDisabled = !config.PROJECT_INFO_UPDATE_ENABLED;

  return (
    <>
      <EuiFlexGroup direction="column">
        <EuiForm>
          <EuiFlexItem>
            <EuiTitle size="s">
              <h3>Stream</h3>
            </EuiTitle>
            <EuiSpacer size="m" />
            <EuiText size="s" color="subdued">
              <p>Product stream the project belongs to</p>
            </EuiText>
            <EuiSpacer size="m" />
            <Stream
              stream={project.stream}
              setStream={setStream}
              isValidStream={isValidStream}
              setIsValidStream={setIsValidStream}
              isDisabled={isDisabled}
            />
            <EuiSpacer size="xl" />
            <EuiTitle size="s">
              <h3>Team</h3>
            </EuiTitle>
            <EuiSpacer size="m" />
            <EuiText size="s" color="subdued">
              <p>The owner of the project</p>
            </EuiText>
            <EuiSpacer size="m" />
            <Team
              team={project.team}
              setTeam={setTeam}
              stream={project.stream}
              isValidTeam={isValidTeam}
              setIsValidTeam={setIsValidTeam}
              isDisabled={isDisabled}
            />
            <EuiSpacer size="xl" />
            <EuiTitle size="s">
              <h3>Labels</h3>
            </EuiTitle>
            <EuiSpacer size="m" />
            <EuiText size="s" color="subdued">
              <p>Additional Labels</p>
            </EuiText>
            <EuiSpacer size="m" />
            <Labels
              labels={project.labels}
              setLabels={setLabels}
              setIsValidLabels={setIsValidLabels}
              isValidLabels={isValidLabels}
              isDisabled={isDisabled}
            />
          </EuiFlexItem>
          <EuiSpacer size="m" />
          <EuiFlexItem>
            <SubmitProjectInfoForm
              project={project}
              fetchUpdates={fetchUpdates}
              isValidTeam={isValidTeam}
              isValidStream={isValidStream}
              isValidLabels={isValidLabels}
              isDisabled={isDisabled}
              originalProject={originalProject}
            />
          </EuiFlexItem>
        </EuiForm>
      </EuiFlexGroup>
    </>
  );
};

export default ProjectInfoForm;
