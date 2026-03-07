// src/infrastructure/d1-repository.ts

export interface AnalyticsEvent {
  path: string;
  userAgent: string;
  country: string;
  referer: string;
}

export class D1Repository {
  constructor(private readonly db: D1Database) {}

  async getPostBySlug(slug: string): Promise<PostRecord | null> {
    const result = await this.db
      .prepare('SELECT * FROM posts WHERE slug = ?')
      .bind(slug)
      .first<PostRecord>();
    return result;
  }

  async getAllPosts(): Promise<PostRecord[]> {
    const result = await this.db
      .prepare('SELECT * FROM posts ORDER BY published_at DESC')
      .all<PostRecord>();
    return result.results;
  }

  async recordAnalytics(event: AnalyticsEvent): Promise<void> {
    await this.db
      .prepare(
        'INSERT INTO analytics (path, timestamp, user_agent, country, referer) VALUES (?, ?, ?, ?, ?)'
      )
      .bind(
        event.path,
        new Date().toISOString(),
        event.userAgent,
        event.country,
        event.referer,
      )
      .run();
  }

  async getRelatedPosts(slug: string, limit: number = 5): Promise<PostRecord[]> {
    const current = await this.getPostBySlug(slug);
    if (!current || !current.tags) return [];

    // 現在の記事のタグをパース
    let currentTags: string[];
    try {
      currentTags = JSON.parse(current.tags);
    } catch {
      return [];
    }

    if (currentTags.length === 0) return [];

    // json_each() でタグJSON配列を行に展開し、
    // 共通タグ数の多い順に関連記事を取得する
    const placeholders = currentTags.map(() => '?').join(', ');
    const result = await this.db
      .prepare(
        `SELECT p.*, COUNT(DISTINCT jt.value) AS shared_tags
         FROM posts p, json_each(p.tags) jt
         WHERE p.slug != ?
           AND jt.value IN (${placeholders})
         GROUP BY p.slug
         ORDER BY shared_tags DESC, p.published_at DESC
         LIMIT ?`
      )
      .bind(slug, ...currentTags, limit)
      .all<PostRecord>();
    return result.results;
  }
}

interface PostRecord {
  id: string;
  slug: string;
  title: string;
  published_at: string;
  updated_at: string | null;
  summary: string | null;
  word_count: number;
  language: string;
  tags: string;
  gemtext_hash: string;
}
