import React, { useContext } from "react";
import { PagePlaceholder, ProjectsContext } from "@caraml-dev/ui-lib";
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
