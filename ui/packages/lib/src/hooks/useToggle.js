import { useCallback, useState } from "react";

export const useToggle = (initialState = false) => {
  const [isShowing, setIsShowing] = useState(initialState);

  const toggle = useCallback(() => setIsShowing(isShowing => !isShowing), []);

  return [isShowing, toggle];
};
