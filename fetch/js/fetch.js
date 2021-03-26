import Response from "./Response";

async function fetch(url, opt) {
  /* global _goFetchSync */
  const res = await _goFetchSync(url, opt);
  return new Response(res);
}

module.exports = fetch;
