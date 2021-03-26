const fs = require("fs");
const path = require("path");

const pkg = require("./package.json");

class WriteVersionPlugin {
  constructor ({path, content}) {
    this.path = path;
    this.content = content;
  }

  apply(compiler) {
    compiler.hooks.done.tap('Write Version Plugin', () => {
      if (this.path && this.content) {
        fs.writeFileSync(path.resolve(this.path), this.content, 'utf-8')
      }
    });
  }
}

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
  plugins: [
    new WriteVersionPlugin({
      path: "internal/version.txt",
      content: pkg.version
    })
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
}
