import { get, normalizePath, set } from "@gojek/mlp-ui/src/utils";

export const extractErrors = (validationError) => {
  let errors = {};
  if (validationError.inner) {
    for (let err of validationError.inner) {
      const path = normalizePath(err.path);
      const fieldsErrors = get(errors, path) || [];
      set(errors, path, [...fieldsErrors, err.message]);
    }
  }
  return errors;
};
