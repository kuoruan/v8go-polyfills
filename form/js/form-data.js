// https://github.com/jimmywarting/FormData/blob/master/FormData.js
// https://github.com/web-std/io/blob/main/packages/form-data/src/form-data.js

function normalizeLineFeeds(value) {
  return value.replace(/\r?\n|\r/g, "\r\n");
}

/**
 * @param {string|Blob|File} value
 * @param {string} [filename]
 * @returns {FormDataEntryValue}
 */
const toEntryValue = (value, filename) => {
  /* global File */
  if (value instanceof File) {
    return filename != null ? new File([value], filename, value) : value;
    /* global Blob */
  } else if (value instanceof Blob) {
    return new File([value], filename != null ? filename : "blob");
  } else {
    if (filename != null) {
      throw new TypeError(
        "filename is only supported when value is Blob or File"
      );
    }
    return value;
  }
};

export default class FormData {
  #entries;

  /**
   * FormData class
   *
   * @param {HTMLFormElement=} form
   */
  constructor(form) {
    /** @type {[string, string|File][]} */
    this.#entries = [];

    if (form !== undefined) {
      throw new TypeError("HTMLFormElement parameter is not supported.");
    }
  }

  /**
   * Appends a new value onto an existing key inside a FormData object, or adds
   * the key if it does not already exist.
   *
   * The difference between `set` and `append` is that if the specified key
   * already exists, `set` will overwrite all existing values with the new one,
   * whereas `append` will append the new value onto the end of the existing
   * set of values.
   *
   * @param {string} name
   * @param {string|Blob|File} value - The name of the field whose data is
   * contained in value.
   * @param {string} [filename] - The filename reported to the server, when a
   * value is a `Blob` or a `File`. The default filename for a `Blob` objects is
   * `"blob"`. The default filename for a `File` is the it's name.
   */
  append(name, value, filename) {
    if (value === undefined) {
      throw new TypeError("2 argument required, but only 1 present.");
    }
    this.#entries.push([name, toEntryValue(value, filename)]);
  }

  /**
   * Deletes a key and all its values from a FormData object.
   *
   * @param {string} name
   */
  delete(name) {
    if (name === undefined) {
      throw new TypeError("1 argument required, but none present.");
    }

    let index;
    if (
      (index = this.#entries.findIndex(([entryName]) => entryName === name)) >
      -1
    ) {
      this.#entries.splice(index, 1);
    }
  }

  /**
   * Returns the first value associated with a given key from within a
   * FormData object.
   *
   * @param {string} name
   * @returns {FormDataEntryValue|null}
   */
  get(name) {
    if (name === undefined) {
      throw new TypeError("1 argument required, but none present.");
    }

    for (const [entryName, value] of this.#entries) {
      if (entryName === name) {
        return value;
      }
    }
    return null;
  }

  /**
   * Returns an array of all the values associated with a given key from within
   * a FormData.
   *
   * @param {string} name
   * @returns {FormDataEntryValue[]}
   */
  getAll(name) {
    return this.#entries
      .map(([entryName, value]) => entryName === name && value)
      .filter(Boolean);
  }

  /**
   * Returns a boolean stating whether a FormData object contains a certain key.
   *
   * @param {string} name
   */
  has(name) {
    return this.#entries.some(([entryName]) => entryName === name);
  }

  /**
   * Sets a new value for an existing key inside a FormData object, or adds the
   * key/value if it does not already exist.
   *
   * @param {string} name
   * @param {string|Blob|File} value
   * @param {string} [filename]
   */
  set(name, value, filename) {
    if (value === undefined) {
      throw new TypeError("2 argument required, but 1 present.");
    }

    const entryValue = toEntryValue(value, filename);

    let index = 0;
    let wasSet = false;
    while (index < this.#entries.length) {
      const entry = this.#entries[index];

      if (entry[0] === name) {
        if (wasSet) {
          this.#entries.splice(index, 1);
        } else {
          wasSet = true;
          entry[1] = entryValue;
          index++;
        }
      } else {
        index++;
      }
    }

    if (!wasSet) {
      this.#entries.push([name, entryValue]);
    }
  }

  /**
   * Method returns an iterator allowing to go through all key/value pairs
   * contained in this object.
   */
  entries() {
    return this.#entries.values();
  }

  /**
   * Returns an iterator allowing to go through all keys of the key/value pairs
   * contained in this object.
   *
   * @returns {IterableIterator<string>}
   */
  *keys() {
    for (const [name] of this.#entries) {
      yield name;
    }
  }

  /**
   * Returns an iterator allowing to go through all values contained in this
   * object.
   *
   * @returns {IterableIterator<FormDataEntryValue>}
   */
  *values() {
    for (const [, value] of this.#entries) {
      yield value;
    }
  }

  /**
   * @param {(value: FormDataEntryValue, key: string, parent: globalThis.FormData) => void} fn
   * @param {any} [thisArg]
   * @returns {void}
   */
  forEach(fn, thisArg) {
    for (const [key, value] of this.#entries) {
      fn.call(thisArg, value, key, this);
    }
  }

  ["_blob"]() {
    const boundary = "----formdata-polyfill-" + Math.random();
    const p = `--${boundary}\r\nContent-Disposition: form-data; name="`;

    const chunks = [];
    this.forEach((value, name) => {
      name = encodeURIComponent(normalizeLineFeeds(name));

      return typeof value == "string"
        ? chunks.push(`${p}${name}"\r\n\r\n${normalizeLineFeeds(value)}\r\n`)
        : chunks.push(
            `${p}${name}"; filename="${encodeURIComponent(
              value.name
            )}"\r\nContent-Type: ${
              value.type || "application/octet-stream"
            }\r\n\r\n`,
            value,
            `\r\n`
          );
    });
    chunks.push(`--${boundary}--`);

    return new Blob(chunks, {
      type: "multipart/form-data; boundary=" + boundary,
    });
  }

  [Symbol.iterator]() {
    return this.#entries.values();
  }

  get [Symbol.toStringTag]() {
    return "FormData";
  }
}
