const ModuleFederationPlugin = require("webpack").container
  .ModuleFederationPlugin;
const paths = require("react-scripts/config/paths");
const deps = require("./package.json").dependencies;

module.exports = {
  plugins: [
    {
      plugin: {
        // Background:
        // https://github.com/facebook/create-react-app/issues/9510#issuecomment-902536147
        overrideWebpackConfig: ({ webpackConfig }) => {
          const htmlWebpackPlugin = webpackConfig.plugins.find(
            plugin => plugin.constructor.name === "HtmlWebpackPlugin"
          );
          htmlWebpackPlugin.userOptions = {
            ...htmlWebpackPlugin.userOptions,
            // Set custom publicPath
            publicPath: paths.publicUrlOrPath,
            // Exclude the exposed app for hot reloading to work
            excludeChunks: ["mlp"]
          };

          return webpackConfig;
        }
      }
    }
  ],
  webpack: {
    output: {
      publicPath: "auto"
    },
    plugins: {
      add: [
        new ModuleFederationPlugin({
          name: "mlp",
          filename: "remoteEntry.js",
          exposes: {
            ".": "./src/AppRoutes"
          },
          shared: {
            "@emotion/react": {
              singleton: true,
              requiredVersion: deps["@emotion/react"]
            },
            "@gojek/mlp-ui": {
              singleton: true,
              requiredVersion: deps["@gojek/mlp-ui"]
            },
            react: {
              shareScope: "default",
              singleton: true,
              requiredVersion: deps["react"]
            },
            "react-dom": {
              singleton: true,
              requiredVersion: deps["react-dom"]
            },
            "react-router-dom": {
              singleton: true,
              requiredVersion: deps["react-router-dom"]
            }
          }
        })
      ]
    }
  }
};
