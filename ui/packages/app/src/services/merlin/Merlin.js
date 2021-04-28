import { useContext, useEffect, useState } from "react";
import { ApplicationsContext } from "@gojek/mlp-ui";

export const MerlinApp = () => {
  const { apps } = useContext(ApplicationsContext);
  const [merlinApp, setMerlinApp] = useState({});

  useEffect(() => {
    if (apps) {
      const merlin = apps.find(app => app.name === "Merlin");
      setMerlinApp(merlin);
    }
  }, [apps]);

  return merlinApp;
};
