import NodeCache from 'node-cache';
import logger from '../utils/logger';

let cache: NodeCache;

export const initializeStore = () => {
  logger.info('Initializing store...');
  cache = new NodeCache();
};

export const getCache = (): NodeCache => {
  return cache;
};
