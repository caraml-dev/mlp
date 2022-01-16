import { useState } from "react";

// usePromise returns a lazy promise (a callback that returns a Promise),
// along with a callback function that will either resolve or reject the promise,
// based on the execution of the promise.
export const usePromise = (callback) => {
  // Init state with empty function call
  const [resolveOrReject, setResolveOrReject] = useState(() => () => { });

  const promise = () =>
    new Promise((resolve, reject) => {
      try {
        callback();
        setResolveOrReject(() => resolve);
      } catch (error) {
        console.error(error);
        setResolveOrReject(() => reject);
      }
    });

  return [promise, resolveOrReject];
};
