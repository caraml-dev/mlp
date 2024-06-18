import React, { useState } from "react";

import {
  EuiButtonEmpty,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButton,
  EuiFormRow
} from "@elastic/eui";
import { isValidK8sLabelKeyValue } from "../../validation/validation";
import config from "../../config";

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
          isValueValid: true,
          existsBefore: true
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
        isValueValid: false,
        existsBefore: false
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
        isKeyValid: isValidK8sLabelKeyValue(newKey)
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
        isValueValid: isValidK8sLabelKeyValue(newValue)
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
      delete element.existsBefore;
      return element;
    });

    setLabels(newLabels);
  };

  return (
    <EuiFormRow isInvalid={!isValidLabels} error={labelError}>
      <EuiFlexGroup direction="column" gutterSize="m">
        {items.map((element, idx) => {
          const isFieldDisabled =
            isDisabled ||
            (element.existsBefore && config.LABELS_BLACKLIST[element.key]);
          return (
            <EuiFlexItem key={idx}>
              <EuiFlexGroup gutterSize="s">
                <EuiFlexItem grow={1}>
                  <EuiFieldText
                    placeholder="key"
                    value={element.key}
                    onChange={onKeyChange(idx)}
                    isInvalid={!element.isKeyValid}
                    disabled={isFieldDisabled}
                  />
                </EuiFlexItem>
                <EuiFlexItem grow={1}>
                  <EuiFieldText
                    placeholder="value"
                    value={element.value}
                    onChange={onValueChange(idx)}
                    isInvalid={!element.isValueValid}
                    disabled={isFieldDisabled}
                  />
                </EuiFlexItem>
                <EuiFlexItem grow={false}>
                  <EuiButtonEmpty
                    iconType="trash"
                    onClick={removeElement(idx)}
                    color="danger"
                    disabled={isFieldDisabled}
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
              isDisabled ||
              (items.length === 0
                ? false
                : items.reduce((addButtonDisabled, currentValue) => {
                    return (
                      addButtonDisabled ||
                      !currentValue.isKeyValid ||
                      !currentValue.isValueValid
                    );
                  }, false))
            }
          />
        </EuiFlexItem>
      </EuiFlexGroup>
    </EuiFormRow>
  );
};
