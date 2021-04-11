import React from "react";
import {
  EuiDescriptionList,
  EuiHorizontalRule,
  EuiSplitPanel,
  EuiTitle
} from "@elastic/eui";

export const Panel = ({ title, items, type, actions }) => {
  return (
    <EuiSplitPanel.Outer grow>
      <EuiSplitPanel.Inner grow>
        <EuiTitle size="xs">
          <span>{title}</span>
        </EuiTitle>
        <EuiHorizontalRule size="full" margin="s" />

        <EuiDescriptionList
          compressed
          listItems={items}
          type={type ? type : "responsiveColumn"}
        />
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
