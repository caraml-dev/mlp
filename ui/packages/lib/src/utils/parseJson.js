export const parseJson = response => {
  return response.json().catch(error => {
    // for responses without body
    if (error.name === "SyntaxError") {
      return {};
    }
    throw new Error(error.message);
  });
};
