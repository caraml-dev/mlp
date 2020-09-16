import React, { useState } from "react";
import { EuiButton, EuiText, EuiFlexGroup } from "@elastic/eui";

import SubmitUserRoleForm from "./SubmitUserRoleForm";

const UserRoleForm = ({ project, fetchUpdates }) => {
  const [isAdd, setAdd] = useState(false);
  return isAdd ? (
    <SubmitUserRoleForm
      project={project}
      fetchUpdates={fetchUpdates}
      toggleAdd={() => setAdd(!isAdd)}
    />
  ) : (
    <EuiFlexGroup gutterSize="xs">
      <EuiButton fill size="s" onClick={() => setAdd(!isAdd)}>
        <EuiText size="s">Add User</EuiText>
      </EuiButton>
    </EuiFlexGroup>
  );
};
export default UserRoleForm;
