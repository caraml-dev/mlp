import React, { useEffect } from "react";
import { replaceBreadcrumbs, useMlpApi } from "@mlp/ui";
import { EuiCallOut, EuiLoadingChart, EuiTextAlign } from "@elastic/eui";
import { SettingsSection } from "./ProjectSetting";
import SecretListTable from "./secret/SecretListTable";

const SecretSetting = ({ projectId }) => {
  const [{ data, isLoaded, error }, fetchSecrets] = useMlpApi(
    `/projects/${projectId}/secrets`
  );

  useEffect(() => {
    replaceBreadcrumbs([
      {
        text: "Secrets Management"
      }
    ]);
  }, [projectId]);

  return (
    <SettingsSection title="Secrets">
      {!isLoaded ? (
        <EuiTextAlign textAlign="center">
          <EuiLoadingChart size="xl" mono />
        </EuiTextAlign>
      ) : error ? (
        <EuiCallOut
          title="Sorry, there was an error"
          color="danger"
          iconType="alert">
          <p>{error.message}</p>
        </EuiCallOut>
      ) : (
        <SecretListTable
          secrets={data}
          projectId={projectId}
          fetchUpdates={() => fetchSecrets()}
        />
      )}
    </SettingsSection>
  );
};

export default SecretSetting;
