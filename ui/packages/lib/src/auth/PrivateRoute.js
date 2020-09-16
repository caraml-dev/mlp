import React, { useContext, useMemo } from "react";
import PropTypes from "prop-types";
import AuthContext from "./context";
import { StatefulRedirect } from "../components";

export const PrivateRoute = ({ redirectTo = "/login", render, ...props }) => {
  const {
    state: { isAuthenticated }
  } = useContext(AuthContext);

  const referer = useMemo(
    () => props.location.pathname + props.location.search,
    [props.location.pathname, props.location.search]
  );

  return isAuthenticated ? (
    render(props)
  ) : (
    <StatefulRedirect to={redirectTo} state={{ referer: referer }} />
  );
};

PrivateRoute.propTypes = {
  path: PropTypes.string.isRequired,
  redirectTo: PropTypes.string,
  render: PropTypes.func.isRequired
};
