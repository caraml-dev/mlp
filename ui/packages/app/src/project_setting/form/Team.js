import React, { useState, useEffect, useMemo } from "react";
import { EuiFormRow } from "@elastic/eui";
import { EuiComboBoxSelect } from "@caraml-dev/ui-lib";
import { isValidK8sLabelKeyValue } from "../../validation/validation";
import config from "../../config";

export const Team = ({
  team,
  setTeam,
  stream,
  isValidTeam,
  setIsValidTeam,
  isDisabled = false
}) => {
  const teamOptions = useMemo(() => {
    return (config.STREAMS[stream] || [])
      .sort((a, b) => a.localeCompare(b))
      .map(team => ({ label: team.trim() }));
  }, [stream]);

  const [teamError, setTeamError] = useState("");

  const onTeamChange = team => {
    let isValid = isValidK8sLabelKeyValue(team);
    if (!isValid) {
      setTeamError(
        "Team name is invalid. It should contain only lowercase alphanumeric and dash (-) or underscore (_) or period (.), and must start and end with an alphanumeric character"
      );
    }
    setIsValidTeam(isValid);
    setTeam(team);
  };

  useEffect(() => {
    if (!team) {
      setIsValidTeam(false);
    }
  }, [team, setIsValidTeam]);

  return (
    <EuiFormRow isInvalid={!isValidTeam} error={teamError}>
      <EuiComboBoxSelect
        value={team}
        options={teamOptions}
        onChange={onTeamChange}
        onCreateOption={config.ALLOW_CUSTOM_TEAM ? onTeamChange : undefined}
        isDiasbled={isDisabled}
      />
    </EuiFormRow>
  );
};

export default Team;
