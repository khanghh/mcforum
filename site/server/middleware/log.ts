// /server/middleware/logger.ts
import { consola } from 'consola'
import { getRequestURL, getRequestIP } from 'h3'

export default defineEventHandler((event) => {
  event.node.res.on('finish', () => {
    consola.info(
      event.node.req.method,
      getRequestURL(event).pathname,
      event.node.res.statusCode,
      getRequestIP(event, { xForwardedFor: true })
    )
  })
})
