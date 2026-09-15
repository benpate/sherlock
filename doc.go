/*
Package sherlock inspects a URL for any and all available metadata, using as
many methods as possible: ActivityStreams, Open Graph, Twitter Cards, oEmbed,
and native HTML signals.

The work happens in the sub-packages. Webpage metadata extraction is the
metadata sub-package, whose Get function fetches a URL once and returns an
oEmbed-adjacent Preview. ActivityStreams loading is a stack of client
middlewares -- activitypub, tombstone, webfinger, bridgyfed, and tagspub --
each implementing hannibal's streams.Client interface.

This root package holds only what those sub-packages and their callers share:
AuthorizedFetch, which signs an outbound request according to the ActivityPub
Authorized Fetch convention, and the address classifiers that decide whether a
value looks like something Sherlock could resolve.
*/
package sherlock
