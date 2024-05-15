// Test whether the value follow RFC1123 format
const DNS1123LabelMaxLength = 63;
export const isDNS1123Label = value => {
  const expression = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  if (value === undefined || value.length > DNS1123LabelMaxLength) {
    return false;
  }

  return expression.test(value);
};

export const validateSecretName = name => {
  const expression = /^(?!\s*$).+/;
  return expression.test(String(name).toLowerCase());
};

export const validateSecretData = key => {
  const expression = /^(?!\s*$).+/;
  return expression.test(String(key).toLowerCase());
};

export const validateEmail = email => {
  const expression = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/;
  return expression.test(email);
};

// See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set for full details
export const isValidK8sLabelValue = value => {
  const expression = /^[a-z0-9]([_.\-a-z0-9]*[a-z0-9])?$/;
  if (value === undefined || value.length > DNS1123LabelMaxLength) {
    return false;
  }
  return expression.test(value);
};
