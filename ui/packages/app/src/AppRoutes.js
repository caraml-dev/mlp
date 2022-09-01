import React from "react";
import { Home, Project } from "./pages";
import { ProjectCreation } from "./project_setting/ProjectCreation";
import ProjectSetting from "./project_setting/ProjectSetting";
import { Navigate, Route, Routes } from "react-router-dom";

export const AppRoutes = () => (
  <Routes>
    {/* LANDING */}
    <Route index element={<Home />} />

    <Route path="projects">
      {/* PROJECT LANDING PAGE */}
      <Route path=":projectId" element={<Project />} />
      {/* PROJECT SETTING */}
      <Route path=":projectId/settings/*" element={<ProjectSetting />} />
      {/* New Project */}
      <Route path="create" element={<ProjectCreation />} />
    </Route>

    {/* DEFAULT */}
    <Route path="*" element={<Navigate to="/pages/404" />} />
  </Routes>
);

export default AppRoutes;
