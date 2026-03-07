// src/presentation/middleware/analytics.ts
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

export function analyticsTracker() {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    await next();

    // Analytics Engineにイベントを非同期記録（レスポンスをブロックしない）
    c.executionCtx.waitUntil(
      recordEvent(c.env.ANALYTICS, {
        path: c.req.path,
        userAgent: c.req.header('User-Agent') ?? '',
        country: c.req.header('CF-IPCountry') ?? '',
        referer: c.req.header('Referer') ?? '',
        timestamp: Date.now(),
      })
    );
  });
}

async function recordEvent(
  analytics: AnalyticsEngineDataset,
  event: {
    path: string;
    userAgent: string;
    country: string;
    referer: string;
    timestamp: number;
  }
): Promise<void> {
  analytics.writeDataPoint({
    blobs: [event.path, event.userAgent, event.country, event.referer],
    doubles: [event.timestamp],
    indexes: [event.path],
  });
}
