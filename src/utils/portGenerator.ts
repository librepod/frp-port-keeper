import logger from './logger'

let portIterator: Generator<number, void, unknown>

export const initPortsGenerator = (allowPorts: string) => {
  portIterator = createAllowPortsGenerator(allowPorts)
}

export const getNextPort = (): number | null => {
  const result = portIterator.next()
  return result.done ? null : result.value
}

const createAllowPortsGenerator = function* (portsRange: string) {
  logger.info('Initializing ports generator...')
  const rangeSlices = portsRange.split(',')
  const ranges: Array<{ start: number; end: number }> = []

  for (const r of rangeSlices) {
    if (r.includes('-')) {
      const [startStr, endStr] = r.split('-').map(s => s.trim())
      const start = parseInt(startStr, 10)
      const end = parseInt(endStr, 10)

      if (start > end) {
        throw new Error('😱 invalid range supplied')
      }

      ranges.push({ start, end })
    } else {
      const port = parseInt(r.trim(), 10)
      ranges.push({ start: port, end: port })
    }
  }

  for (const range of ranges) {
    for (let i = range.start; i <= range.end; i++) {
      yield i
    }
  }
}
