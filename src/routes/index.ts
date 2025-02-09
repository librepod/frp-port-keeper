import { Router } from 'express'
import pingRouter from './ping'
import portRegistrationsRouter from './portRegistrations'

const router = Router()

router.use('/ping', pingRouter)
router.use('/port-registrations', portRegistrationsRouter)

export default router
