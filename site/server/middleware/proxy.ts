// /server/middleware/proxy.ts
import { proxyRequest, getRequestIP, getRequestProtocol } from 'h3'

const backendUrl = process.env.BBSGO_BACKEND_URL || 'http://localhost:3001/api/'

export default defineEventHandler((event) => {
  const req = event.node.req
  if (!req.url?.startsWith('/api')) return

  const target = new URL(req.url, backendUrl)

  return proxyRequest(event, target.toString(), {
    headers: {
      host: target.host,
      'X-Forwarded-For': getRequestIP(event, { xForwardedFor: true }) || '',
      'X-Forwarded-Proto': getRequestProtocol(event),
      'X-Forwarded-Host': req.headers['host'],
    },
  })
})
