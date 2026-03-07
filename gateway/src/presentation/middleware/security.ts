// src/presentation/middleware/security.ts
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

export function securityHeaders() {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    await next();

    c.res.headers.set('X-Content-Type-Options', 'nosniff');
    c.res.headers.set('X-Frame-Options', 'DENY');
    c.res.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');
    c.res.headers.set(
      'Content-Security-Policy',
      "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'"
    );
    c.res.headers.set(
      'Strict-Transport-Security',
      'max-age=31536000; includeSubDomains'
    );
    c.res.headers.set(
      'Permissions-Policy',
      'camera=(), microphone=(), geolocation=()'
    );
  });
}
