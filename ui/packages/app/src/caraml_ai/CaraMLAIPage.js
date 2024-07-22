import React from "react";
import { EuiPageTemplate } from "@elastic/eui";
import config from "./../config";

export const CaraMLAIPage = () => {
  const iframe = `<iframe src=${config.CARAML_AI_STREAMLIT_HOMEPAGE} style="height: 80vh; width: 100%;"></iframe>`

  function Iframe(props) {
    return (<div dangerouslySetInnerHTML={ {__html:  props.iframe?props.iframe:""}} />);
  }
  return (
    <EuiPageTemplate restrictWidth="90%" panelled={false}>
      <EuiPageTemplate.Header
        bottomBorder={false}
        iconType="timelionApp"
        pageTitle="CaraML AI"
      />

      <EuiPageTemplate.Section paddingSize="none">
        <Iframe iframe={iframe} />
      </EuiPageTemplate.Section>
    </EuiPageTemplate>
  );
};
