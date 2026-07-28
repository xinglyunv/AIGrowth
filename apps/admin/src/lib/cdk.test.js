import assert from 'node:assert/strict';
import test from 'node:test';
import { toRFC3339 } from './cdk.ts';

test('converts datetime-local input to an RFC3339 timestamp', () => {
  assert.equal(toRFC3339('2026-07-31T00:28'), '2026-07-31T00:28:00.000Z');
});

test('keeps empty expiry values empty', () => {
  assert.equal(toRFC3339(''), null);
});
