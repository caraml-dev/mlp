import { GoogleOAuthProvider } from "@react-oauth/google";
import { jwtDecode } from "jwt-decode";
import PropTypes from "prop-types";
import React, { useCallback, useEffect, useState } from "react";

const emptyState = {
  isAuthenticated: false,
  profileObj: null,
  jwt: ""
};

const AuthContext = React.createContext({
  state: emptyState
});

export const AuthProvider = ({ clientId, children }) => {
  const [state, setState] = useState(emptyState);

  const onJWTUpdate = useCallback(
    jwt => {
      try {
        const profileObj = jwtDecode(jwt);
        if (profileObj.exp * 1000 > Date.now()) {
          setState({
            isAuthenticated: true,
            profileObj: profileObj,
            jwt: jwt
          });
        }
      } catch (error) {
        console.debug("failed to decode JWT: \n", jwt);
      }
    },
    [setState]
  );

  useEffect(() => {
    onJWTUpdate(localStorage.getItem("auth"));
  }, [onJWTUpdate]);

  const login = ({ credential }) => {
    onJWTUpdate(credential);
    localStorage.setItem("auth", credential);
  };

  const logout = () => {
    setState(emptyState);
    localStorage.removeItem("auth");
  };

  return !!state ? (
    <AuthContext.Provider
      value={{
        state,
        clientId,
        onLogin: login,
        onLogout: logout
      }}>
      <GoogleOAuthProvider clientId={clientId}>{children}</GoogleOAuthProvider>
    </AuthContext.Provider>
  ) : null;
};

AuthProvider.propTypes = {
  clientId: PropTypes.string.isRequired
};

export default AuthContext;
