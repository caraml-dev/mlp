import React, { Fragment, useEffect, useState } from "react";
import {
  EuiButtonEmpty,
  EuiButtonIcon,
  EuiFlexGroup,
  EuiFlexItem,
  EuiForm,
  EuiFormRow,
  EuiSuperSelect,
  EuiText,
  EuiTextColor,
  EuiToolTip
} from "@elastic/eui";

const UserRoleSelection = ({ chosenRoles = [], roleOptions, onChange }) => {
  const noAvailableRoles = chosenRoles.length === roleOptions.length;
  chosenRoles = noAvailableRoles
    ? chosenRoles
    : chosenRoles.length === 0
    ? [...chosenRoles, ""]
    : [...chosenRoles];
  const [selectedRoles, setSelectedRoles] = useState(chosenRoles);

  const removeRoleOption = idx => {
    selectedRoles.splice(idx, 1);
    setSelectedRoles(selectedRoles.slice());
  };

  const [enableAddRole, setEnableAddRole] = useState(!noAvailableRoles);
  useEffect(() => {
    const enable =
      selectedRoles[selectedRoles.length - 1] !== "" && !noAvailableRoles;
    setEnableAddRole(enable);
  }, [selectedRoles, noAvailableRoles]);

  const addRoleOption = _ => {
    setSelectedRoles([...selectedRoles, ""]);
  };

  useEffect(() => {
    onChange(selectedRoles);
  }, [selectedRoles, onChange]);

  const roleOptionsDropDown = () => {
    return roleOptions.map(opt => {
      const isDisabled = selectedRoles.indexOf(opt) >= 0;
      return {
        value: opt,
        inputDisplay: opt,
        disabled: isDisabled,
        dropdownDisplay: <RoleDropdownOption role={opt} disabled={isDisabled} />
      };
    });
  };

  const onChangeRow = idx => {
    return e => {
      selectedRoles[idx] = e;
      setSelectedRoles(selectedRoles.slice());
    };
  };

  const roleSelection = selectedRoles.map((role, idx) => (
    <EuiFormRow
      fullWidth
      display="columnCompressed"
      key={role + "-" + idx}
      label="Select role">
      <EuiFlexGroup direction="row" alignItems="center">
        <EuiFlexItem grow={false} style={{ minWidth: 200 }}>
          <EuiSuperSelect
            options={roleOptionsDropDown()}
            valueOfSelected={role || ""}
            hasDividers
            onChange={onChangeRow(idx)}
          />
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiButtonIcon
            size="s"
            color="danger"
            iconType="trash"
            onClick={() => removeRoleOption(idx)}
            aria-label="Remove variable"
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiFormRow>
  ));
  return (
    <EuiForm>
      {selectedRoles && roleSelection}
      <EuiFormRow fullWidth>
        <EuiButtonEmpty
          onClick={() => addRoleOption()}
          size="s"
          color="text"
          iconType="plusInCircle"
          style={{ marginLeft: -12 }}
          disabled={!enableAddRole}>
          <EuiText size="s">Add another role</EuiText>
        </EuiButtonEmpty>
      </EuiFormRow>
    </EuiForm>
  );
};

const RoleDropdownOption = ({ role, disabled }) => {
  const option = (
    <Fragment>
      <EuiText size="s">
        <EuiTextColor color="subdued">{role}</EuiTextColor>
      </EuiText>
    </Fragment>
  );

  return disabled ? (
    <EuiToolTip position="left" content={`${role} already selected`}>
      {option}
    </EuiToolTip>
  ) : (
    option
  );
};

export default UserRoleSelection;
