import { Low } from 'lowdb';
import { JSONFileSyncPreset } from 'lowdb/node';
import logger from '../utils/logger';

type Data = {
  proxies: { [key: string]: number };
};

let db: Low<Data>;

export const initializeStore = async () => {
  logger.info('Initializing store...');
  const defaultData: Data = {
    proxies: {}
  }
  return JSONFileSyncPreset('db.json', defaultData)
};

export const getDb = (): Low<Data> => {
  return db;
};
