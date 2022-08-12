import babel from "@rollup/plugin-babel";
import resolve from "@rollup/plugin-node-resolve";
import eslint from "@rollup/plugin-eslint";
import sass from "rollup-plugin-sass";
import peerDepsExternal from "rollup-plugin-peer-deps-external";
import ignoreImport from "rollup-plugin-ignore-import";

import pkg from "./package.json";

const PATH_NODE_MODULES = "../../node_modules";

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
        includePaths: [PATH_NODE_MODULES],
        importer(url, _) {
          return {
            file: url.replace(/^~/, `${PATH_NODE_MODULES}/`)
          };
        }
      }
    }),
    resolve(),
    babel({
      babelHelpers: "bundled",
      exclude: "node_modules/**"
    })
  ]
};
