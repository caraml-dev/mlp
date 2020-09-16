import React, { useState } from "react";
import { EuiButton, EuiFlexGroup, EuiText } from "@elastic/eui";

import SubmitSecretForm from "./SubmitSecretForm";

const SecretForm = ({ projectId, fetchUpdates }) => {
  const [isAdd, setAdd] = useState(false);

  return isAdd ? (
    <SubmitSecretForm
      projectId={projectId}
      fetchUpdates={fetchUpdates}
      toggleAdd={() => setAdd(!isAdd)}
    />
  ) : (
    <EuiFlexGroup gutterSize="xs">
      <EuiButton fill size="s" onClick={() => setAdd(!isAdd)}>
        <EuiText size="s">Create New Secret</EuiText>
      </EuiButton>
    </EuiFlexGroup>
  );
};
export default SecretForm;
