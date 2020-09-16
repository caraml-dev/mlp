import React, { useCallback, useEffect, useState } from "react";
import {
  EuiButton,
  EuiButtonEmpty,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiForm,
  EuiFormRow,
  EuiIcon,
  EuiPanel,
  EuiSpacer,
  EuiText,
  EuiTitle,
  EuiToolTip
} from "@elastic/eui";
import { validateEmail } from "../../validation/validateEmail";
import { addToast, useMlpApi } from "@mlp/ui";
import UserRoleSelection from "./UserRoleSelection";
import { ROLE_OPTIONS } from "./UserRoleListTable";

const SubmitUserRoleForm = ({ userRole, project, fetchUpdates, toggleAdd }) => {
  const [request, setRequest] = useState({
    roles: userRole ? userRole.roles : [],
    user: userRole ? userRole.user : ""
  });

  const [submissionResponse, submitForm] = useMlpApi(
    `/projects/${project.id}`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" }
    },
    {},
    false
  );

  useEffect(() => {
    if (submissionResponse.isLoaded && !submissionResponse.error) {
      addToast({
        id: "submit-success-create",
        title: userRole ? "User Role changed!" : "New User Role is created!",
        color: "success",
        iconType: "check"
      });
      fetchUpdates();

      if (!userRole) {
        toggleAdd();
      }
    }
  }, [submissionResponse, fetchUpdates, toggleAdd, userRole]);

  const onChange = field => {
    return value => {
      setRequest(r => ({ ...r, [field]: value }));
    };
  };

  const [isValidUser, setValidUser] = useState(validateEmail(request.user));

  const onUserChanges = e => {
    setValidUser(validateEmail(e.target.value));
    onChange("user")(e.target.value);
  };

  const saveAction = () => {
    submitForm({
      body: JSON.stringify(convertRequestToJSONPayload(request, project))
    });
  };

  const onRolesChanges = useCallback(onChange("roles"), []);

  return (
    <EuiPanel paddingSize="m">
      <EuiFlexGroup justifyContent="spaceAround" direction="column">
        <EuiFlexItem>
          <EuiTitle size="xs">
            <h5> {userRole ? "Edit" : "Add a"} User Role</h5>
          </EuiTitle>
        </EuiFlexItem>
        <EuiFlexItem>
          <EuiForm
            isInvalid={!!submissionResponse.error}
            error={
              submissionResponse.error ? [submissionResponse.error.message] : ""
            }>
            <EuiFlexGroup direction="column">
              <EuiFlexItem grow={false}>
                {!userRole && (
                  <EuiFormRow
                    fullWidth
                    label={
                      <EuiToolTip content="Specify user for project role">
                        <span>
                          User{" "}
                          <EuiIcon type="questionInCircle" color="subdued" />
                        </span>
                      </EuiToolTip>
                    }
                    display="columnCompressed">
                    <EuiFieldText
                      placeholder="e.g system@google.com"
                      value={request.user}
                      onChange={e => onUserChanges(e)}
                      name="user"
                      isInvalid={!isValidUser}
                    />
                  </EuiFormRow>
                )}
                <EuiSpacer size="s" />
                <UserRoleSelection
                  roleOptions={ROLE_OPTIONS}
                  chosenRoles={request.roles}
                  onChange={onRolesChanges}
                />
              </EuiFlexItem>
            </EuiFlexGroup>
          </EuiForm>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiFlexGroup direction="row">
            <EuiFlexItem grow={false}>
              <EuiButton
                size="s"
                color="primary"
                fill
                onClick={() => saveAction()}
                disabled={submissionResponse.isLoading || !isValidUser}>
                <EuiText size="s">Save</EuiText>
              </EuiButton>
            </EuiFlexItem>
            {!userRole && (
              <EuiFlexItem grow={false}>
                <EuiButtonEmpty
                  size="s"
                  onClick={() => toggleAdd()}
                  disabled={submissionResponse.isLoading}>
                  <EuiText size="s">Cancel</EuiText>
                </EuiButtonEmpty>
              </EuiFlexItem>
            )}
          </EuiFlexGroup>
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPanel>
  );
};

export default SubmitUserRoleForm;

function convertRequestToJSONPayload(request, project) {
  let copyOfProject = Object.assign({}, project);
  const requestUser = request.user;
  const requestRoles = request.roles;
  const projectRoleKeyLookup = {
    Administrator: "administrators",
    Reader: "readers"
  };

  const userRoleMapping = ROLE_OPTIONS.reduce((map, obj) => {
    const roleKey = projectRoleKeyLookup[obj] || "";
    if (roleKey === "") {
      return map;
    }

    const found = requestRoles.indexOf(obj) >= 0;
    map[roleKey] = found ? requestUser : null;
    return map;
  }, {});

  Object.keys(userRoleMapping).forEach(fieldRole => {
    let roleUsers = Object.assign([], copyOfProject[fieldRole]);
    const userWillAssignToRole = userRoleMapping[fieldRole] != null;

    if (userWillAssignToRole) {
      roleUsers = roleUsers.concat(requestUser);
      roleUsers = Array.from(new Set(roleUsers));
    } else {
      const indexOfUser = roleUsers.indexOf(requestUser);

      if (indexOfUser >= 0) {
        roleUsers.splice(indexOfUser, 1);
      }
    }
    copyOfProject[fieldRole] = roleUsers;
  });

  return {
    name: copyOfProject.name,
    administrators: copyOfProject.administrators,
    readers: copyOfProject.readers,
    stream: copyOfProject.stream,
    team: copyOfProject.team
  };
}
