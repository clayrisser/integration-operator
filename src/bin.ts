import IntegrationOperator from './integrationOperator';
import config from './config';

const logger = console;

(async () => {
  const integrationOperator = new IntegrationOperator(config);
  function exit(_reason: string) {
    integrationOperator.stop();
    process.exit(0);
  }
  process
    .on('SIGTERM', () => exit('SIGTERM'))
    .on('SIGINT', () => exit('SIGINT'));
  await integrationOperator.start();
})().catch(logger.error);
