import React from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiLoadingContent,
  EuiPageTemplate
} from "@elastic/eui";

export const PagePlaceholder = () => (
  <EuiPageTemplate restrictWidth="90%">
    <EuiPageTemplate.Section grow={false}>
      <EuiLoadingContent lines={2} />
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section grow={false}>
      <EuiLoadingContent lines={5} />
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup direction="row">
        <EuiFlexItem>
          <EuiLoadingContent lines={10} />
        </EuiFlexItem>

        <EuiFlexItem>
          <EuiLoadingContent lines={10} />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </EuiPageTemplate>
);
