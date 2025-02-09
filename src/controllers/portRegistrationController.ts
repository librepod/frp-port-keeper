import type { RequestHandler } from 'express'
import { getNextPort } from '../utils/portGenerator'
import { getDb } from '../services/store'
import logger from '../utils/logger'

export const portRegistrationsHandler: RequestHandler = async (req, res) => {
  try {
    const { proxy_name } = req.body.content

    // Check if port already assigned for this proxy
    const db = getDb()
    // console.log('db:', db)
    const cachedPort = db.data.proxies[proxy_name]
    console.log('cachedPort:', cachedPort)

    if (cachedPort) {
      logger.info(`Found cached port ${cachedPort} for proxy ${proxy_name}`)
      res.json({
        unchange: false,
        content: {
          ...req.body.content,
          remote_port: cachedPort,
        },
      })
      return
    }

    // Get new port
    let port = getNextPort()
    while (port !== null && db.data.ports[port]) {
      // Port is already taken, get the next port
      port = getNextPort()
    }

    if (port === null) {
      res.status(400).json({
        reject: true,
        reject_reason: 'NO_MORE_FREE_PORTS_LEFT',
      })
      return
    }

    // Save the proxy_name & port
    db.data.proxies[proxy_name] = port
    db.data.ports[port] = proxy_name
    db.write()

    logger.info(`Assigned port ${port} to proxy ${proxy_name}`)

    res.json({
      unchange: false,
      content: {
        ...req.body.content,
        remote_port: port,
      },
    })
    return
  } catch (error) {
    logger.error(`Error in portRegistrationsHandler: ${error}`)
    res.status(500).json({
      reject: true,
      reject_reason: 'INTERNAL_SERVER_ERROR',
    })
  }
}
