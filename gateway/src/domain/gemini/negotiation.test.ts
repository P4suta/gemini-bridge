// src/domain/gemini/negotiation.test.ts
import { describe, it, expect } from 'vitest';
import { negotiateFormat } from './negotiation';

describe('negotiateFormat', () => {
  it('returns html for browser Accept header', () => {
    const accept = 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8';
    expect(negotiateFormat(accept)).toBe('html');
  });

  it('returns gemtext for Gemini client Accept header', () => {
    expect(negotiateFormat('text/gemini')).toBe('gemtext');
  });

  it('returns json for JSON Accept header', () => {
    expect(negotiateFormat('application/json')).toBe('json');
  });

  it('returns html for null Accept header', () => {
    expect(negotiateFormat(null)).toBe('html');
  });

  it('returns html for wildcard Accept header', () => {
    expect(negotiateFormat('*/*')).toBe('html');
  });

  it('respects quality values', () => {
    const accept = 'text/gemini;q=0.9, text/html;q=1.0';
    expect(negotiateFormat(accept)).toBe('html');
  });
});
