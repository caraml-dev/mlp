import React, { useEffect, useState } from "react";
import PropTypes from "prop-types";

const emptyState = {
  isAuthenticated: false,
  profileObj: null,
  accessToken: "",
  expiresAt: 0
};

const AuthContext = React.createContext({
  state: emptyState
});

export const AuthProvider = ({ clientId, children }) => {
  const [state, setState] = useState(undefined);

  useEffect(() => {
    // TODO: not safe!
    const localStateStr = localStorage.getItem("auth");

    if (localStateStr) {
      const authObject = JSON.parse(localStateStr);
      setState(authObject.expiresAt > Date.now() ? authObject : emptyState);
    } else {
      setState(emptyState);
    }
  }, []);

  const setAuth = data => {
    setState(data);
    localStorage.setItem("auth", JSON.stringify(data));
  };

  const login = ({ profileObj, accessToken, tokenObj }) => {
    setAuth({
      ...state,
      isAuthenticated: true,
      profileObj: profileObj,
      accessToken: accessToken,
      expiresAt: tokenObj.expires_at
    });
  };

  const logout = () => {
    setAuth(emptyState);
  };

  return !!state ? (
    <AuthContext.Provider
      value={{
        state,
        clientId,
        onLogin: login,
        onLogout: logout
      }}>
      {children}
    </AuthContext.Provider>
  ) : null;
};

AuthProvider.propTypes = {
  clientId: PropTypes.string.isRequired
};

export default AuthContext;
