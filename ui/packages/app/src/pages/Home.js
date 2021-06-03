import React, { useContext } from "react";
import { EuiLoadingContent, EuiPage, EuiPageBody } from "@elastic/eui";
import { CurrentProjectContext } from "@gojek/mlp-ui";
import { Redirect } from "@reach/router";

const Home = () => {
  const { projectId } = useContext(CurrentProjectContext);

  return projectId === null ? (
    <EuiPage>
      <EuiPageBody>
        <EuiLoadingContent lines={5} />
      </EuiPageBody>
    </EuiPage>
  ) : (
    <Redirect to={`/projects/${projectId}`} noThrow />
  );
};

export default Home;
