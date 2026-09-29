package activitypub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/stretchr/testify/require"
)

// These tests serve one httptest server under two names.  The cache and this package treat
// "127.0.0.1" and "localhost" as different hosts, so one server can play both victim and attacker.

// TestTrusted_SameHostIsAccepted confirms that a document whose id is on the host that served it
// loads with a single request.
func TestTrusted_SameHostIsAccepted(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/actor"), "PEM-GENUINE"))

	result, err := server.client().Load(server.victimURL("/actor"))

	require.NoError(t, err)
	require.Equal(t, "PEM-GENUINE", result.PublicKey().PublicKeyPEM())
	require.Equal(t, 1, server.requests("/actor"))
}

// TestTrusted_NoIDIsAccepted confirms that a document with no id, which claims nothing, loads.
func TestTrusted_NoIDIsAccepted(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/anonymous", map[string]any{vocab.PropertyType: vocab.ObjectTypeNote})

	result, err := server.client().Load(server.victimURL("/anonymous"))

	require.NoError(t, err)
	require.Equal(t, vocab.ObjectTypeNote, result.Type())
}

// TestTrusted_ForeignCopyIsReplacedByTheOriginal confirms that a document claiming another host's id
// is exchanged for the copy that host serves itself.
func TestTrusted_ForeignCopyIsReplacedByTheOriginal(t *testing.T) {

	// BUG-223: the forged copy would be cached under the victim's id, key and all.
	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/actor"), "PEM-GENUINE"))
	server.serve("/forged", server.document(server.victimURL("/actor"), "PEM-ATTACKER"))

	result, err := server.client().Load(server.attackerURL("/forged"))

	require.NoError(t, err)
	require.Equal(t, server.victimURL("/actor"), result.ID())
	require.Equal(t, "PEM-GENUINE", result.PublicKey().PublicKeyPEM())
	require.Equal(t, 1, server.requests("/actor"))
}

// TestTrusted_DeniedByItsOwnHost confirms that a document is rejected when the host its id names
// serves a different id there, and that the rejection never reaches the inner client.
func TestTrusted_DeniedByItsOwnHost(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/someone-else"), "PEM-GENUINE"))
	server.serve("/forged", server.document(server.victimURL("/actor"), "PEM-ATTACKER"))

	inner := &fakeClient{loadResult: streams.NewDocument(map[string]any{vocab.PropertyID: "inner"})}
	client := New(WithAllowPrivateIPs(true), WithInnerClient(inner))

	result, err := client.Load(server.attackerURL("/forged"))

	requireForbidden(t, err)
	require.True(t, result.IsNil())
	require.Empty(t, inner.loadedID, "a rejected document must not fall through to the inner client")
}

// TestTrusted_MissingAtItsOwnHost confirms that a document is rejected when the host its id names
// answers 404, and that the rejection is Forbidden rather than NotFound.
func TestTrusted_MissingAtItsOwnHost(t *testing.T) {

	// A NotFound here would let a forger make a cache delete whatever URL it was loading
	server := newOriginServer(t)
	server.serve("/forged", server.document(server.victimURL("/nobody"), "PEM-ATTACKER"))

	_, err := server.client().Load(server.attackerURL("/forged"))

	requireForbidden(t, err)
	require.False(t, derp.IsNotFound(err))
}

// TestTrusted_OwnHostServesHTML confirms that a confirmation which is not an ActivityPub document is
// a rejection.
func TestTrusted_OwnHostServesHTML(t *testing.T) {

	server := newOriginServer(t)
	server.serveHTML("/page")
	server.serve("/forged", server.document(server.victimURL("/page"), "PEM-ATTACKER"))

	_, err := server.client().Load(server.attackerURL("/forged"))

	requireForbidden(t, err)
}

// TestTrusted_OwnHostServesOtherJSON confirms that a confirmation must declare an ActivityPub content
// type, even when its body is JSON naming the right id.
func TestTrusted_OwnHostServesOtherJSON(t *testing.T) {

	server := newOriginServer(t)
	server.serveAs("/actor", "application/feed+json", server.document(server.victimURL("/actor"), "PEM-GENUINE"))
	server.serve("/forged", server.document(server.victimURL("/actor"), "PEM-ATTACKER"))

	_, err := server.client().Load(server.attackerURL("/forged"))

	requireForbidden(t, err)
}

