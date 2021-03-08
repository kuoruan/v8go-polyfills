import Response from "./Response";

async function fetch(url, opt = {}) {
  if (typeof url !== "string") url = String(url);
  /* global _goFetchSync */
  const res = await _goFetchSync(url, JSON.stringify(opt));
  return new Response(JSON.parse(res));
}

module.exports = fetch;
