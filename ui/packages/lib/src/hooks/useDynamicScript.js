import React from "react";

const urlCache = new Set();

// See: https://github.com/module-federation/module-federation-examples/blob/7fc92f8f7678a7af76f925b9d5d1d03eb472109e/advanced-api/dynamic-remotes/app1/src/App.js#L17-L61
export const useDynamicScript = url => {
  const [ready, setReady] = React.useState(false);
  const [errorLoading, setErrorLoading] = React.useState(false);

  React.useEffect(() => {
    if (!url) return;

    if (urlCache.has(url)) {
      setReady(true);
      setErrorLoading(false);
      return;
    }

    setReady(false);
    setErrorLoading(false);

    const element = document.createElement("script");

    element.src = url;
    element.type = "text/javascript";
    element.async = true;

    element.onload = () => {
      urlCache.add(url);
      setReady(true);
    };

    element.onerror = () => {
      setReady(false);
      setErrorLoading(true);
    };

    document.head.appendChild(element);

    return () => {
      urlCache.delete(url);
      document.head.removeChild(element);
    };
  }, [url]);

  return {
    errorLoading,
    ready
  };
};
