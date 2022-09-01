import React, { useContext, useEffect } from "react";
import { EuiFlexGroup, EuiFlexItem, EuiPageTemplate } from "@elastic/eui";
import { GoogleLogin } from "@react-oauth/google";
import { get } from "../utils";
import AuthContext from "../auth/context";
import { useLocation, useNavigate } from "react-router-dom";

export const Login = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const {
    state: { isAuthenticated },
    onLogin
  } = useContext(AuthContext);

  useEffect(() => {
    if (isAuthenticated) {
      navigate(get(location, "state.referer") || "/");
    }
  }, [isAuthenticated, location, navigate]);

  const onFailure = () => {
    console.log("Login Failed");
  };

  return (
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
