const pkg = require("./package.json");
const WriteVersionPlugin = require("./write-version-plugin");

module.exports = {
  mode: "production",
  entry: {
    fetch: "./fetch/js/index.js",
    url: "./url/js/index.js",
  },
  output: {
    path: __dirname,
    filename: "[name]/bundle.js",
  },
  target: "es2020",
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
