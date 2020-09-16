import React, { Fragment } from "react";
import { EuiButton, EuiEmptyPrompt } from "@elastic/eui";

export const Empty = () => (
  <EuiEmptyPrompt
    iconType="gisApp"
    title={<h2>Error 404</h2>}
    body={
      <Fragment>
        <h3>You have no spice</h3>
        <p>
          Navigators use massive amounts of spice to gain a limited form of
          prescience. This allows them to safely navigate interstellar space,
          enabling trade and travel throughout the galaxy.
        </p>
        <p>You&rsquo;ll need spice to rule Arrakis, young Atreides.</p>
      </Fragment>
    }
    actions={
      <EuiButton color="primary" fill href="/">
        Embark Home
      </EuiButton>
    }
  />
);
