import { Router } from 'express'
import type { Request, Response } from 'express'

const router = Router()

router.get('/', (_: Request, res: Response) => {
  res.send('pong')
})

export default router
