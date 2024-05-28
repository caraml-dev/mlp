import babel from '@rollup/plugin-babel';
import resolve from "@rollup/plugin-node-resolve";
import terser from "@rollup/plugin-terser";
import path from "path";
import peerDepsExternal from 'rollup-plugin-peer-deps-external';
import sass from "rollup-plugin-sass";
import pkg from "./package.json";

module.exports = {
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
    resolve(),
    babel({
      babelHelpers: "bundled",
      exclude: "node_modules/**",
      presets: [['@babel/preset-react', { "runtime": "automatic" }]],
      extensions: ['.js', '.jsx']
    }),
    sass({
      output: true,
      options: {
        quietDeps: true,
        includePaths: [path.resolve(__dirname, "../../node_modules")]
      }
    }),
    terser(),
  ]
}
