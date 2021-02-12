export interface Config {
  debug: boolean;
}

const { env } = process;
const config: Config = {
  debug: env.DEBUG_OPERATOR?.toLowerCase() === 'true'
};

export default config;