// TestTrusted_ConfirmationIsNeverChased confirms that a confirmation naming yet another id is rejected
// without a third request.
func TestTrusted_ConfirmationIsNeverChased(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/forged", server.document(server.victimURL("/actor"), "PEM-ATTACKER"))
	server.serve("/actor", server.document(server.attackerURL("/elsewhere"), "PEM-OTHER"))
	server.serve("/elsewhere", server.document(server.attackerURL("/elsewhere"), "PEM-OTHER"))

	_, err := server.client().Load(server.attackerURL("/forged"))

	requireForbidden(t, err)
	require.Equal(t, 1, server.requests("/actor"))
	require.Zero(t, server.requests("/elsewhere"))
}

// TestTrusted_RedirectToTheIDsHostIsAccepted confirms that a redirect ending on the id's own host
// needs no confirmation, because the final URL is what vouches.
func TestTrusted_RedirectToTheIDsHostIsAccepted(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/actor"), "PEM-GENUINE"))
	server.redirect("/moved", server.victimURL("/actor"))

	result, err := server.client().Load(server.attackerURL("/moved"))

	require.NoError(t, err)
	require.Equal(t, "PEM-GENUINE", result.PublicKey().PublicKeyPEM())
	require.Equal(t, 1, server.requests("/actor"))
}

// TestTrusted_OpenRedirectCannotVouch confirms that an open redirect on the victim's host, landing on
// the attacker, cannot make the attacker's document speak for an id on the victim's host.
func TestTrusted_OpenRedirectCannotVouch(t *testing.T) {

	// The forged document claims the redirect URL itself as its id, so the confirmation fetch
	// follows the same redirect back to the attacker.  Only the final URL tells the two apart.
	server := newOriginServer(t)
	server.redirect("/go", server.attackerURL("/forged"))
	server.serve("/forged", server.document(server.victimURL("/go"), "PEM-ATTACKER"))

	_, err := server.client().Load(server.victimURL("/go"))

	requireForbidden(t, err)
}

// TestTrusted_UnusableIDs confirms that an id which is not a string, not a URL, or not what
// Document.ID reports, is rejected even when it names the serving host.
func TestTrusted_UnusableIDs(t *testing.T) {

	server := newOriginServer(t)

	test := func(name string, id any) {
		t.Run(name, func(t *testing.T) {
			server.serve("/"+name, map[string]any{vocab.PropertyID: id, vocab.PropertyType: vocab.ObjectTypeNote})
			_, err := server.client().Load(server.victimURL("/" + name))
			requireForbidden(t, err)
		})
	}

	test("number", 42)
	test("list", []any{server.victimURL("/list")})
	test("object", map[string]any{"href": server.victimURL("/object")})
	test("relative", "/relative")
	test("scheme", "ftp://127.0.0.1/scheme")
	test("comment", server.victimURL("/comment<!-- -->"))
	test("markup", server.victimURL("/markup<b>x</b>"))
}

// TestTrusted_AlternateIDIsChecked confirms that a document using "@id" in place of "id" cannot
// skip the check.
func TestTrusted_AlternateIDIsChecked(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/actor"), "PEM-GENUINE"))
	server.serve("/forged", map[string]any{vocab.PropertyID_Alternate: server.victimURL("/actor"), vocab.PropertyType: vocab.ActorTypePerson})

	result, err := server.client().Load(server.attackerURL("/forged"))

	// Either outcome is safe: the victim's own copy, or a rejection
	if err != nil {
		requireForbidden(t, err)
		return
	}

	require.Equal(t, "PEM-GENUINE", result.PublicKey().PublicKeyPEM())
}

// TestTrusted_ConfirmationCarriesCallerOptions confirms that the confirming fetch is sent with the
// same caller options as the first.
func TestTrusted_ConfirmationCarriesCallerOptions(t *testing.T) {

	server := newOriginServer(t)
	server.serve("/actor", server.document(server.victimURL("/actor"), "PEM-GENUINE"))
	server.serve("/forged", server.document(server.victimURL("/actor"), "PEM-ATTACKER"))

	_, err := server.client().Load(server.attackerURL("/forged"), stampHeaderOption())

	require.NoError(t, err)
	require.Equal(t, "reached", server.header("/actor", "X-Sherlock-Test"))
}

