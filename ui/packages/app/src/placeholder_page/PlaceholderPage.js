import React from "react";
import { EuiPageTemplate } from "@elastic/eui";

export const PlaceholderPage = ({ app }) => {
  const iframe = `<iframe src=${app.placeholder_page_config.url} style="height: 80vh; width: 100%;"></iframe>`

  function Iframe(props) {
    return (<div dangerouslySetInnerHTML={ {__html:  props.iframe?props.iframe:""}} />);
  }
  return (
    <EuiPageTemplate restrictWidth="90%" panelled={false}>
      <EuiPageTemplate.Header
        bottomBorder={false}
        iconType={app.config.icon}
        pageTitle={app.name}
      />

      <EuiPageTemplate.Section paddingSize="none">
        <Iframe iframe={iframe} />
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};
