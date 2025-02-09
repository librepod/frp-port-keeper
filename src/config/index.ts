import 'dotenv/config'
import * as fs from 'fs'
import * as yaml from 'yaml'

const FALLBACK_PORT_RANGE = '8000-65535'

interface RangePort {
  start: number
  end: number
}

interface SinglePort {
  single: number
}

type PortRange = RangePort | SinglePort

export const PORT = process.env['PORT'] || '8080'
export let ALLOW_PORTS = process.env['ALLOW_PORTS'] || ''

if (!ALLOW_PORTS) {
  const frpsConfigPath = process.env['FRPS_CONFIG_PATH'] || './frps.yaml'

  try {
    const frpsConfigContent = fs.readFileSync(frpsConfigPath, 'utf-8')
    const frpsConfig = yaml.parse(frpsConfigContent)
    if (frpsConfig.allowPorts && frpsConfig.allowPorts.length > 0) {
      ALLOW_PORTS = frpsConfig.allowPorts
        .map((portRange: PortRange) => {
          if ('start' in portRange && 'end' in portRange) {
            return `${portRange.start}-${portRange.end}`
          } else if ('single' in portRange) {
            return `${portRange.single}`
          } else {
            return ''
          }
        })
        .filter(Boolean) // Remove empty strings
        .join(',')
      console.log(`Loaded allowPorts from ${frpsConfigPath}: ${ALLOW_PORTS}`)
    } else {
      console.log(
        `⚠ allowPorts not specified in ${frpsConfigPath}, falling back to ${FALLBACK_PORT_RANGE} port range`,
      )
      ALLOW_PORTS = FALLBACK_PORT_RANGE
    }
  } catch (err) {
    console.error(`Error reading ${frpsConfigPath}: ${err}`)
    // Fallback to default
    ALLOW_PORTS = FALLBACK_PORT_RANGE
  }
}
