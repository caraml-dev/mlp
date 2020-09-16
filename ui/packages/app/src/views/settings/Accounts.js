import React, { Fragment, useContext, useEffect } from "react";
import { EuiButton, EuiText } from "@elastic/eui";
import gitlabSvg from "./gitlab.svg";
import { AuthContext, useMlpApi } from "@mlp/ui";

export const Accounts = props => {
  const {
    state: { profileObj }
  } = useContext(AuthContext);
  const gitlabSessionToken = sessionStorage.getItem("gitlab");

  const [gitlabAuth] = useMlpApi(`/users/authorize`, {}, "", true);
  const [newToken, generateToken] = useMlpApi(
    `/users/token/generate?code=${props.search.code}&state=${props.search.state}`,
    {},
    "",
    false
  );
  const [savedToken, retrieveToken] = useMlpApi(
    `/users/token/retrieve`,
    {
      headers: {
        "User-Email": profileObj.email
      }
    },
    "",
    false
  );

  useEffect(() => {
    if (
      typeof props.search.code !== "undefined" &&
      typeof props.search.state !== "undefined" &&
      gitlabAuth.isLoaded &&
      !newToken.isLoading &&
      !newToken.isLoaded
    ) {
      generateToken();
    }
  }, [gitlabAuth, props.search, newToken, generateToken]);

  useEffect(() => {
    if (newToken.isLoaded && !savedToken.isLoading && !savedToken.isLoaded) {
      retrieveToken();
    }
  }, [newToken, savedToken, retrieveToken]);

  useEffect(() => {
    if (savedToken.isLoaded && !savedToken.error) {
      sessionStorage.setItem("gitlab", savedToken.data);
    }
  }, [savedToken]);

  return (
    <Fragment>
      <EuiText>
        <h4>Connected Accounts</h4>
        <hr />
        <br />
      </EuiText>
      {!gitlabSessionToken && (
        <EuiButton
          iconType={gitlabSvg}
          color="warning"
          iconSide="left"
          size="s"
          href={gitlabAuth.data}
          disabled={!gitlabAuth.isLoaded}>
          Authorize Gitlab
        </EuiButton>
      )}
      {gitlabSessionToken && (
        <EuiButton iconType="check" iconSide="left" size="s" disabled>
          Connected to Gitlab
        </EuiButton>
      )}
    </Fragment>
  );
};

export default Accounts;
