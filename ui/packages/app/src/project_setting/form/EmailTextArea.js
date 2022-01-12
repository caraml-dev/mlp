import React, { useState } from "react";
import { EuiTextArea } from "@elastic/eui";
import { validateEmail } from "../../validation/validation";

export const EmailTextArea = ({ onChange, onValidChange }) => {
  const onValueChange = e => {
    let emails = e.target.value.replace(/\s/g, "").split(",");
    let valid = true;

    for (var i = 0; i < emails.length; i++) {
      let email = emails[i];
      if (email === "" || !validateEmail(email)) {
        valid = false;
        break;
      }
    }
    setIsValid(valid);
    setValue(e.target.value);
    onValidChange(valid);
    onChange(emails);
  };

  const [isValid, setIsValid] = useState(true);
  const [value, setValue] = useState();

  return (
    <EuiTextArea
      resize="none"
      placeholder="foo@go-jek.com, bar@go-jek.com"
      value={value}
      isInvalid={!isValid}
      onChange={onValueChange}
    />
  );
};

export default EmailTextArea;
