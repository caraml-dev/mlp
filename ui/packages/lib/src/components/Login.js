import React, { Fragment, useContext } from "react";
import {
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiFlexGroup,
  EuiFlexItem,
  EuiEmptyPrompt
} from "@elastic/eui";
import { GoogleLogin } from "@react-oauth/google";
import { Redirect } from "@reach/router";
import { get } from "../utils";
import AuthContext from "../auth/context";

export const Login = ({ location }) => {
  const {
    state: { isAuthenticated },
    onLogin
  } = useContext(AuthContext);

  const onFailure = () => {
    console.log("Login Failed");
  };

  return isAuthenticated ? (
    <Redirect to={get(location, "state.referer") || "/"} noThrow />
  ) : (
    <EuiPage className="eui-fullHeight">
      <EuiPageBody>
        <EuiPageContent verticalPosition="center" horizontalPosition="center">
          <EuiEmptyPrompt
            iconType="machineLearningApp"
            title={<h2>Machine Learning Platform</h2>}
            body={
              <Fragment>
                <p>Use your Google account to sign in</p>
              </Fragment>
            }
          />
          <EuiFlexGroup direction="column" alignItems="center">
            <EuiFlexItem grow={false} style={{ maxWidth: "200px" }}>
              <GoogleLogin
                onSuccess={onLogin}
                onError={onFailure}
                useOneTap={true}
              />
            </EuiFlexItem>
          </EuiFlexGroup>
        </EuiPageContent>
      </EuiPageBody>
    </EuiPage>
  );
};
