// Test whether the subdomain follow RFC1123 format
export const validateSubdomain = subdomain => {
  const expression = /^[a-zA-Z][a-zA-Z0-9-]+[a-zA-Z0-9]$/;
  return expression.test(String(subdomain.toLowerCase()));
};
