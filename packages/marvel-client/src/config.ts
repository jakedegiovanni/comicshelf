export interface Config {
  API_PUBLIC_KEY: string;
  API_PRIVATE_KEY: string;
  API_HOST: string;
}

export const configFromEnv = (): Config => {
  const { MARVEL_API_HOST, MARVEL_API_PRIVATE_KEY, MARVEL_API_PUBLIC_KEY } =
    process.env;

  if (!MARVEL_API_HOST || !MARVEL_API_PRIVATE_KEY || !MARVEL_API_PUBLIC_KEY) {
    throw new Error(
      `Environment must containe all of: MARVEL_API_HOST , MARVEL_API_PRIVATE_KEY , MARVEL_API_PUBLIC_KEY`,
    );
  }

  return {
    API_PUBLIC_KEY: MARVEL_API_PUBLIC_KEY,
    API_PRIVATE_KEY: MARVEL_API_PRIVATE_KEY,
    API_HOST: MARVEL_API_HOST,
  };
};
