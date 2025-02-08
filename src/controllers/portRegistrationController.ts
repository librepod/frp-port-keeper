import { RequestHandler } from 'express';
import { getNextPort } from '../utils/portGenerator';
import { getCache } from '../services/store';
import logger from '../utils/logger';

export const portRegistrationsHandler: RequestHandler = (req, res) => {
  try {
    const { proxy_name } = req.body.content;
    
    // Check if port already assigned for this proxy
    const cache = getCache();
    const cachedPort = cache.get<number>(proxy_name);
    
    if (cachedPort) {
      logger.info(`Found cached port ${cachedPort} for proxy ${proxy_name}`);
      res.json({
        unchange: false,
        content: {
          ...req.body.content,
          remote_port: cachedPort
        }
      });
      return;
    }

    // Get new port
    const port = getNextPort();
    if (!port) {
      res.status(400).json({
        reject: true,
        reject_reason: 'NO_MORE_FREE_PORTS_LEFT'
      });
      return;
    }

    // Cache the port
    cache.set(proxy_name, port);
    cache.set(port.toString(), proxy_name);

    logger.info(`Assigned port ${port} to proxy ${proxy_name}`);
    
    res.json({
      unchange: false,
      content: {
        ...req.body.content,
        remote_port: port
      }
    });
    return;
  } catch (error) {
    logger.error(`Error in portRegistrationsHandler: ${error}`);
    res.status(500).json({
      reject: true,
      reject_reason: 'INTERNAL_SERVER_ERROR'
    });
  }
};