// TestFinalURL_NoResponse confirms that a transaction with no response reports the requested URL.
func TestFinalURL_NoResponse(t *testing.T) {
	require.Equal(t, "https://example.com/a", finalURL(remote.Get("https://example.com/a"), "https://example.com/a"))
}

// TestSameHost confirms that hosts compare without case or port, and that a URL with no host vouches
// for nothing.
func TestSameHost(t *testing.T) {
	require.True(t, sameHost("https://example.com/a", "https://EXAMPLE.com:8443/b"))
	require.False(t, sameHost("https://example.com/a", "https://example.com.evil.com/a"))
	require.False(t, sameHost("https://user@example.com/a", "https://user.com/a"))
	require.False(t, sameHost("", ""))
	require.False(t, sameHost("https:///a", "https:///a"))
}

/******************************************
 * Helpers
 ******************************************/

// originServer is an httptest server that answers with canned documents, redirects, or HTML by path,
// and records every request it receives.
type originServer struct {
	server   *httptest.Server
	mutex    sync.Mutex
	handlers map[string]http.HandlerFunc
	counts   map[string]int
	headers  map[string]http.Header
}

// newOriginServer starts an originServer that closes when the test ends
func newOriginServer(t *testing.T) *originServer {

	t.Helper()

	result := &originServer{
		handlers: make(map[string]http.HandlerFunc),
		counts:   make(map[string]int),
		headers:  make(map[string]http.Header),
	}

	result.server = httptest.NewServer(http.HandlerFunc(result.handle))
	t.Cleanup(result.server.Close)

	return result
}

// client returns an activitypub Client that may reach this loopback server
func (server *originServer) client() streams.Client {
	return New(WithAllowPrivateIPs(true))
}

// victimURL returns a URL on this server addressed as 127.0.0.1
func (server *originServer) victimURL(path string) string {
	return server.server.URL + path
}

// attackerURL returns a URL on this server addressed as localhost, which is a different host
func (server *originServer) attackerURL(path string) string {
	return strings.Replace(server.server.URL, "127.0.0.1", "localhost", 1) + path
}

// document returns a minimal Person with the given id, publishing the given key
func (server *originServer) document(id string, publicKeyPEM string) map[string]any {
	return map[string]any{
		vocab.PropertyID:   id,
		vocab.PropertyType: vocab.ActorTypePerson,
		vocab.PropertyPublicKey: map[string]any{
			vocab.PropertyID:           id + "#main-key",
			vocab.PropertyPublicKeyPEM: publicKeyPEM,
		},
	}
}

// serve answers a path with a JSON document, as ActivityPub
func (server *originServer) serve(path string, document map[string]any) {
	server.serveAs(path, vocab.ContentTypeActivityPub, document)
}

// serveAs answers a path with a JSON document, under the given content type
func (server *originServer) serveAs(path string, contentType string, document map[string]any) {
	server.handleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_ = json.NewEncoder(w).Encode(document)
	})
}

// serveHTML answers a path with a web page
func (server *originServer) serveHTML(path string) {
	server.handleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body>Not ActivityPub</body></html>"))
	})
}

// redirect answers a path with a 302 to the target
func (server *originServer) redirect(path string, target string) {
	server.handleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusFound)
	})
}

// handleFunc registers a handler for a path
func (server *originServer) handleFunc(path string, handler http.HandlerFunc) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.handlers[path] = handler
}

// handle counts the request and passes it to the handler for its path, or answers 404
func (server *originServer) handle(w http.ResponseWriter, r *http.Request) {

	server.mutex.Lock()
	server.counts[r.URL.Path]++
	server.headers[r.URL.Path] = r.Header.Clone()
	handler, exists := server.handlers[r.URL.Path]
	server.mutex.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	handler(w, r)
}

// requests returns how many requests reached a path
func (server *originServer) requests(path string) int {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	return server.counts[path]
}

// header returns a header from the most recent request to a path
func (server *originServer) header(path string, name string) string {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	return server.headers[path].Get(name)
}

// stampHeaderOption returns a caller option that marks each request it reaches with a test header
func stampHeaderOption() remote.Option {
	return remote.Option{
		BeforeRequest: func(transaction *remote.Transaction) error {
			transaction.Header("X-Sherlock-Test", "reached")
			return nil
		},
	}
}

// requireForbidden fails the test unless err is a Forbidden error
func requireForbidden(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	require.Equal(t, http.StatusForbidden, derp.ErrorCode(err), "want Forbidden, got %v", err)
}
