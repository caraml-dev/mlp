import { navigate } from "@reach/router";
import { slugify } from "@gojek/mlp-ui/src/utils";

export const pages = {};

export const navigation = [
  {
    name: "Settings",
    items: [
      {
        name: "Connected Accounts",
        isSelected: true
      }
    ]
  }
].map(({ name, items, ...rest }) => ({
  name,
  id: slugify(name),
  items: items.map(({ name: itemName, ...rest }) => {
    const id = `/${slugify(name)}/${slugify(itemName)}`;

    return {
      id: id,
      name: itemName,
      onClick: () => navigate(id),
      ...rest
    };
  }),
  ...rest
}));
