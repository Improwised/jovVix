export function constructPath(basePath, params, queryParams) {
  let path = basePath;

  // Append parameters to the path
  if (params) {
    Object.keys(params).forEach((key) => {
      path += `/${params[key]}`;
    });
  }

  // Append query parameters to the path
  if (queryParams) {
    const queryString = new URLSearchParams(queryParams).toString();
    if (queryString) {
      path += `?${queryString}`;
    }
  }

  return path;
}

export function dividePath(fullPath) {
  if (!fullPath) {
    return {
      path: undefined,
      query: undefined,
    };
  }
  const pathAndQuery = fullPath.split("?");
  const path = pathAndQuery[0];
  const queryString = pathAndQuery.slice(1).join("?");
  const queryParams = {};

  if (queryString) {
    const params = new URLSearchParams(queryString);
    for (const [key, value] of params.entries()) {
      queryParams[key] = value;
    }
  }

  return {
    path,
    query: queryParams,
  };
}
