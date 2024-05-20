import React, { useState } from "react";
import { EuiButton, EuiText, EuiFlexGroup, EuiFlexItem } from "@elastic/eui";
import UpdateProjectInfoModal from "./UpdateProjectInfoModal";

const SubmitProjectInfoForm = ({
  project,
  isValidTeam,
  isValidStream,
  isValidLabels,
  fetchUpdates,
  isDisabled,
  originalProject
}) => {
  const [showUpdateModal, setShowUpdateModal] = useState(false);

  return (
    <EuiFlexGroup gutterSize="s">
      <EuiFlexItem grow={false}>
        <EuiButton
          size="s"
          color="primary"
          onClick={() => setShowUpdateModal(true)}
          disabled={
            isDisabled || !(isValidTeam && isValidStream && isValidLabels)
          }
          fill>
          <EuiText size="s">Submit</EuiText>
        </EuiButton>
        {showUpdateModal && (
          <UpdateProjectInfoModal
            project={project}
            closeModal={() => setShowUpdateModal(false)}
            fetchUpdates={fetchUpdates}
            originalProject={originalProject}
          />
        )}
      </EuiFlexItem>
    </EuiFlexGroup>
  );
};

export default SubmitProjectInfoForm;
