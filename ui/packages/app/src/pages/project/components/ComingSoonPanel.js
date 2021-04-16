import React from "react";
import {
  EuiCard,
  EuiIcon,
} from "@elastic/eui";

export const ComingSoonPanel = ({ title, iconType }) => {
  return (
    <EuiCard
      icon={<EuiIcon size="xl" type={iconType} />}
      title={title}
      description="Coming soon."
      layout="horizontal"
      // betaBadgeLabel={badges[index]}
      // betaBadgeTooltipContent={
      //   badges[index]
      //     ? 'This module is not GA. Please help us by reporting any bugs.'
      //     : undefined
      // }
      onClick={() => {}}
    />
  );
};
