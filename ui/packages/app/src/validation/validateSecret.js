export const validateSecretName = name => {
  const expression = /^(?!\s*$).+/;
  return expression.test(String(name).toLowerCase());
};

export const validateSecretKey = key => {
  const expression = /^(?!\s*$).+/;
  return expression.test(String(key).toLowerCase());
};
