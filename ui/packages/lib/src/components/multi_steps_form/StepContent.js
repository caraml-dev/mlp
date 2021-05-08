import React from "react";
import { EuiFlexGroup, EuiFlexItem } from "@elastic/eui";

export const StepContent = ({ children, width = "75%" }) => (
  <EuiFlexGroup direction="row" justifyContent="center">
    <EuiFlexItem grow={false} style={{ width }}>
      {children}
    </EuiFlexItem>
  </EuiFlexGroup>
);
