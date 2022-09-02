/* global __webpack_init_sharing__ */
/* global __webpack_share_scopes__ */
import React, { useMemo } from "react";
import { useDynamicScript } from "./useDynamicScript";

// See: https://github.com/module-federation/module-federation-examples/blob/7fc92f8f7678a7af76f925b9d5d1d03eb472109e/advanced-api/dynamic-remotes/app1/src/App.js#L64-L84
function loadComponent(scope, module) {
  return async () => {
    // Initializes the share scope. This fills it with known provided modules from this build and all remotes
    await __webpack_init_sharing__("default");
    const container = window[scope]; // or get the container somewhere else
    // Initialize the container, it may provide shared modules
    await container.init(__webpack_share_scopes__.default);
    const factory = await window[scope].get(module);
    const Module = factory();
    return Module;
  };
}

const componentCache = new Map();
export const useFederatedComponent = (remoteUrl, scope, module) => {
  const { ready, errorLoading } = useDynamicScript(remoteUrl);

  const Component = useMemo(() => {
    if (ready) {
      const key = `${remoteUrl}-${scope}-${module}`;

      if (!componentCache.has(key)) {
        const component = React.lazy(loadComponent(scope, module));
        componentCache.set(key, component);
      }
      return componentCache.get(key);
    }
    return null;
  }, [module, ready, remoteUrl, scope]);

  return { Component, errorLoading };
};
