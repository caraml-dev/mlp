export const slugify = str => {
  const parts = str
    .toLowerCase()
    .replace(/[-]+/g, " ")
    .replace(/[^\w^\s]+/g, "")
    .replace(/ +/g, " ")
    .split(" ");
  return parts.join("-");
};
