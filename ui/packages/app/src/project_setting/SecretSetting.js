import React from "react";
import { useMlpApi } from "@gojek/mlp-ui";
import { EuiCallOut, EuiLoadingChart, EuiTextAlign } from "@elastic/eui";
import SecretListTable from "./secret/SecretListTable";
import { useParams } from "react-router-dom";

const SecretSetting = () => {
  const { projectId } = useParams();
  const [{ data, isLoaded, error }, fetchSecrets] = useMlpApi(
    `/projects/${projectId}/secrets`
  );

  return (
    <>
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
          fetchUpdates={fetchSecrets}
        />
      )}
    </>
  );
};

export default SecretSetting;
