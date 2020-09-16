import React from "react";
import { EuiText, EuiToolTip } from "@elastic/eui";

const moment = require("moment");

export const DateFromNow = ({ date, size }) => {
  const momentDate = moment(date, "YYYY-MM-DDTHH:mm.SSZ");

  return (
    <EuiToolTip position="top" content={momentDate.toLocaleString()}>
      <EuiText size={size}>{momentDate.fromNow()}</EuiText>
    </EuiToolTip>
  );
};
