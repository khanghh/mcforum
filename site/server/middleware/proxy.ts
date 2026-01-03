// /server/middleware/proxy.ts
import { consola } from 'consola'

import { proxyRequest } from 'h3'

export default defineEventHandler((event) => {
  const url = event.node.req.url
  consola.info('Received request for URL:', url)
  consola.log("Condition:", !url?.startsWith('/api/'))
  if (!url?.startsWith('/api/')) return

  const { backendUrl } = useRuntimeConfig().app
  const target = new URL(url, backendUrl)

  const req = event.node.req

  // Resolve client IP
  const clientIp =
    (req.headers['x-forwarded-for'] as string)?.split(',')[0] ||
    req.socket.remoteAddress ||
    ''
  consola.info('Proxying request to:', target.toString(), 'from client IP:', clientIp)

  return proxyRequest(event, target.toString(), {
    headers: {
      host: target.host,
      'X-Forwarded-For': clientIp,
      'X-Forwarded-Proto': req.headers['x-forwarded-proto'] as string || 'http',
      'X-Forwarded-Host': req.headers['host'] || '',
    },
  })
})
