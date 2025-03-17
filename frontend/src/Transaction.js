export async function doQuery(baseUrl, querySpec, defaultQueryString) {
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
  if (querySpec.limit) {
    specs.push("Limit=" + querySpec.limit);
  }
  let path = baseUrl;
  if (specs.length > 0) {
    path += "?" + specs[0]
    specs.slice(1).forEach((e) => { path += "&" + e })
  } else if (defaultQueryString) {
    path += defaultQueryString;
  }
  return fetch(encodeURI(path))
    .then(res => res.json())
}

export async function Query(querySpec) {
  return doQuery("/api/transactions", querySpec, "?Limit=50");
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

export async function Categories(querySpec) {
  return doQuery("/api/categories", querySpec);
}

export async function Summary(querySpec) {
  return doQuery("/api/summary", querySpec);
}

export async function SummaryChart(summaryChartSpec) {
  let specs = []
  if (summaryChartSpec.dateFrom) {
    specs.push("DateFrom=" + summaryChartSpec.dateFrom);
  }
  if (summaryChartSpec.dateUntil) {
    specs.push("DateUntil=" + summaryChartSpec.dateUntil);
  }
  if (summaryChartSpec.chartType) {
    specs.push("ChartType=" + summaryChartSpec.chartType);
  }
  let path = "/api/summaryChart";
  if (specs.length > 0) {
    path += "?" + specs[0]
    specs.slice(1).forEach((e) => { path += "&" + e })
  } else if (defaultQueryString) {
    path += defaultQueryString;
  }
  return fetch(encodeURI(path))
    .then(res => res.json())
  }