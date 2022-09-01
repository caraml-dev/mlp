import React from "react";
import { useFederatedComponent } from "../../hooks";
import PropTypes from "prop-types";
import { ErrorBox } from "../ErrorBoundary";

export const FederatedComponent = ({
  scope,
  url,
  module = ".",
  fallback = "Loading...",
  error = e => <ErrorBox error={e} />,
  ...props
}) => {
  const { Component: FederatedComponent, errorLoading } = useFederatedComponent(
    url,
    scope,
    module
  );

  return (
    <React.Suspense fallback={fallback}>
      {!!errorLoading
        ? typeof error === "function"
          ? error(
              new Error(
                `Failed to load remote component "${module}" from "${scope}@${url}"`
              )
            )
          : error
        : !!FederatedComponent && <FederatedComponent {...props} />}
    </React.Suspense>
  );
};

FederatedComponent.propTypes = {
  scope: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  module: PropTypes.string,
  fallback: PropTypes.node,
  error: PropTypes.oneOfType([PropTypes.node, PropTypes.func])
};
