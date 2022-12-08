export function Query(querySpec) {
  let specs = []
  if (querySpec.dateFrom) {
    specs.push("DateFrom=" + querySpec.dateFrom);
  }
  if (querySpec.dateUntil) {
    specs.push("DateUntil=" + querySpec.dateUntil);
  }
  if (querySpec.descrLike && querySpec.descrLike !== "") {
    specs.push("DescrLike=" + querySpec.descrLike);
  }
  if (querySpec.category && querySpec.category !== "") {
    specs.push("Category=" + querySpec.category);
  }
  if (querySpec.subcategory && querySpec.subcategory !== "") {
    specs.push("Subcategory=" + querySpec.subcategory);
  }
  if (querySpec.state && querySpec.state !== "") {
    specs.push("State=" + querySpec.state);
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

export async function Delete(id) {
  const response = await fetch(`/api/transactions/${id}`, {
    method: 'DELETE',
  });
  if (response.length > 0) {
    const r = response.json();
    if (r.message) {
      throw new Error(r.message);
    }
  }
}

export async function Update(t) {
  if (!t?.id) {
    throw new Error(`No id specified in ${t}`)
  }

  const response = await fetch(`/api/transactions/${t.id}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(t)
  });
  if (response.length > 0) {
    const r = response.json();
    if (r.message) {
      throw new Error(r.message);
    }
  }
}

export async function Import(blob) {
  const response = await fetch(`/api/import`, {
    method: 'POST',
    headers: {
      'Content-Type': 'text/csv'
    },
    body: blob
  });
  if (response.length > 0) {
    const r = response.json();
    if (r.message) {
      throw new Error(r.message);
    }
  }
}
