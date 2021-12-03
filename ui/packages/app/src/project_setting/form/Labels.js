import React, { useState } from "react";

import {
  EuiButtonEmpty,
  EuiFieldText,
  EuiFlexGroup,
  EuiFlexItem,
  EuiButton
} from "@elastic/eui";
import { isDNS1123Label } from "../../validation/validation";

export const Labels = ({ onChange }) => {
  const [items, setItems] = useState([]);

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

  return (
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
                />
              </EuiFlexItem>
              <EuiFlexItem grow={1}>
                <EuiFieldText
                  placeholder="value"
                  value={element.value}
                  onChange={onValueChange(idx)}
                  isInvalid={!element.isValueValid}
                />
              </EuiFlexItem>
              <EuiFlexItem grow={false}>
                <EuiButtonEmpty
                  iconType="trash"
                  onClick={removeElement(idx)}
                  color="danger"
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
            items.length === 0
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
  );
};
