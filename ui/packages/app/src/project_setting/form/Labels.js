import React, { useState } from "react";

import {
  EuiButtonEmpty,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButton,
  EuiFormRow
} from "@elastic/eui";
import { isDNS1123Label } from "../../validation/validation";

export const Labels = ({
  labels,
  setLabels,
  setIsValidLabels,
  isValidLabels,
  isDisabled = false
}) => {
  const [items, setItems] = useState(
    (e => {
      if (e) {
        return e.map((label, idx) => ({
          ...label,
          idx,
          isKeyValid: true,
          isValueValid: true
        }));
      } else {
        return [];
      }
    })(labels)
  );

  const addItem = () => {
    const newItems = [
      ...items,
      {
        idx: items.length,
        isKeyValid: false,
        isValueValid: false
      }
    ];
    setItems(newItems);
    onChange(newItems);
  };

  const onKeyChange = idx => {
    return e => {
      const newItems = [...items];
      const newKey = e.target.value.trim();
      newItems[idx] = {
        ...newItems[idx],
        key: newKey,
        isKeyValid: isDNS1123Label(newKey)
      };
      setItems(newItems);
      onChange(newItems);
    };
  };

  const onValueChange = idx => {
    return e => {
      const newItems = [...items];
      const newValue = e.target.value.trim();
      newItems[idx] = {
        ...newItems[idx],
        value: newValue,
        isValueValid: isDNS1123Label(newValue)
      };
      setItems(newItems);
      onChange(newItems);
    };
  };

  const removeElement = idx => {
    return e => {
      const newItems = [...items];
      newItems.splice(idx, 1);
      setItems(newItems);
      onChange(newItems);
    };
  };

  const [labelError, setLabelError] = useState("");
  const onChange = labels => {
    const labelsValid =
      labels.length === 0
        ? true
        : labels.reduce((labelsValid, label) => {
            return labelsValid && label.isKeyValid && label.isValueValid;
          }, true);
    setIsValidLabels(labelsValid);
    if (!labelsValid) {
      setLabelError(
        "Invalid labels. Both key and value of a label must contain only lowercase alphanumeric and dash (-), and must start and end with an alphanumeric character"
      );
    }

    //deep copy
    let newLabels = JSON.parse(JSON.stringify(labels));
    newLabels = newLabels.map(element => {
      delete element.isKeyValid;
      delete element.isValueValid;
      delete element.idx;
      return element;
    });

    setLabels(newLabels);
  };

  return (
    <EuiFormRow isInvalid={!isValidLabels} error={labelError}>
      <EuiFlexGroup direction="column" gutterSize="m">
        {items.map((element, idx) => {
          return (
            <EuiFlexItem>
              <EuiFlexGroup gutterSize="s">
                <EuiFlexItem grow={1}>
                  <EuiFieldText
                    placeholder="key"
                    value={element.key}
                    onChange={onKeyChange(idx)}
                    isInvalid={!element.isKeyValid}
                    disabled={isDisabled}
                  />
                </EuiFlexItem>
                <EuiFlexItem grow={1}>
                  <EuiFieldText
                    placeholder="value"
                    value={element.value}
                    onChange={onValueChange(idx)}
                    isInvalid={!element.isValueValid}
                    disabled={isDisabled}
                  />
                </EuiFlexItem>
                <EuiFlexItem grow={false}>
                  <EuiButtonEmpty
                    iconType="trash"
                    onClick={removeElement(idx)}
                    color="danger"
                    disabled={isDisabled}
                  />
                </EuiFlexItem>
              </EuiFlexGroup>
            </EuiFlexItem>
          );
        })}
        <EuiFlexItem>
          <EuiButton
            iconType="plusInCircle"
            onClick={addItem}
            disabled={
              isDisabled
                ? true
                : items.length === 0
                ? false
                : items.reduce((addButtonDisabled, currentValue) => {
                    return (
                      addButtonDisabled ||
                      !currentValue.isKeyValid ||
                      !currentValue.isValueValid
                    );
                  }, false)
            }
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiFormRow>
  );
};
