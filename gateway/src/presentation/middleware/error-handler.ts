// src/presentation/middleware/error-handler.ts
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

export function errorHandler() {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    try {
      await next();
    } catch (error) {
      console.error('Unhandled error:', error);

      // Gemini 40 TEMPORARY FAILURE に相当
      c.res = new Response(
        '<h1>500 Internal Server Error</h1><p>An unexpected error occurred in the Gemini-HTTP bridge.</p>',
        {
          status: 500,
          headers: {
            'Content-Type': 'text/html; charset=utf-8',
            'X-Gemini-Status': '40',
            'X-Gemini-Meta': 'Internal bridge error',
          },
        }
      );
    }
  });
}
