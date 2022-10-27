import React, { useContext } from "react";
import PropTypes from "prop-types";
import AuthContext from "./context";
import { Navigate, useLocation } from "react-router-dom";

export const PrivateRoute = ({ redirectTo = "/login", children }) => {
  const location = useLocation();

  const {
    state: { isAuthenticated }
  } = useContext(AuthContext);

  return isAuthenticated ? (
    children
  ) : (
    <Navigate
      to={redirectTo}
      state={{ referer: location.pathname + location.search }}
      replace={true}
    />
  );
};

PrivateRoute.propTypes = {
  redirectTo: PropTypes.string
};
