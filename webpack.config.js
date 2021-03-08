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
    entrypoints: true,
    chunkGroups: true,
    timings: true,
    errors: true,
  },
}
