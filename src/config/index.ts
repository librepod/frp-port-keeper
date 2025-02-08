import 'dotenv/config'
import * as fs from 'fs';
import { parse } from 'ini';

export const PORT = process.env.PORT || '8080';
export let ALLOW_PORTS = process.env.ALLOW_PORTS || '';

if (!ALLOW_PORTS) {
  const frpsIniPath = process.env.FRPS_INI_PATH || './frps.ini';

  try {
    const frpsConfig = parse(fs.readFileSync(frpsIniPath, 'utf-8'));
    if (frpsConfig.common && frpsConfig.common.allow_ports) {
      ALLOW_PORTS = frpsConfig.common.allow_ports;
      console.log(`Loaded allow_ports from ${frpsIniPath}: ${ALLOW_PORTS}`);
    } else {
      console.log(
        `⚠ common.allow_ports not specified in ${frpsIniPath}, falling back to 1000-65535 port range`
      );
      ALLOW_PORTS = '1000-65535';
    }
  } catch (err) {
    console.error(`Error reading ${frpsIniPath}: ${err}`);
    // Fallback to default
    ALLOW_PORTS = '1000-65535';
  }
}
