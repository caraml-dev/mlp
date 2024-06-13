import React, { useEffect } from "react";
import { addToast, useMlpApi } from "@caraml-dev/ui-lib";
import {
  EuiConfirmModal,
  EuiOverlayMask,
  EuiCode,
  EuiBasicTable,
  EuiFlexItem,
  EuiFlexGroup,
  EuiLink
} from "@elastic/eui";
import { useNavigate } from "react-router-dom";

const UpdateProjectInfoModal = ({
  project,
  closeModal,
  fetchUpdates,
  originalProject
}) => {
  const [submissionResponse, submitForm] = useMlpApi(
    `/v1/projects/${project.id}`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" }
    },
    {},
    false
  );

  const navigate = useNavigate();

  useEffect(() => {
    if (submissionResponse.isLoaded && !submissionResponse.error) {
      const response = submissionResponse.data.update_status_url;
      closeModal();
      addToast({
        id: "submit-success-create",
        title: "Project Info Updated!",
        color: "success",
        iconType: "check",
        text: response ? (
          <p>
            Your project's resource(s) are currently being redeployed, you can
            check the workflow here:{" "}
            <EuiLink href={response} target="_blank">
              Link
            </EuiLink>
          </p>
        ) : null
      });
      fetchUpdates();
      navigate(`/projects/${project.id}/settings/project-info`);
    }
  }, [navigate, submissionResponse, project, fetchUpdates, closeModal]);

  const handleUpdate = () => {
    const updatedProjectInfo = {
      ...project,
      stream: project.stream,
      team: project.team,
      labels: project.labels
    };

    submitForm({ body: JSON.stringify(updatedProjectInfo) });
  };

  const columns = [
    {
      field: "key",
      name: "Key"
    },
    {
      field: "value",
      name: "Value"
    }
  ];

  const rows = project.labels.map(label => ({
    ...label
  }));

  const originalRows = originalProject.labels.map(label => ({
    ...label
  }));

  return (
    <EuiOverlayMask>
      <EuiConfirmModal
        title="Update Project Info"
        onCancel={() => closeModal()}
        onConfirm={handleUpdate}
        cancelButtonText="Cancel"
        confirmButtonText="Update"
        buttonColor="primary"
        defaultFocusedButton="confirm">
        <p>
          You are about to update the project info for the project{" "}
          <b>{project.name}</b>
          , are you sure?
          <br />
          <b>Note:</b> Project info changes will take approximately 10 minutes.
        </p>
        <EuiFlexGroup gutterSize="l">
          <EuiFlexItem>
            <h3>Current Project Info</h3>
            <p>
              Stream: <EuiCode>{originalProject.stream}</EuiCode>
              <br />
              Team: <EuiCode>{originalProject.team}</EuiCode>
              <br />
              Labels:
              <EuiBasicTable
                items={originalRows}
                columns={columns}
                rowHeader="Key"
              />
            </p>
          </EuiFlexItem>
          <EuiFlexItem>
            <h3>New Project Info</h3>
            <p>
              Stream: <EuiCode>{project.stream}</EuiCode>
              <br />
              Team: <EuiCode>{project.team}</EuiCode>
              <br />
              Labels:
              <EuiBasicTable items={rows} columns={columns} rowHeader="Key" />
            </p>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiConfirmModal>
    </EuiOverlayMask>
  );
};

export default UpdateProjectInfoModal;
