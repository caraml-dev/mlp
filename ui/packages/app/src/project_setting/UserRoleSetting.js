import React, { useContext } from "react";
import UserRoleListTable from "./user_role/UserRoleListTable";
import { CurrentProjectContext } from "@gojek/mlp-ui";
import { EuiLoadingChart, EuiTextAlign } from "@elastic/eui";

const UserRoleSetting = () => {
  const { project, refresh } = useContext(CurrentProjectContext);

  return (
    <>
      {!project ? (
        <EuiTextAlign textAlign="center">
          <EuiLoadingChart size="xl" mono />
        </EuiTextAlign>
      ) : (
        <UserRoleListTable project={project} fetchUpdates={() => refresh()} />
      )}
    </>
  );
};

export default UserRoleSetting;
