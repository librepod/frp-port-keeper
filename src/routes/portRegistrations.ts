import { Router } from 'express';
import { portRegistrationsHandler } from '../controllers/portRegistrationController';

const router = Router();

router.post('/', portRegistrationsHandler);

export default router;
