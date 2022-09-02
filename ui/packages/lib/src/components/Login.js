import React, { useContext } from "react";
import { EuiFlexGroup, EuiFlexItem, EuiPageTemplate } from "@elastic/eui";
import { GoogleLogin } from "@react-oauth/google";
import AuthContext from "../auth/context";
import { Navigate, useLocation } from "react-router-dom";

export const Login = () => {
  const location = useLocation();
  const { state, onLogin } = useContext(AuthContext);

  const onFailure = () => {
    console.log("Login Failed");
  };

  return !!state.isAuthenticated ? (
    <Navigate to={location?.state?.referer || "/"} replace={true} />
  ) : (
    <EuiPageTemplate>
      <EuiPageTemplate.EmptyPrompt
        iconType="machineLearningApp"
        title={<h2>Machine Learning Platform</h2>}
        body={<p>Use your Google account to sign in</p>}
        actions={
          <EuiFlexGroup
            direction="column"
            alignItems="center"
            gutterSize="none">
            <EuiFlexItem grow={false} style={{ maxWidth: "200px" }}>
              <GoogleLogin
                onSuccess={onLogin}
                onError={onFailure}
                useOneTap={true}
              />
            </EuiFlexItem>
          </EuiFlexGroup>
        }
      />
    </EuiPageTemplate>
  );
};
