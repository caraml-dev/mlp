import React, { useContext } from "react";
import { Home, Project } from "./pages";
import { ProjectCreation } from "./project_setting/ProjectCreation";
import ProjectSetting from "./project_setting/ProjectSetting";
import { Navigate, Route, Routes } from "react-router-dom";
import { PlaceholderPage } from "./placeholder_page/PlaceholderPage";
import { ApplicationsContext } from "@caraml-dev/ui-lib";

export const AppRoutes = () => {
  const { apps, isLoaded } = useContext(ApplicationsContext);

  // If the apps have not been loaded yet, we do not render any of the app related routes - the additional placeholder
  // apps need to be retrieved from the MLP API v2/applications endpoint first before we generate each route for them.
  return isLoaded && (
    <Routes>
      {/* LANDING */}
      <Route index element={<Home />} />

      {apps?.map(app => !!app.placeholder_page_config &&
        <Route key={app.name} path={app.homepage} element={<PlaceholderPage app={app} />} />)
      }

      <Route path="projects">
        {/* PROJECT LANDING PAGE */}
        <Route path=":projectId" element={<Project />} />
        {/* PROJECT SETTING */}
        <Route path=":projectId/settings/*" element={<ProjectSetting />} />
        {/* New Project */}
        <Route path="create" element={<ProjectCreation />} />
      </Route>

      {/* DEFAULT */}
      <Route path="*" element={<Navigate to="/pages/404" replace={true} />} />
    </Routes>
  )
};

export default AppRoutes;
