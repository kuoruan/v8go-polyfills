const pkg = require("./package.json");
const WriteVersionPlugin = require("./write-version-plugin");

module.exports = {
  mode: "production",
  entry: {
    url: "./url/js/index.js",
  },
  output: {
    chunkFormat: "commonjs",
    path: __dirname,
    filename: "[name]/bundle.js",
  },
  target: "es2021",
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
