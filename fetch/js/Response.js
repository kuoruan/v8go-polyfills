class Response {
  constructor(res) {
    this.headers = res.headers;
    this.ok = res.ok;
    this.redirected = res.redirected;
    this.status = res.status;
    this.statusText = res.statusText;
    this.url = res.url;
    this.body = res.body;
  }

  text() {
    if (typeof this.body !== "string") {
      return Promise.reject("response body is not text.");
    }
    return Promise.resolve(this.body);
  }

  json() {
    return this.text().then((v) => JSON.parse(v));
  }
}

module.exports = Response;
