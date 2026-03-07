// presentation/middleware/rate-limit.ts
// KVベースの簡易レート制限（D1を使用しない）
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

export function rateLimit(maxRequests: number = 100, windowSeconds: number = 60) {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    const ip = c.req.header('CF-Connecting-IP') ?? 'unknown';
    const key = `ratelimit:${ip}`;

    const current = await c.env.METADATA.get(key);
    const count = current ? parseInt(current, 10) : 0;

    if (count >= maxRequests) {
      return c.json(
        { error: 'Too many requests' },
        429,
        {
          'Retry-After': String(windowSeconds),
          'X-Gemini-Status': '44',
        }
      );
    }

    // カウンターをインクリメント（TTL付きで自動失効）
    await c.env.METADATA.put(key, String(count + 1), {
      expirationTtl: windowSeconds,
    });

    await next();
  });
}
