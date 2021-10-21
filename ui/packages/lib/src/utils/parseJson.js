const jsonBig = require(`json-bigint`);

export const parseJson = response => {
  // using the json-bigint library to parse the data instead of the Response.json() Web API,
  // to parse BigInt without losing precision.
  return response.text().then(
    text => jsonBig.parse(text)
  ).catch(error => {
    // for responses without body
    if (error.name === "SyntaxError") {
      return {};
    }
    throw new Error(error.message);
  });
};
