import babel from "rollup-plugin-babel";
import resolve from "rollup-plugin-node-resolve";
import sass from "rollup-plugin-sass";
import peerDepsExternal from "rollup-plugin-peer-deps-external";
import ignoreImport from "rollup-plugin-ignore-import";
import nodeSass from "node-sass";
import { eslint } from "rollup-plugin-eslint";

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
        includePaths: ["../../node_modules/"],
        importer(path) {
          return { file: path[0] !== "~" ? path : path.slice(1) };
        }
      },
      runtime: nodeSass
    }),
    resolve(),
    babel({
      exclude: "node_modules/**"
    })
  ]
};
