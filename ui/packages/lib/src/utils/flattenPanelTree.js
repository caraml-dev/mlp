export const flattenPanelTree = (tree, array = []) => {
  array.push(tree);
  if (tree.items) {
    tree.items.forEach(item => {
      if (item.panel) {
        flattenPanelTree(item.panel, array);
        item.panel = item.panel.id;
      }
    });
  }
  return array;
};
