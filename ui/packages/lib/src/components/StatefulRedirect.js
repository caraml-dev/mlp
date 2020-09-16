import { useEffect } from "react";
import { navigate } from "@reach/router";
import PropTypes from "prop-types";

// To replicate React Router's stateful redirect
export const StatefulRedirect = ({ to, state }) => {
  useEffect(() => {
    navigate(to, { state: state, replace: true });
  }, [to, state]);

  return null;
};

StatefulRedirect.propTypes = {
  to: PropTypes.string.isRequired,
  state: PropTypes.object
};
