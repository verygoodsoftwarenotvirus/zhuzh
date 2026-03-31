import { buildServerSideLogger } from '@zhuzh/logger';

const isDev = process.env.NODE_ENV !== 'production';

export const logger = buildServerSideLogger('admin', isDev);
