// src/infrastructure/kv-cache.ts
import type { ContentCache } from '../application/serve-content';

export class KVContentCache implements ContentCache {
  constructor(private readonly kv: KVNamespace) {}

  async get(key: string): Promise<string | null> {
    return await this.kv.get(key);
  }

  async set(key: string, value: string, ttlSeconds: number): Promise<void> {
    await this.kv.put(key, value, { expirationTtl: ttlSeconds });
  }
}
