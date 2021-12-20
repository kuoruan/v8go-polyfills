/*
 * Copyright (c) 2021 Xingwang Liao
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

export default class Headers {
  #map;

  /**
   * Headers class
   *
   * @param headers
   */
  constructor(headers) {
    this.#map = Object.create(null);

    if (headers) {
      if (typeof headers[Symbol.iterator] === "function") {
        for (let [name, value] of headers) {
          this.append(name, value);
        }
      } else {
        for (let name of Object.keys(headers)) {
          this.append(name, headers[name]);
        }
      }
    }
  }

  /**
   *  Append a value onto existing header
   *
   * @param name the header name
   * @param value the header value
   */
  append(name, value) {
    name = normalizeName(name);
    value = normalizeValue(value);

    if (!this.#map[name]) {
      this.#map[name] = [];
    }

    this.#map[name].push(value);
  }

  /**
   * Delete all header values given name
   *
   * @param name name  Header name
   */
  delete(name) {
    delete this.#map[normalizeName(name)];
  }

  /**
   * Iterate over all headers as [name, value]
   *
   * @return  Iterator
   */
  *entries() {
    for (const name in this.#map) {
      yield [name, this.#map[name].join(",")];
    }
  }

  /**
   * Return first header value given name
   *
   * @param name  Header name
   * @return Mixed
   */
  get(name) {
    name = normalizeName(name);

    return this.#map[name] ? this.#map[name][0] : null;
  }

  /**
   * Check for header name existence
   *
   * @param name  Header name
   * @return Boolean
   */
  has(name) {
    return normalizeName(name) in this.#map;
  }

  /**
   * Iterate over all keys
   *
   * @return  Iterator
   */
  *keys() {
    for (let [name] of this) {
      yield name;
    }
  }

  /**
   * Overwrite header values given name
   *
   * @param name The header name
   * @param value The header value
   */
  set(name, value) {
    this.#map[normalizeName(name)] = [normalizeValue(value)];
  }

  /**
   * Iterate over all values
   *
   * @return  Iterator
   */
  *values() {
    for (const [, value] of this) {
      yield value;
    }
  }

  /**
   * The class itself is iterable
   * alies for headers.entries()
   *
   * @return  Iterator
   */
  [Symbol.iterator]() {
    return this.entries();
  }

  /**
   * Create the default string description.
   * It is accessed internally by the Object.prototype.toString().
   *
   * @return  String  [Object Headers]
   */
  get [Symbol.toStringTag]() {
    return "Headers";
  }
}

function normalizeName(name) {
  if (typeof name !== "string") {
    name = String(name);
  }

  if (/[^a-z0-9\-#$%&'*+.^_`|~]/i.test(name)) {
    throw new TypeError("Invalid character in header field name");
  }

  return name.toLowerCase();
}

function normalizeValue(value) {
  if (typeof value !== "string") {
    value = String(value);
  }

  return value;
}
