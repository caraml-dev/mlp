import React, { useCallback, useContext, useEffect, useState } from "react";
import { FormContext } from "../context";
import FormValidationContext from "./context";
import { extractErrors } from "./errors";
import debounce from "lodash/debounce";
import zip from "lodash/zip";

// Debounced method would only be run at most once every configured duration.
// Setting a larger value improves the performance, however increases the
// delay between form value changes and feedback to the user.
const DEBOUNCE_INTERVAL_MS = 300;

const debouncedValidate = debounce(
  (schemas, contexts, formData, setErrors, setIsValidated) => {
    Promise.all(
      zip(schemas, contexts).map(([schema, ctx]) => {
        return !!schema
          ? new Promise((resolve, reject) => {
            schema
              .validate(formData, {
                abortEarly: false,
                context: ctx,
              })
              .then(
                () => resolve({}),
                (err) => resolve(extractErrors(err))
              );
          })
          : Promise.resolve({});
      })
    )
      .then(setErrors)
      .then(() => setIsValidated(true));
  },
  DEBOUNCE_INTERVAL_MS
);


export const MultiSectionFormValidationContextProvider = ({
  schemas,
  contexts,
  /* The result of executing onSubmit will be resolved using Promise.resolve.
     If a thenable, the resetting of the isSubmitting state will be chained to it.
  */
  onSubmit,
  children
}) => {
  const { data: formData } = useContext(FormContext);

  // identifies if user tried to submit this form
  const [isTouched, setIsTouched] = useState(false);

  // identifies if the form was validated
  const [isValidated, setIsValidated] = useState(false);

  // identifies if the form is in submission state
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState([]);

  const isValid = (errors) =>
    errors.reduce(
      (isValid, errors) => isValid && !Object.keys(errors).length,
      true
    );

  const onStartSubmitting = (event) => {
    event && event.preventDefault();
    setIsTouched(true);
    setIsSubmitting(true);
  };

  const onFinishSubmitting = useCallback(() => {
    setIsTouched(false);
    setIsValidated(false);
    // Execute the onSubmit callback and at last, reset the submitting status.
    // If the onSubmit was defined as a lazy Promise, we must chain the reset action to the Promise.
    // This will ensure that downstream actions (such as re-enabling the Submit button) are paused
    // until we have a success/failure response from the onSubmit call.
    Promise
      .resolve(onSubmit())
      .finally(() => { setIsSubmitting(false); })
  }, [onSubmit]);

  useEffect(() => {
    if (isTouched) {
      if (schemas) {
        debouncedValidate(
          schemas,
          contexts,
          formData,
          setErrors,
          setIsValidated
        );
      } else {
        setIsValidated(true);
      }
    }
  }, [isTouched, schemas, contexts, formData]);

  useEffect(() => {
    if (isSubmitting && isValidated) {
      isValid(errors) ? onFinishSubmitting() : setIsSubmitting(false);
    }
  }, [isSubmitting, isValidated, errors, onFinishSubmitting]);

  return (
    <FormValidationContext.Provider
      value={{ onSubmit: onStartSubmitting, isSubmitting, errors }}>
      {children}
    </FormValidationContext.Provider>
  );
};
