// src/infrastructure/workers-ai.ts

export interface AISummaryResult {
  summary: string;
  language: string;
  model: string;
}

export class WorkersAIClient {
  constructor(private readonly ai: Ai) {}

  /**
   * gemtextコンテンツから指定言語の要約を生成する。
   */
  async generateSummary(
    gemtext: string,
    targetLanguage: 'ja' | 'en' | 'zh',
  ): Promise<AISummaryResult> {
    const languageNames: Record<string, string> = {
      ja: 'Japanese',
      en: 'English',
      zh: 'Chinese',
    };

    const response = await this.ai.run('@cf/meta/llama-3.1-8b-instruct', {
      messages: [
        {
          role: 'system',
          content: `You are a technical blog summarizer. Generate a concise summary (2-3 sentences) in ${languageNames[targetLanguage]}. The input is in Gemini protocol's gemtext format. Focus on the key technical concepts and takeaways.`,
        },
        {
          role: 'user',
          content: gemtext,
        },
      ],
      max_tokens: 256,
    });

    return {
      summary: (response as { response: string }).response,
      language: targetLanguage,
      model: '@cf/meta/llama-3.1-8b-instruct',
    };
  }

  /**
   * OGP description用の短い要約を生成する。
   */
  async generateOGPDescription(gemtext: string): Promise<string> {
    const response = await this.ai.run('@cf/meta/llama-3.1-8b-instruct', {
      messages: [
        {
          role: 'system',
          content:
            'Generate a concise meta description (max 160 characters) in the same language as the input. Focus on what the article is about. Do not include quotes or formatting.',
        },
        {
          role: 'user',
          content: gemtext,
        },
      ],
      max_tokens: 64,
    });

    return (response as { response: string }).response;
  }
}
