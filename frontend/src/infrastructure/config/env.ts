import { z } from 'zod'

const envSchema = z.object({
  VITE_API_URL: z.string().url().default('http://localhost:8080'),
  VITE_DEFAULT_MODEL: z.string().default('flux'),
  MODE: z.enum(['development', 'production', 'test']).default('development'),
})

export type Env = z.infer<typeof envSchema>

export const validateEnv = (): Env => {
  const env = {
    VITE_API_URL: import.meta.env.VITE_API_URL,
    VITE_DEFAULT_MODEL: import.meta.env.VITE_DEFAULT_MODEL,
    MODE: import.meta.env.MODE,
  }

  const result = envSchema.safeParse(env)

  if (!result.success) {
    console.error('❌ Invalid environment variables:', result.error.format())
    throw new Error('Invalid environment variables')
  }

  return result.data
}

export const env = validateEnv()
