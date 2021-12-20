export default class Body {
  #bodyInit;
  #options;

  #bodyText;
  #bodyBlob;
  #bodyFormData;

  bodyUsed = false;

  #initBody(body, options) {
    this.#bodyInit = body;

    if (!body) {
      this.#bodyText = "";
    } else if (typeof body === "string") {
      this.#bodyText = body;
      /* global Blob */
    } else if (body instanceof Blob) {
      this.#bodyBlob = body;
      this.#options = options;
      /* global FormData */
    } else if (body instanceof FormData) {
      this.#bodyFormData = body;
    } else {
      throw new TypeError("unsupported body type");
    }
  }

  async blob() {
    if (this.#bodyBlob) {
      return this.#bodyBlob;
    } else if (this.#bodyFormData) {
      throw new Error("could not read FormData body as blob");
    } else {
      return new Blob([this.#bodyText]);
    }
  }

  async text() {
    if (this.#bodyBlob) {
      return readBlobAsText(this._bodyBlob, this._options);
    } else if (this.#bodyFormData) {
      throw new Error("could not read FormData body as text");
    } else {
      return this.#bodyText;
    }
  }
}
