import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiPageTemplate,
  EuiSkeletonText
} from "@elastic/eui";
import React from "react";

export const PagePlaceholder = () => (
  <EuiPageTemplate restrictWidth="90%">
    <EuiPageTemplate.Section grow={false}>
      <EuiSkeletonText lines={2} />
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section grow={false}>
      <EuiSkeletonText lines={5} />
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiFlexGroup direction="row">
        <EuiFlexItem>
          <EuiSkeletonText lines={10} />
        </EuiFlexItem>

        <EuiFlexItem>
          <EuiSkeletonText lines={10} />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiPageTemplate.Section>
  </EuiPageTemplate>
);
