export default class TextDecoder {
  constructor(utfLabel = "utf-8", options = { fatal: false }) {
    if (
      !["utf-8", "utf8", "unicode-1-1-utf-8"].includes(utfLabel.toLowerCase())
    ) {
      throw new RangeError(
        `Failed to construct 'TextDecoder': The encoding label provided ('${utfLabel}') is invalid.`
      );
    }
    if (options.fatal) {
      throw new Error(
        `Failed to construct 'TextDecoder': the 'fatal' option is unsupported.`
      );
    }
  }

  decode(buf, options = { stream: false }) {
    if (options["stream"]) {
      throw new Error(`Failed to decode: the 'stream' option is unsupported.`);
    }

    let bytes;

    if (buf instanceof Uint8Array) {
      // Accept Uint8Array instances as-is.
      bytes = buf;
    } else if (buf.buffer instanceof ArrayBuffer) {
      // Look for ArrayBufferView, which isn't a real type, but basically
      // represents all the valid TypedArray types plus DataView. They all have
      // ".buffer" as an instance of ArrayBuffer.
      bytes = new Uint8Array(buf.buffer);
    } else {
      // The only other valid argument here is that "buffer" is an ArrayBuffer.
      // We also try to convert anything else passed to a Uint8Array, as this
      // catches anything that's array-like. Native code would throw here.
      bytes = new Uint8Array(buf);
    }

    let inputIndex = 0;

    // Create a working buffer for UTF-16 code points, but don't generate one
    // which is too large for small input sizes. UTF-8 to UCS-16 conversion is
    // going to be at most 1:1, if all code points are ASCII. The other extreme
    // is 4-byte UTF-8, which results in two UCS-16 points, but this is still 50%
    // fewer entries in the output.
    const pendingSize = Math.min(256 * 256, bytes.length + 1);
    const pending = new Uint16Array(pendingSize);
    const chunks = [];
    let pendingIndex = 0;

    for (;;) {
      const more = inputIndex < bytes.length;

      // If there's no more data or there'd be no room for two UTF-16 values,
      // create a chunk. This isn't done at the end by simply slicing the data
      // into equal sized chunks as we might hit a surrogate pair.
      if (!more || pendingIndex >= pendingSize - 1) {
        // nb. .apply and friends are *really slow*. Low-hanging fruit is to
        // expand this to literally pass pending[0], pending[1], ... etc, but
        // the output code expands pretty fast in this case.
        chunks.push(
          String.fromCharCode.apply(null, pending.subarray(0, pendingIndex))
        );

        if (!more) {
          return chunks.join("");
        }

        // Move the buffer forward and create another chunk.
        bytes = bytes.subarray(inputIndex);
        inputIndex = 0;
        pendingIndex = 0;
      }

      // The native TextDecoder will generate "REPLACEMENT CHARACTER" where the
      // input data is invalid. Here, we blindly parse the data even if it's
      // wrong: e.g., if a 3-byte sequence doesn't have two valid continuations.

      const byte1 = bytes[inputIndex++];
      if ((byte1 & 0x80) === 0) {
        // 1-byte or null
        pending[pendingIndex++] = byte1;
      } else if ((byte1 & 0xe0) === 0xc0) {
        // 2-byte
        const byte2 = bytes[inputIndex++] & 0x3f;
        pending[pendingIndex++] = ((byte1 & 0x1f) << 6) | byte2;
      } else if ((byte1 & 0xf0) === 0xe0) {
        // 3-byte
        const byte2 = bytes[inputIndex++] & 0x3f;
        const byte3 = bytes[inputIndex++] & 0x3f;
        pending[pendingIndex++] = ((byte1 & 0x1f) << 12) | (byte2 << 6) | byte3;
      } else if ((byte1 & 0xf8) === 0xf0) {
        // 4-byte
        const byte2 = bytes[inputIndex++] & 0x3f;
        const byte3 = bytes[inputIndex++] & 0x3f;
        const byte4 = bytes[inputIndex++] & 0x3f;

        // this can be > 0xffff, so possibly generate surrogates
        let codepoint =
          ((byte1 & 0x07) << 0x12) | (byte2 << 0x0c) | (byte3 << 0x06) | byte4;
        if (codepoint > 0xffff) {
          // codepoint &= ~0x10000;
          codepoint -= 0x10000;
          pending[pendingIndex++] = ((codepoint >>> 10) & 0x3ff) | 0xd800;
          codepoint = 0xdc00 | (codepoint & 0x3ff);
        }
        pending[pendingIndex++] = codepoint;
      } else {
        // invalid initial byte
      }
    }
  }

  get encoding() {
    return "utf-8";
  }

  get fatal() {
    return false;
  }

  get ignoreBOM() {
    return false;
  }
}
