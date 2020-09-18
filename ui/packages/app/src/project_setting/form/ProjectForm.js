import React, { useContext, useState, useEffect } from "react";
import { navigate } from "@reach/router";
import {
  EuiPanel,
  EuiFormRow,
  EuiFieldText,
  EuiDescribedFormGroup,
  EuiFlexGroup,
  EuiFlexItem,
  EuiTitle,
  EuiButton,
  EuiForm,
  EuiIcon
} from "@elastic/eui";
import { addToast, useMlpApi } from "@gojek/mlp-ui";
import { ProjectFormContext } from "./context";
import { SingleSelectionComboBox } from "./SingleSelectionComboBox";
import { EmailTextArea } from "./EmailTextArea";
import { Labels } from "./Labels";
import { validateSubdomain } from "../../validation/validateSubdomain";

const ProjectForm = () => {
  const {
    project,
    setName,
    setTeam,
    setStream,
    setAdmin,
    setReader,
    setLabels
  } = useContext(ProjectFormContext);
  const [isValidProjectName, setValidProjectName] = useState(
    project.name ? validateSubdomain(project.name) : true
  );

  const onProjectNameChange = e => {
    const newValue = e.target.value;
    setValidProjectName(validateSubdomain(newValue));
    setName(newValue);
  };

  const teamOptions = [];
  const onTeamChange = selectedTeam => {
    setTeam(selectedTeam.label);
  };
  const [teamValid, setTeamValid] = useState();

  const streamOptions = [];
  const onStreamChange = selectedStream => {
    setStream(selectedStream.label);
  };
  const [streamValid, setStreamValid] = useState();

  const onAdminValueChange = emails => {
    setAdmin(emails);
  };
  const [adminValid, setAdminValid] = useState();

  const onReaderValueChange = emails => {
    setReader(emails);
  };
  const [readerValid, setReaderValid] = useState();

  const onLabelChange = labels => {
    const labelsValid =
      labels.length === 0
        ? true
        : labels.reduce((labelsValid, label) => {
            return labelsValid && label.isKeyValid && label.isValueValid;
          }, true);
    setLabelsValid(labelsValid);

    //deep copy
    let newLabels = JSON.parse(JSON.stringify(labels));
    newLabels = newLabels.map(element => {
      delete element.isKeyValid;
      delete element.isValueValid;
      delete element.idx;
      return element;
    });

    setLabels(newLabels);
  };
  const [labelsValid, setLabelsValid] = useState(true);

  const onSubmit = () => {
    submitForm({ body: JSON.stringify(project) });
  };
  const [submissionResponse, submitForm] = useMlpApi(
    "/projects",
    {
      method: "POST",
      headers: { "Content-Type": "application/json" }
    },
    {},
    false
  );

  useEffect(() => {
    if (submissionResponse.isLoaded && !submissionResponse.error) {
      addToast({
        id: "create-project-success-toast",
        title: "Project Created",
        color: "success",
        iconType: "check"
      });
      navigate(`/projects/${submissionResponse.data.id}`);
    }
  }, [submissionResponse]);

  return (
    <EuiForm>
      <EuiFlexGroup>
        <EuiFlexItem grow={false}>
          <EuiTitle size="m">
            <h1>
              <EuiIcon type="folderClosed" size="xl" />
              &nbsp; Create New Project
            </h1>
          </EuiTitle>
        </EuiFlexItem>
      </EuiFlexGroup>

      <EuiFlexGroup direction="column">
        <EuiFlexItem grow={false}>
          <EuiPanel grow={false}>
            <EuiDescribedFormGroup
              title={<h3>Name</h3>}
              description="Project name should be a valid subdomain">
              <EuiFormRow>
                <EuiFieldText
                  name="name"
                  placeholder="my-new-project"
                  value={project.name}
                  onChange={onProjectNameChange}
                  isInvalid={!isValidProjectName}
                  aria-label="Project Name"
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Team</h3>}
              description="Owner of the project">
              <EuiFormRow>
                <SingleSelectionComboBox
                  options={teamOptions}
                  onChange={onTeamChange}
                  onValidChange={setTeamValid}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Stream</h3>}
              description="Product stream the project belongs to">
              <EuiFormRow>
                <SingleSelectionComboBox
                  options={streamOptions}
                  onChange={onStreamChange}
                  onValidChange={setStreamValid}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Project Members</h3>}
              description="Comma separated list of user / service account email. Administrators have full access to the project, whereas Readers have read-only access.">
              <EuiFormRow label="Administrators">
                <EmailTextArea
                  onChange={onAdminValueChange}
                  onValidChange={setAdminValid}
                />
              </EuiFormRow>
              <EuiFormRow label="Readers">
                <EmailTextArea
                  onChange={onReaderValueChange}
                  onValidChange={setReaderValid}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Labels</h3>}
              description="Additional Labels">
              <Labels onChange={onLabelChange} />
            </EuiDescribedFormGroup>
          </EuiPanel>
        </EuiFlexItem>

        <EuiFlexItem grow={false}>
          <EuiFlexGroup direction="row" justifyContent="flexEnd">
            <EuiFlexItem grow={false}>
              <EuiButton
                size="s"
                color="primary"
                onClick={onSubmit}
                disabled={
                  !(
                    isValidProjectName &&
                    teamValid &&
                    streamValid &&
                    adminValid &&
                    readerValid &&
                    labelsValid
                  )
                }
                fill>
                Submit
              </EuiButton>
            </EuiFlexItem>
          </EuiFlexGroup>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiForm>
  );
};

export default ProjectForm;
