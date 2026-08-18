/*
Package sherlock inspects a URL for any and all available metadata, using as
many methods as possible: ActivityStreams, Open Graph, Twitter Cards, oEmbed,
Microformats2, RSS/Atom/JSON feeds, and native HTML signals.

It has two first-class APIs built on shared machinery: Load returns an
ActivityStreams document (native AS2 when the server offers it, synthesized
from everything else when it doesn't), and Metadata returns an oEmbed-adjacent
metadata.Preview built by the metadata sub-package's extraction engine.
*/
package sherlock
