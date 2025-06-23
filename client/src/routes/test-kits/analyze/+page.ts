import { requireAuth } from "$lib/api";

export function load({ url }) {
  return requireAuth(url);
}
