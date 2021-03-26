module.exports = {
  mode: "production",
  entry: {
    "fetch": "./fetch/js/index.js",
    "url": "./url/js/index.js",
  },
  output: {
    path: __dirname,
    filename: "[name]/bundle.js",
  },
  target: "es2020",
  stats: {
    all: false,
    assets: true,
    assetsSort: "size",
    entrypoints: true,
    errors: true,
    timings: true,
    warnings: true,
  },
}
