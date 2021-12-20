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

    if (options === null) options = {};

    // Simulate WebIDL type casting for NaN value in lastModified option.
    const lastModified =
      options.lastModified === undefined
        ? Date.now()
        : Number(options.lastModified);
    if (!Number.isNaN(lastModified)) {
      this.#lastModified = lastModified;
    }

    this.#name = String(fileName);
  }

  get name() {
    return this.#name;
  }

  get lastModified() {
    return this.#lastModified;
  }

  get [Symbol.toStringTag]() {
    return "File";
  }
}
