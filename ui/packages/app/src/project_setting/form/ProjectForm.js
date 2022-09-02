import React, { useContext, useState, useEffect } from "react";
import {
  EuiPanel,
  EuiFormRow,
  EuiFieldText,
  EuiDescribedFormGroup,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButton,
  EuiForm
} from "@elastic/eui";
import { addToast, useMlpApi } from "@gojek/mlp-ui";
import { ProjectFormContext } from "./context";
import { SingleSelectionComboBox } from "./SingleSelectionComboBox";
import { EmailTextArea } from "./EmailTextArea";
import { Labels } from "./Labels";
import { isDNS1123Label } from "../../validation/validation";
import config from "../../config";
import { useNavigate } from "react-router-dom";

const ProjectForm = () => {
  const navigate = useNavigate();

  const {
    project,
    setName,
    setTeam,
    setStream,
    setAdmin,
    setReader,
    setLabels
  } = useContext(ProjectFormContext);

  const teamOptions = config.TEAMS.map(team => {
    return {
      label: team
    };
  });

  const streamOptions = config.STREAMS.map(stream => {
    return {
      label: stream
    };
  });

  const [projectError, setProjectError] = useState("");
  const [isValidProject, setIsValidProject] = useState(false);
  const onProjectChange = e => {
    const newValue = e.target.value;
    let isValid = isDNS1123Label(newValue);
    if (!isValid) {
      setProjectError(
        "Project name is invalid. It should contain only lowercase alphanumeric and dash ('-')"
      );
    }
    setIsValidProject(isValid);
    setName(newValue);
  };

  const [teamError, setTeamError] = useState("");
  const [isValidTeam, setIsValidTeam] = useState(false);
  const onTeamChange = selectedTeam => {
    let isValid = isDNS1123Label(selectedTeam.label);
    if (!isValid) {
      setTeamError(
        "Team name is invalid. It should contain only lowercase alphanumeric and dash (-)"
      );
    }
    setIsValidTeam(isValid);
    setTeam(selectedTeam.label);
  };

  const [streamError, setStreamError] = useState("");
  const [isValidStream, setIsValidStream] = useState(false);
  const onStreamChange = selectedStream => {
    let isValid = isDNS1123Label(selectedStream.label);
    if (!isValid) {
      setStreamError(
        "Stream name is invalid. It should contain only lowercase alphanumeric and dash (-)"
      );
    }
    setIsValidStream(isValid);
    setStream(selectedStream.label);
  };

  const onAdminValueChange = emails => {
    setAdmin(emails);
  };
  const [isValidAdmin, setIsValidAdmin] = useState(false);
  const [adminError, setAdminError] = useState("");
  const onValidAdminChange = valid => {
    setIsValidAdmin(valid);
    if (!valid) {
      setAdminError("Invalid email address");
    }
  };

  const onReaderValueChange = emails => {
    setReader(emails);
  };
  const [isValidReader, setIsValidReader] = useState(false);
  const [readerError, setReaderError] = useState("");
  const onValidReaderChange = valid => {
    setIsValidReader(valid);
    if (!valid) {
      setReaderError("Invalid email address");
    }
  };

  const [isValidLabels, setIsValidLabels] = useState(true);
  const [labelError, setLabelError] = useState("");
  const onLabelChange = labels => {
    const labelsValid =
      labels.length === 0
        ? true
        : labels.reduce((labelsValid, label) => {
            return labelsValid && label.isKeyValid && label.isValueValid;
          }, true);
    setIsValidLabels(labelsValid);
    if (!labelsValid) {
      setLabelError(
        "Invalid labels. Both key and value of a label must contain only alphanumeric and dash (-)"
      );
    }

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
  }, [navigate, submissionResponse]);

  return (
    <EuiForm>
      <EuiFlexGroup direction="column">
        <EuiFlexItem grow={false}>
          <EuiPanel grow={false}>
            <EuiDescribedFormGroup
              title={<h3>Name</h3>}
              description="Project name should contain only lowercase alphanumeric and dash ('-')">
              <EuiFormRow isInvalid={!isValidProject} error={projectError}>
                <EuiFieldText
                  name="name"
                  placeholder="my-new-project"
                  value={project.name}
                  onChange={onProjectChange}
                  isInvalid={!isValidProject}
                  aria-label="Project Name"
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Team</h3>}
              description="Owner of the project">
              <EuiFormRow isInvalid={!isValidTeam} error={teamError}>
                <SingleSelectionComboBox
                  options={teamOptions}
                  onChange={onTeamChange}
                  onValidChange={setIsValidTeam}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Stream</h3>}
              description="Product stream the project belongs to">
              <EuiFormRow isInvalid={!isValidStream} error={streamError}>
                <SingleSelectionComboBox
                  options={streamOptions}
                  onChange={onStreamChange}
                  onValidChange={setIsValidStream}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Project Members</h3>}
              description="Comma separated list of user / service account email. Administrators have full access to the project, whereas Readers have read-only access.">
              <EuiFormRow
                label="Administrators"
                isInvalid={!isValidAdmin}
                error={adminError}>
                <EmailTextArea
                  onChange={onAdminValueChange}
                  onValidChange={onValidAdminChange}
                />
              </EuiFormRow>
              <EuiFormRow
                label="Readers"
                isInvalid={!isValidReader}
                error={readerError}>
                <EmailTextArea
                  onChange={onReaderValueChange}
                  onValidChange={onValidReaderChange}
                />
              </EuiFormRow>
            </EuiDescribedFormGroup>
            <EuiDescribedFormGroup
              title={<h3>Labels</h3>}
              description="Additional Labels">
              <EuiFormRow isInvalid={!isValidLabels} error={labelError}>
                <Labels onChange={onLabelChange} />
              </EuiFormRow>
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
                    isValidProject &&
                    isValidTeam &&
                    isValidStream &&
                    isValidAdmin &&
                    isValidReader &&
                    isValidLabels
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
