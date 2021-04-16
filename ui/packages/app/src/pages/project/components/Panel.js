import React from "react";
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiDescriptionList,
  EuiHorizontalRule,
  EuiIcon,
  EuiSpacer,
  EuiSplitPanel,
  EuiTitle
} from "@elastic/eui";

export const Panel = ({ title, items, type, iconType, actions }) => {
  return (
    <EuiSplitPanel.Outer grow>
      <EuiSplitPanel.Inner grow>
        <EuiFlexGroup gutterSize="m">
          <EuiFlexItem grow={false}>
            <EuiIcon size="xl" type={iconType} />
          </EuiFlexItem>
          <EuiFlexItem grow={9}>
            <EuiTitle size="s">
              <span>{title}</span>
            </EuiTitle>
            <EuiSpacer size="m" />
            <EuiDescriptionList
              compressed
              listItems={items}
              type={type ? type : "responsiveColumn"}
            />
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiSplitPanel.Inner>

      {actions && (
        <EuiSplitPanel.Inner grow={false}>
          <EuiHorizontalRule size="full" margin="s" />
          {actions}
        </EuiSplitPanel.Inner>
      )}
    </EuiSplitPanel.Outer>
  );
};
