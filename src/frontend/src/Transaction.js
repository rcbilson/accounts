export function transactionQuery(querySpec) {
  let specs = []
  if (querySpec.dateFrom) {
    specs.push("DateFrom=" + querySpec.dateFrom);
  }
  if (querySpec.dateUntil) {
    specs.push("DateUntil=" + querySpec.dateUntil);
  }
  if (querySpec.descrLike !== "") {
    specs.push("DescrLike=" + querySpec.descrLike);
  }
  if (querySpec.category !== "") {
    specs.push("Category=" + querySpec.category);
  }
  if (querySpec.subcategory !== "") {
    specs.push("Subcategory=" + querySpec.subcategory);
  }
  let path = "/api/transactions";
  if (specs.length > 0) {
    path += "?" + specs[0]
    specs.slice(1).forEach((e) => { path += "&" + e })
  } else {
    // no filter, limit to the first 50
    path += "?Limit=50"
  }
  return fetch(encodeURI(path))
    .then(res => res.json())
}
