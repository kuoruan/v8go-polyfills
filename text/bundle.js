(()=>{"use strict";var e={d:(t,r)=>{for(var o in r)e.o(r,o)&&!e.o(t,o)&&Object.defineProperty(t,o,{enumerable:!0,get:r[o]})},o:(e,t)=>Object.prototype.hasOwnProperty.call(e,t),r:e=>{"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})}},t={};e.r(t),e.d(t,{TextDecoder:()=>o,TextEncoder:()=>r});class r{encode(e,t={stream:!1}){if(t.stream)throw new Error("Failed to encode: the 'stream' option is unsupported.");let r=0;const o=e.length;let n=0,i=Math.max(32,o+(o>>>1)+7),a=new Uint8Array(i>>>3<<3);for(;r<o;){let t=e.charCodeAt(r++);if(t>=55296&&t<=56319){if(r<o){const o=e.charCodeAt(r);56320==(64512&o)&&(++r,t=((1023&t)<<10)+(1023&o)+65536)}if(t>=55296&&t<=56319)continue}if(n+4>a.length){i+=8,i*=1+r/e.length*2,i=i>>>3<<3;const t=new Uint8Array(i);t.set(a),a=t}if(0!=(4294967168&t)){if(0==(4294965248&t))a[n++]=t>>>6&31|192;else if(0==(4294901760&t))a[n++]=t>>>12&15|224,a[n++]=t>>>6&63|128;else{if(0!=(4292870144&t))continue;a[n++]=t>>>18&7|240,a[n++]=t>>>12&63|128,a[n++]=t>>>6&63|128}a[n++]=63&t|128}else a[n++]=t}return a.slice(0,n)}get encoding(){return"utf-8"}}class o{constructor(e="utf-8",t={fatal:!1}){if(!["utf-8","utf8","unicode-1-1-utf-8"].includes(e.toLowerCase()))throw new RangeError(`Failed to construct 'TextDecoder': The encoding label provided ('${e}') is invalid.`);if(t.fatal)throw new Error("Failed to construct 'TextDecoder': the 'fatal' option is unsupported.")}decode(e,t={stream:!1}){if(t.stream)throw new Error("Failed to decode: the 'stream' option is unsupported.");let r;r=e instanceof Uint8Array?e:e.buffer instanceof ArrayBuffer?new Uint8Array(e.buffer):new Uint8Array(e);let o=0;const n=Math.min(65536,r.length+1),i=new Uint16Array(n),a=[];let l=0;for(;;){const e=o<r.length;if(!e||l>=n-1){if(a.push(String.fromCharCode.apply(null,i.subarray(0,l))),!e)return a.join("");r=r.subarray(o),o=0,l=0}const t=r[o++];if(0==(128&t))i[l++]=t;else if(192==(224&t)){const e=63&r[o++];i[l++]=(31&t)<<6|e}else if(224==(240&t)){const e=63&r[o++],n=63&r[o++];i[l++]=(31&t)<<12|e<<6|n}else if(240==(248&t)){let e=(7&t)<<18|(63&r[o++])<<12|(63&r[o++])<<6|63&r[o++];e>65535&&(e-=65536,i[l++]=e>>>10&1023|55296,e=56320|1023&e),i[l++]=e}}}get encoding(){return"utf-8"}get fatal(){return!1}get ignoreBOM(){return!1}}var n=window;for(var i in t)n[i]=t[i];t.__esModule&&Object.defineProperty(n,"__esModule",{value:!0})})();