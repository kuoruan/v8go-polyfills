const pkg = require("./package.json");
const WriteVersionPlugin = require("./write-version-plugin");

/**
 * @type {import("webpack").Configuration}
 */
module.exports = {
  mode: "production",
  target: "browserslist",
  entry: {
    url: "./url/js/index.js",
    streams: "./streams/js/index.js",
    text: "./text/js/index.js",
    file: "./file/js/index.js",
  },
  output: {
    path: __dirname,
    filename: "[name]/bundle.js",
    library: {
      // emit name
      // https://webpack.js.org/configuration/output/#expose-via-object-assignment
      type: "window",
    },
  },
  plugins: [
    new WriteVersionPlugin({
      path: "internal/version.txt",
      content: pkg.version,
    }),
  ],
  stats: {
    all: false,
    assets: true,
    assetsSort: "size",
    entrypoints: true,
    errors: true,
    timings: true,
    warnings: true,
  },
};
