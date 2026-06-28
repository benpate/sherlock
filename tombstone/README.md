# tombstone

A [Sherlock](../README.md) client middleware that turns a "Gone" (HTTP 410) response into a synthetic ActivityStreams Tombstone document, so deleted remote objects resolve to a stable placeholder instead of an error.

## What matters here

- **Only a `derp.IsGone` error is converted; every other error passes through.** A 410 becomes a Tombstone returned WITHOUT an error; all other failures are returned unchanged. Don't broaden this to other status codes — a Tombstone is a claim that the object is *permanently gone*, not merely unreachable.

- **Returning the Tombstone with a nil error is intentional — it lets the result be cached.** The whole point is to stop hammering a server that keeps answering 410. A success-shaped return is what places the placeholder in the cache.

- **`New` wires itself as the inner client's root client.** It calls `innerClient.SetRootClient(self)` at construction; a parent may later override this. Keep that call — some hannibal clients make recursive calls back through the root.

- **An existing Tombstone is not overwritten.** If the inner result is already a Tombstone, it is returned as-is rather than rebuilt.
