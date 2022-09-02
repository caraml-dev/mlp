import React, { useState } from "react";
import { EuiHeaderBreadcrumbs } from "@elastic/eui";
import "./Breadcrumbs.scss";

let replaceBreadcrumbsHandler = () => {};

export const replaceBreadcrumbs = breadcrumbs => {
  replaceBreadcrumbsHandler(breadcrumbs);
};

export const Breadcrumbs = props => {
  const [breadcrumbs, setBreadcrumbs] = useState([]);

  replaceBreadcrumbsHandler = setBreadcrumbs;

  return (
    <EuiHeaderBreadcrumbs
      className="euiBreadcrumbs--header"
      breadcrumbs={breadcrumbs}
      {...props}
    />
  );
};
