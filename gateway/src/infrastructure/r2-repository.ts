// src/infrastructure/r2-repository.ts
import type { ContentRepository, PostMeta } from '../application/serve-content';

export class R2ContentRepository implements ContentRepository {
  constructor(private readonly bucket: R2Bucket) {}

  async getHtml(slug: string): Promise<string | null> {
    const object = await this.bucket.get(`posts/${slug}/index.html`);
    return object ? await object.text() : null;
  }

  async getGemtext(slug: string): Promise<string | null> {
    const object = await this.bucket.get(`posts/${slug}/index.gmi`);
    return object ? await object.text() : null;
  }

  async getMeta(slug: string): Promise<PostMeta | null> {
    const object = await this.bucket.get(`posts/${slug}/meta.json`);
    if (!object) return null;
    return await object.json<PostMeta>();
  }
}
