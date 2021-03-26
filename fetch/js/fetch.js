import Response from "./Response";

async function fetch(url, opt) {
  if (typeof url !== "string") url = String(url);
  /* global _goFetchSync */
  const res = await _goFetchSync(url, opt);
  return new Response(res);
}

module.exports = fetch;
