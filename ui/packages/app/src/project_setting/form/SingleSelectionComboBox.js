import React, { useState } from "react";
import { EuiComboBox } from "@elastic/eui";

export const SingleSelectionComboBox = ({
  options,
  onChange,
  onValidChange
}) => {
  const [allOptions] = useState(
    options.sort((a, b) => (a.label < b.label ? -1 : a.label > b.label ? 1 : 0))
  );
  const [selectedOption, setSelectedOption] = useState([]);
  const [isValid, setIsValid] = useState(true);

  const onCreate = (searchValue, flattenedOptions = []) => {
    const normalizedSearchValue = searchValue.trim().toLowerCase();
    if (!normalizedSearchValue) {
      return;
    }
    const newValue = {
      label: searchValue
    };

    if (
      flattenedOptions.findIndex(
        option => option.label.trim().toLowerCase() === normalizedSearchValue
      ) === -1
    ) {
      allOptions.push(newValue);
    }
    let isValid = newValue !== "";
    setIsValid(isValid);
    onValidChange(isValid);
    onSelectionChange([newValue]);
  };

  const onSelectionChange = e => {
    let isValid = e.length >= 1;
    setSelectedOption(e);
    setIsValid(isValid);
    onValidChange(isValid);
    if (isValid) {
      onChange(e[0]);
    } else {
      onChange("");
    }
  };

  return (
    <EuiComboBox
      options={allOptions}
      singleSelection={{ asPlainText: true }}
      isClearable={true}
      onChange={onSelectionChange}
      onCreateOption={onCreate}
      selectedOptions={selectedOption}
      isInvalid={!isValid}
    />
  );
};
