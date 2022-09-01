import React, { useContext } from "react";
import { PagePlaceholder, ProjectsContext } from "@gojek/mlp-ui";
import { Navigate } from "react-router-dom";

const Home = () => {
  const { currentProject } = useContext(ProjectsContext);

  return !currentProject ? (
    <PagePlaceholder />
  ) : (
    <Navigate to={`/projects/${currentProject.id}`} replace={true} />
  );
};

export default Home;
