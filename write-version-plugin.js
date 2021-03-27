const fs = require("fs");
const path = require("path");

class WriteVersionPlugin {
  constructor({ path, content = "" }) {
    this.path = path;
    this.content = content;
  }

  apply(compiler) {
    compiler.hooks.done.tap("Write Version Plugin", () => {
      if (this.path) {
        fs.writeFileSync(path.resolve(this.path), this.content, "utf-8");
      }
    });
  }
}

module.exports = WriteVersionPlugin;
