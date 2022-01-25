import Blob from "./blob";

export default class File extends Blob {
  #lastModified = 0;
  #name = "";

  /**
   * @param {*[]} fileBits
   * @param {string} fileName
   * @param {{lastModified?: number, type?: string}} options
   */ // @ts-ignore
  constructor(fileBits, fileName, options = {}) {
    if (arguments.length < 2) {
      throw new TypeError(
        `Failed to construct 'File': 2 arguments required, but only ${arguments.length} present.`
      );
    }
    super(fileBits, options);

    // Simulate WebIDL type casting for NaN value in lastModified option.
    this.#lastModified = Number(options?.lastModified) || Date.now();
    this.#name = String(fileName);
  }

  get name() {
    return this.#name;
  }

  get lastModified() {
    return this.#lastModified;
  }

  /**
   * The path the URL of the File is relative to.
   * @type {string}
   */
  get webkitRelativePath() {
    return "";
  }

  get [Symbol.toStringTag]() {
    return "File";
  }
}
