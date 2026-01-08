import { useCookie } from '#app'

// Request body encapsulation
function applyOptions(options: any = {}): any {
  const config = useRuntimeConfig()
  options.baseURL = config.public.BASE_URL as string
  options.initialCache = options.initialCache ?? false
  options.headers = options.headers || {}
  options.method = options.method || 'GET'

  const token = useCookie('sid')
  if (token.value) {
    options.headers['X-User-Token'] = token.value
  }

  // options.params = options.params || {}
  // options.params.userToken = token.value

  return options
}

export async function useHttp(url: string, options: any = {}): Promise<any> {
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

export function useHttpPut(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'PUT',
  })
}

export function useHttpPatch(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'PATCH',
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

// POST request (application/x-www-form-urlencoded)
export function useHttpPostForm(url: string, options: any = {}): Promise<any> {
  return useHttp(url, {
    ...options,
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
  })
}

function upload(url: string, options: any = {}): Promise<any> {
  options = applyOptions(options)
  const fullUrl = options.baseURL ? (options.baseURL + url) : url
  const formData = options.body
  const onUploadProgress = options.onUploadProgress
  const method = options.method

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open(method, fullUrl)

    if (options.headers) {
      Object.keys(options.headers).forEach((key) => {
        // FormData automatically sets Content-Type with boundary
        if (key.toLowerCase() !== 'content-type') {
          xhr.setRequestHeader(key, options.headers[key])
        }
      })
    }

    if (onUploadProgress) {
      xhr.upload.onprogress = onUploadProgress
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const response = JSON.parse(xhr.responseText)
          if (!response.error) {
            resolve(response.data)
          } else {
            reject(response.error)
          }
        } catch (e) {
          reject(e)
        }
      } else {
        reject(new Error(xhr.statusText || 'Upload failed'))
      }
    }

    xhr.onerror = () => {
      reject(new Error('Network Error'))
    }

    xhr.send(formData)
  })
}

// File upload
export function useHttpPostUpload(url: string, options: any = {}): Promise<any> {
  return upload(url, { ...options, method: 'POST' })
}

// File upload (PATCH)
export function useHttpPatchUpload(url: string, options: any = {}): Promise<any> {
  return upload(url, { ...options, method: 'PATCH' })
}

