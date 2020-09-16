export const validateEmail = email => {
  const expression = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/;
  return expression.test(String(email).toLowerCase());
};
