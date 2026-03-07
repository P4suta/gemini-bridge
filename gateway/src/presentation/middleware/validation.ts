// src/presentation/middleware/validation.ts
import { createMiddleware } from 'hono/factory';
import type { Env } from '../../config/bindings';

/** パスパラメータのバリデーション正規表現 */
const SLUG_PATTERN = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;

/** ディレクトリトラバーサル検出パターン */
const TRAVERSAL_PATTERN = /(?:^|[\/\\])\.\.(?:[\/\\]|$)/;

export function validateSlug() {
  return createMiddleware<{ Bindings: Env }>(async (c, next) => {
    const slug = c.req.param('slug');
    if (!slug) {
      await next();
      return;
    }

    // ディレクトリトラバーサル防止
    if (TRAVERSAL_PATTERN.test(slug)) {
      return c.json(
        { error: 'Invalid path parameter' },
        400,
        { 'X-Gemini-Status': '59' }
      );
    }

    // slugフォーマット検証
    if (!SLUG_PATTERN.test(slug)) {
      return c.json(
        { error: 'Invalid slug format. Only lowercase alphanumeric and hyphens allowed.' },
        400,
        { 'X-Gemini-Status': '59' }
      );
    }

    // 長さ制限（URLの実用的上限を考慮）
    if (slug.length > 128) {
      return c.json(
        { error: 'Slug too long' },
        400,
        { 'X-Gemini-Status': '59' }
      );
    }

    await next();
  });
}
