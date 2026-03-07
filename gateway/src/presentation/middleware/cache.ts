// src/presentation/middleware/cache.ts
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

export function cacheControl() {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    await next();

    const path = new URL(c.req.url).pathname;

    // 静的アセットは長期キャッシュ
    if (path.startsWith('/assets/')) {
      c.res.headers.set('Cache-Control', 'public, max-age=31536000, immutable');
      return;
    }

    // フィード
    if (path.startsWith('/feed/')) {
      c.res.headers.set('Cache-Control', 'public, max-age=3600, s-maxage=3600');
      return;
    }

    // コンテンツページ: s-maxageでエッジキャッシュ、ブラウザには短いmax-age
    c.res.headers.set(
      'Cache-Control',
      'public, max-age=300, s-maxage=3600, stale-while-revalidate=86400'
    );
  });
}
