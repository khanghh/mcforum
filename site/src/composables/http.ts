import { useCookie } from '#app'

// 请求体封装
function applyOptions(options: any = {}): any {
  const config = useRuntimeConfig()
  options.baseURL = config.public.BASE_URL as string
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

export function useHttp(url: string, options: any = {}): Promise<any> {
  options = applyOptions(options)
  if (options.headers?.['Content-Type'] === 'application/x-www-form-urlencoded' && options.body) {
    options.body = new URLSearchParams(options.body).toString()
  }
  return new Promise((resolve, reject) => {
    $fetch(url, options)
      .then((resp) => {
        const response = resp as { error?: any; data?: any }
        if (!response.error) {
          resolve(response.data)
        } else {
          reject(response.error)
        }
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export function useHttpPost(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'POST',
  })
}

export function useHttpGet(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'GET',
  })
}

export function useHttpDelete(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'DELETE',
  })
}

export function useHttpPutForm(url: string, options: any = {}): Promise<any> {
  console.log(options)
  return useHttp(url, {
    ...options,
    method: 'PUT',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
  })
}

export function useHttpPatchForm(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
  })
}

// POST请求(application/x-www-form-urlencoded)
export function useHttpPostForm(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
  })
}

// POST请求(multipart/form-data)
export function useHttpPostMultipart(url: string, formData: FormData): Promise<any> {
  return useHttp(url, {
    method: 'POST',
    body: formData,
  })
}
