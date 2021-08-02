import React from "react";
import { EuiBadge, EuiLink, EuiBadgeGroup } from "@elastic/eui";
import EllipsisText from "react-ellipsis-text";
import useCollapse from "react-collapsed";

/**
 * CollapsibleLabelsPanel is a collapsible panel that groups labels for a model version.
 * By default, only 2 labels are shown, and a "Show All" button needs to be clicked
 * to display all the labels.
 *
 * @param labels is a dictionary of label key and value.
 * @param labelOnClick is a callback function that accepts {queryText: "..."} object, that will be triggered when a label is clicked.
 * @param minLabelsCount is the no of labels to show when the panel is collapsed.
 * @param maxLabelLength is the max no of characters of label key / value to show, characters longer than maxLabelLength will be shown as ellipsis.
 * @returns {JSX.Element}
 * @constructor
 */
export const CollapsibleLabelsPanel = ({
  labels,
  labelOnClick,
  minLabelsCount = 2,
  maxLabelLength = 9
}) => {
  const { getToggleProps, isExpanded } = useCollapse();

  return (
    <EuiBadgeGroup>
      {labels &&
        Object.entries(labels).map(
          ([key, val], index) =>
            (isExpanded || index < minLabelsCount) && (
              <EuiBadge
                key={key}
                onClick={() => {
                  const queryText = `labels: ${key} in (${val})`;
                  labelOnClick({ queryText });
                }}
                onClickAriaLabel="search by label">
                <EllipsisText text={key} length={maxLabelLength} />:
                <EllipsisText text={val} length={maxLabelLength} />
              </EuiBadge>
            )
        )}
      {// Toggle collapse button
      !isExpanded && labels && Object.keys(labels).length > minLabelsCount && (
        <EuiLink {...getToggleProps()}>
          {isExpanded
            ? ""
            : `Show All [${Object.keys(labels).length - minLabelsCount}]`}
        </EuiLink>
      )}
    </EuiBadgeGroup>
  );
};
