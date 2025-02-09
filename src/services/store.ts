import { LowSync } from 'lowdb'
import { JSONFileSyncPreset } from 'lowdb/node'
import logger from '../utils/logger'
import { DB_PATH } from '../config'

type Data = {
  proxies: { [key: string]: number }
  ports: { [key: number]: string }
}

export const getDb = (): LowSync<Data> => {
  logger.info('Initializing store...')
  const defaultData: Data = {
    proxies: {},
    ports: {},
  }
  return JSONFileSyncPreset(DB_PATH, defaultData)
}
