import { describe, it, expect } from 'vitest';
import { statusCategory } from './types';
import type { GeminiStatusCode } from './types';

describe('statusCategory', () => {
  it('returns INPUT for 10', () => {
    expect(statusCategory(10)).toBe('INPUT');
  });

  it('returns INPUT for 11', () => {
    expect(statusCategory(11)).toBe('INPUT');
  });

  it('returns SUCCESS for 20', () => {
    expect(statusCategory(20)).toBe('SUCCESS');
  });

  it('returns REDIRECT for 30', () => {
    expect(statusCategory(30)).toBe('REDIRECT');
  });

  it('returns REDIRECT for 31', () => {
    expect(statusCategory(31)).toBe('REDIRECT');
  });

  it('returns TEMPORARY_FAILURE for 40-44', () => {
    const codes: GeminiStatusCode[] = [40, 41, 42, 43, 44];
    for (const code of codes) {
      expect(statusCategory(code)).toBe('TEMPORARY_FAILURE');
    }
  });

  it('returns PERMANENT_FAILURE for 50-59', () => {
    const codes: GeminiStatusCode[] = [50, 51, 52, 53, 59];
    for (const code of codes) {
      expect(statusCategory(code)).toBe('PERMANENT_FAILURE');
    }
  });

  it('returns CLIENT_CERTIFICATE for 60-62', () => {
    const codes: GeminiStatusCode[] = [60, 61, 62];
    for (const code of codes) {
      expect(statusCategory(code)).toBe('CLIENT_CERTIFICATE');
    }
  });
});
