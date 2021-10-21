const jsonBig = require(`json-bigint`);

const parseJsonHandleBigInt = (response, parseBigInt) => {
  // If parseBigInt is set, use the json-bigint library to parse the data instead of
  // the Response.json() Web API which would lose precision.
  return parseBigInt
    ? response.text().then(text => jsonBig.parse(text))
    : response.json();
}

export const parseJson = (response, parseBigInt) => {
  return parseJsonHandleBigInt(response, parseBigInt).catch(error => {
    // for responses without body
    if (error.name === "SyntaxError") {
      return {};
    }
    throw new Error(error.message);
  });
};
