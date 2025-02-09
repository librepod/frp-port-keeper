import express from 'express'
import routes from './routes'
import { PORT, ALLOW_PORTS } from './config'
import logger from './utils/logger'
import { initPortsGenerator } from './utils/portGenerator'

const app = express()

// Initialize the ports generator with ALLOW_PORTS
initPortsGenerator(ALLOW_PORTS)

app.use(express.json())
app.use(routes)

app.listen(PORT, () => {
  logger.info(`Server is running on port ${PORT}`)
})

export default app
