const CONFIG = {
  baseURL: process.env.BASE_URL_TEST,
}

// 请求体封装
function applyOptions(options = {}) {
  options.baseURL = options.baseURL ?? CONFIG.baseURL
  options.initialCache = options.initialCache ?? false
  options.headers = options.headers || {}
  options.method = options.method || 'GET'

  const token = useCookie('bbsgo_token')
  if (token.value) {
    options.headers['X-User-Token'] = token.value
  }

  // options.params = options.params || {}
  // options.params.userToken = token.value

  return options
}

export function useHttp(url, options = {}) {
  options = applyOptions(options)
  if (options.headers?.['Content-Type'] === 'application/x-www-form-urlencoded' && options.body) {
    options.body = new URLSearchParams(options.body).toString()
  }
  return new Promise((resolve, reject) => {
    $fetch(url, options)
      .then((resp) => {
        if (!resp.error) {
          resolve(resp.data)
        } else {
          reject(resp)
        }
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export function useHttpPost(url, options = {}) {
  return useHttp(url, {
    ...options,
    method: 'POST',
  })
}

export function useHttpGet(url, options = {}) {
  return useHttp(url, {
    ...options,
    method: 'GET',
  })
}

export function useHttpDelete(url, options = {}) {
  return useHttp(url, {
    ...options,
    method: 'DELETE',
  })
}

export function useHttpPutForm(url, { body } = {}) {
  return useHttp(url, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body,
  })
}

export function useHttpPatchForm(url, { body } = {}) {
  return useHttp(url, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body,
  })
}

// POST请求(application/x-www-form-urlencoded)
export function useHttpPostForm(url, { body = {} } = {}) {
  return useHttp(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body,
  })
}

// POST请求(multipart/form-data)
export function useHttpPostMultipart(url, formData) {
  return useHttp(url, {
    method: 'POST',
    body: formData,
  })
}
