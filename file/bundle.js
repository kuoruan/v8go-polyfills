(()=>{"use strict";var e={d:(t,r)=>{for(var n in r)e.o(r,n)&&!e.o(t,n)&&Object.defineProperty(t,n,{enumerable:!0,get:r[n]})},o:(e,t)=>Object.prototype.hasOwnProperty.call(e,t),r:e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})}},t={};async function*r(e,t=!0){for(const r of e)if("stream"in r)yield*r.stream();else if(ArrayBuffer.isView(r))if(t){let e=r.byteOffset;const t=r.byteOffset+r.byteLength;for(;e!==t;){const n=Math.min(t-e,65536),o=r.buffer.slice(e,e+n);e+=o.byteLength,yield new Uint8Array(o)}}else yield r;else{let e=0;for(;e!==r.size;){const t=r.slice(e,Math.min(r.size,e+65536)),n=await t.arrayBuffer();e+=n.byteLength,yield new Uint8Array(n)}}}e.r(t),e.d(t,{default:()=>i});class n{#e=[];#t="";#r=0;constructor(e=[],t={}){if("object"!=typeof e||null===e)throw new TypeError("Failed to construct 'Blob': The provided value cannot be converted to a sequence.");if("function"!=typeof e[Symbol.iterator])throw new TypeError("Failed to construct 'Blob': The object must have a callable @@iterator property.");if("object"!=typeof t&&"function"!=typeof t)throw new TypeError("Failed to construct 'Blob': parameter 2 cannot convert to dictionary.");const r=new TextEncoder;for(const t of e){let e;e=ArrayBuffer.isView(t)?new Uint8Array(t.buffer.slice(t.byteOffset,t.byteOffset+t.byteLength)):t instanceof ArrayBuffer?new Uint8Array(t.slice(0)):t instanceof n?t:r.encode(t),this.#r+=ArrayBuffer.isView(e)?e.byteLength:e.size,this.#e.push(e)}const o=void 0===t?.type?"":String(t.type);this.#t=/^[\x20-\x7E]*$/.test(o)?o:""}get size(){return this.#r}get type(){return this.#t}async text(){const e=new TextDecoder;let t="";for await(const n of r(this.#e,!1))t+=e.decode(n,{stream:!0});return t+=e.decode(),t}async arrayBuffer(){const e=new Uint8Array(this.size);let t=0;for await(const n of r(this.#e,!1))e.set(n,t),t+=n.length;return e.buffer}stream(){const e=r(this.#e,!0);return new ReadableStream({type:"bytes",async pull(t){const r=await e.next();r.done?t.close():t.enqueue(r.value)},async cancel(){await e.return()}})}slice(e=0,t=this.size,r=""){const{size:o}=this;let i=e<0?Math.max(o+e,0):Math.min(e,o),s=t<0?Math.max(o+t,0):Math.min(t,o);const a=Math.max(s-i,0),f=this.#e,c=[];let l=0;for(const e of f){if(l>=a)break;const t=ArrayBuffer.isView(e)?e.byteLength:e.size;if(i&&t<=i)i-=t,s-=t;else{let r;ArrayBuffer.isView(e)?(r=e.subarray(i,Math.min(t,s)),l+=r.byteLength):(r=e.slice(i,Math.min(t,s)),l+=r.size),s-=t,c.push(r),i=0}}const y=new n([],{type:String(r).toLowerCase()});return y.#r=a,y.#e=c,y}get[Symbol.toStringTag](){return"Blob"}static[Symbol.hasInstance](e){return e&&"object"==typeof e&&"function"==typeof e.constructor&&("function"==typeof e.stream||"function"==typeof e.arrayBuffer)&&/^(Blob|File)$/.test(e[Symbol.toStringTag])}}class o extends n{#n=0;#o="";constructor(e,t,r={}){if(arguments.length<2)throw new TypeError(`Failed to construct 'File': 2 arguments required, but only ${arguments.length} present.`);super(e,r),this.#n=Number(r?.lastModified)||Date.now(),this.#o=String(t)}get name(){return this.#o}get lastModified(){return this.#n}get webkitRelativePath(){return""}get[Symbol.toStringTag](){return"File"}}const i={Blob:n,File:o};var s=window;for(var a in t)s[a]=t[a];t.__esModule&&Object.defineProperty(s,"__esModule",{value:!0})})();