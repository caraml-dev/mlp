import babel from "@rollup/plugin-babel";
import resolve from "@rollup/plugin-node-resolve";
import eslint from "@rollup/plugin-eslint";
import sass from "rollup-plugin-sass";
import peerDepsExternal from "rollup-plugin-peer-deps-external";
import ignoreImport from "rollup-plugin-ignore-import";
import path from "path";

import pkg from "./package.json";

export default {
  input: "src/index.js",
  output: [
    {
      file: pkg.main,
      format: "cjs",
      sourcemap: true
    },
    {
      file: pkg.module,
      format: "es",
      sourcemap: true
    }
  ],
  plugins: [
    peerDepsExternal({
      includeDependencies: true
    }),
    ignoreImport({
      extensions: [".scss", ".css"]
    }),
    eslint(),
    sass({
      output: true,
      options: {
        quietDeps: true,
        includePaths: [path.resolve(__dirname, "../../node_modules")]
      }
    }),
    resolve(),
    babel({
      babelHelpers: "bundled",
      exclude: "node_modules/**"
    })
  ]
};
