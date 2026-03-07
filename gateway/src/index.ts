// src/index.ts （DI配線を含む完全版）
import { Hono } from 'hono';
import type { Env } from './config/bindings';
import { R2ContentRepository } from './infrastructure/r2-repository';
import { KVContentCache } from './infrastructure/kv-cache';
import { D1Repository } from './infrastructure/d1-repository';
import { WorkersAIClient } from './infrastructure/workers-ai';
import { securityHeaders } from './presentation/middleware/security';
import { cacheControl } from './presentation/middleware/cache';
import { analyticsTracker } from './presentation/middleware/analytics';
import { errorHandler } from './presentation/middleware/error-handler';
import { rateLimit } from './presentation/middleware/rate-limit';
import { serveContent } from './application/serve-content';

// コンテキスト型拡張: DI済みインスタンスをHonoコンテキストに注入
type AppVariables = {
  repository: R2ContentRepository;
  cache: KVContentCache;
  db: D1Repository;
  ai: WorkersAIClient;
};

const app = new Hono<{ Bindings: Env; Variables: AppVariables }>();

// DI配線ミドルウェア: バインディングからインスタンスを生成しコンテキストに注入
app.use('*', async (c, next) => {
  c.set('repository', new R2ContentRepository(c.env.CONTENT));
  c.set('cache', new KVContentCache(c.env.METADATA));
  c.set('db', new D1Repository(c.env.DB));
  c.set('ai', new WorkersAIClient(c.env.AI));
  await next();
});

// グローバルミドルウェア
app.use('*', errorHandler());
app.use('*', securityHeaders());
app.use('*', analyticsTracker());
app.use('*', cacheControl());
app.use('*', rateLimit(100, 60));

// コンテンツ配信ルート
app.get('/posts/:slug{[a-z0-9\\-]+}/', async (c) => {
  const slug = c.req.param('slug');
  const accept = c.req.header('Accept');
  const result = await serveContent(
    slug,
    accept ?? null,
    c.get('repository'),
    c.get('cache'),
  );

  return new Response(result.body, {
    status: result.status,
    headers: {
      'Content-Type': result.contentType,
      ...result.headers,
    },
  });
});

// フィードルート
app.get('/feed/atom.xml', async (c) => {
  const atom = await c.env.CONTENT.get('feed/atom.xml');
  if (!atom) return c.notFound();
  return new Response(await atom.text(), {
    headers: { 'Content-Type': 'application/atom+xml; charset=utf-8' },
  });
});

// APIルート: 記事メタデータ
app.get('/api/posts', async (c) => {
  const db = c.get('db');
  const posts = await db.getAllPosts();
  return c.json(posts);
});

// APIルート: 関連記事
app.get('/api/posts/:slug/related', async (c) => {
  const slug = c.req.param('slug');
  const db = c.get('db');
  const related = await db.getRelatedPosts(slug);
  return c.json(related);
});

// APIルート: AI要約
app.get('/api/posts/:slug/summary', async (c) => {
  const slug = c.req.param('slug');
  const lang = (c.req.query('lang') ?? 'ja') as 'ja' | 'en' | 'zh';
  const cache = c.get('cache');

  const cacheKey = `ai:summary:${lang}:${slug}`;
  const cached = await cache.get(cacheKey);
  if (cached) return c.json(JSON.parse(cached));

  const repo = c.get('repository');
  const gemtext = await repo.getGemtext(slug);
  if (!gemtext) return c.notFound();

  const ai = c.get('ai');
  const summary = await ai.generateSummary(gemtext, lang);
  await cache.set(cacheKey, JSON.stringify(summary), 86400);
  return c.json(summary);
});

// トップページ
app.get('/', async (c) => {
  const index = await c.env.CONTENT.get('index.html');
  if (!index) return c.notFound();
  return new Response(await index.text(), {
    headers: { 'Content-Type': 'text/html; charset=utf-8' },
  });
});

// 404ハンドラー
app.notFound((c) => {
  return new Response(
    '<h1>404 Not Found</h1><p>The requested Gemini content was not found.</p>',
    {
      status: 404,
      headers: {
        'Content-Type': 'text/html; charset=utf-8',
        'X-Gemini-Status': '51',
        'X-Gemini-Meta': 'Not found',
      },
    },
  );
});

export default app;
