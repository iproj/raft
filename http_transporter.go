package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Parts from this transporter were heavily influenced by Peter Bougon's
// raft implementation: https://github.com/peterbourgon/raft

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// An HTTPTransporter is a default transport layer used to communicate between
// multiple servers.
type HTTPTransporter struct {
	DisableKeepAlives    bool
	prefix               string
	appendEntriesPath    string
	requestVotePath      string
	snapshotPath         string
	snapshotRecoveryPath string
	redirectPath         string
	peerJoinPath         string
	peerRemovePath       string
	httpClient           http.Client
	Transport            *http.Transport
}

type HTTPMuxer interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new HTTP transporter with the given path prefix.
func NewHTTPTransporter(prefix string, timeout time.Duration) *HTTPTransporter {
	t := &HTTPTransporter{
		DisableKeepAlives:    false,
		prefix:               prefix,
		appendEntriesPath:    joinPath(prefix, "/appendEntries"),
		requestVotePath:      joinPath(prefix, "/requestVote"),
		snapshotPath:         joinPath(prefix, "/snapshot"),
		snapshotRecoveryPath: joinPath(prefix, "/snapshotRecovery"),
		redirectPath:         joinPath(prefix, "/redirect"),
		peerJoinPath:         joinPath(prefix, "/join"),
		peerRemovePath:       joinPath(prefix, "/remove"),
		Transport:            &http.Transport{DisableKeepAlives: false},
	}
	t.httpClient.Transport = t.Transport
	t.Transport.ResponseHeaderTimeout = timeout
	return t
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the path prefix used by the transporter.
func (t *HTTPTransporter) Prefix() string {
	return t.prefix
}

func (t *HTTPTransporter) RedirectPath() string {
	return t.redirectPath
}

func (t *HTTPTransporter) PeerJoinPath() string {
	return t.peerJoinPath
}

// Retrieves the AppendEntries path.
func (t *HTTPTransporter) AppendEntriesPath() string {
	return t.appendEntriesPath
}

// Retrieves the RequestVote path.
func (t *HTTPTransporter) RequestVotePath() string {
	return t.requestVotePath
}

// Retrieves the Snapshot path.
func (t *HTTPTransporter) SnapshotPath() string {
	return t.snapshotPath
}

// Retrieves the SnapshotRecovery path.
func (t *HTTPTransporter) SnapshotRecoveryPath() string {
	return t.snapshotRecoveryPath
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Installation
//--------------------------------------

// Applies Raft routes to an HTTP router for a given server.
func (t *HTTPTransporter) Install(server Server, mux HTTPMuxer) {
	mux.HandleFunc(t.AppendEntriesPath(), t.appendEntriesHandler(server))
	mux.HandleFunc(t.RequestVotePath(), t.requestVoteHandler(server))
	mux.HandleFunc(t.SnapshotPath(), t.snapshotHandler(server))
	mux.HandleFunc(t.SnapshotRecoveryPath(), t.snapshotRecoveryHandler(server))
	mux.HandleFunc(t.peerJoinPath, t.peerJoinHandler(server))
	mux.HandleFunc(t.peerRemovePath, t.peerRemoveHandler(server))
}

//--------------------------------------
// Outgoing
//--------------------------------------

// Sends an AppendEntries RPC to a peer.
func (t *HTTPTransporter) SendAppendEntriesRequest(server Server, peer *Peer, req *AppendEntriesRequest) *AppendEntriesResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		traceln("transporter.ae.encoding.error:", err)
		return nil
	}

	url := joinPath(peer.ConnectionString, t.AppendEntriesPath())
	traceln(server.Name(), "POST", url)

	httpResp, err := t.httpClient.Post(url, "application/protobuf", &b)
	if httpResp == nil || err != nil {
		traceln("transporter.ae.response.error:", err)
		return nil
	}
	defer httpResp.Body.Close()

	resp := &AppendEntriesResponse{}
	if _, err = resp.Decode(httpResp.Body); err != nil && err != io.EOF {
		traceln("transporter.ae.decoding.error:", err)
		return nil
	}

	return resp
}

func (t *HTTPTransporter) Redirect(server Server, command Command) error {

	bytez, _ := json.Marshal(command)
	peer, ok := server.Peers()[server.Leader()]
	if !ok {
		return fmt.Errorf("Leader: %s has not connectAddr", server.Leader())
	}
	url := fmt.Sprintf("%s%s", peer.ConnectionString, t.redirectPath)
	httpResp, err := http.Post(url, "application/json", bytes.NewReader(bytez))
	if err != nil {
		return fmt.Errorf("Post %s failed: %v", url, err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Invalid http code: %d", httpResp.StatusCode)
	}
	return nil
}

// Sends a RequestVote RPC to a peer.
func (t *HTTPTransporter) SendVoteRequest(server Server, peer *Peer, req *RequestVoteRequest) *RequestVoteResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		traceln("transporter.rv.encoding.error:", err)
		return nil
	}

	url := fmt.Sprintf("%s%s", peer.ConnectionString, t.RequestVotePath())
	traceln(server.Name(), "POST", url)

	httpResp, err := t.httpClient.Post(url, "application/protobuf", &b)
	if httpResp == nil || err != nil {
		traceln("transporter.rv.response.error:", err)
		return nil
	}
	defer httpResp.Body.Close()

	resp := &RequestVoteResponse{}
	if _, err = resp.Decode(httpResp.Body); err != nil && err != io.EOF {
		traceln("transporter.rv.decoding.error:", err)
		return nil
	}

	return resp
}

func joinPath(connectionString, thePath string) string {
	u, err := url.Parse(connectionString)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, thePath)
	return u.String()
}

// Sends a SnapshotRequest RPC to a peer.
func (t *HTTPTransporter) SendSnapshotRequest(server Server, peer *Peer, req *SnapshotRequest) *SnapshotResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		traceln("transporter.rv.encoding.error:", err)
		return nil
	}

	url := joinPath(peer.ConnectionString, t.snapshotPath)
	traceln(server.Name(), "POST", url)

	httpResp, err := t.httpClient.Post(url, "application/protobuf", &b)
	if httpResp == nil || err != nil {
		traceln("transporter.rv.response.error:", err)
		return nil
	}
	defer httpResp.Body.Close()

	resp := &SnapshotResponse{}
	if _, err = resp.Decode(httpResp.Body); err != nil && err != io.EOF {
		traceln("transporter.rv.decoding.error:", err)
		return nil
	}

	return resp
}

// Sends a SnapshotRequest RPC to a peer.
func (t *HTTPTransporter) SendSnapshotRecoveryRequest(server Server, peer *Peer, req *SnapshotRecoveryRequest) *SnapshotRecoveryResponse {
	var b bytes.Buffer
	if _, err := req.Encode(&b); err != nil {
		traceln("transporter.rv.encoding.error:", err)
		return nil
	}

	url := joinPath(peer.ConnectionString, t.snapshotRecoveryPath)
	traceln(server.Name(), "POST", url)

	httpResp, err := t.httpClient.Post(url, "application/protobuf", &b)
	if httpResp == nil || err != nil {
		traceln("transporter.rv.response.error:", err)
		return nil
	}
	defer httpResp.Body.Close()

	resp := &SnapshotRecoveryResponse{}
	if _, err = resp.Decode(httpResp.Body); err != nil && err != io.EOF {
		traceln("transporter.rv.decoding.error:", err)
		return nil
	}

	return resp
}

//--------------------------------------
// Incoming
//--------------------------------------

func (t *HTTPTransporter) peerRemoveHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		debugln(server.Name(), "RECV /remove")
		command := &DefaultLeaveCommand{}
		if err := json.NewDecoder(r.Body).Decode(command); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, ok := server.Peers()[command.Name]
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid peer: %s", command.Name), http.StatusBadRequest)
			return
		}
		if _, err := server.Do(command); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (t *HTTPTransporter) peerJoinHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		debugln(server.Name(), "RECV /join")
		command := &DefaultJoinCommand{}
		if err := json.NewDecoder(r.Body).Decode(command); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, ok := server.Peers()[command.Name]
		if ok {
			http.Error(w, fmt.Sprintf("Already exist: %s", command.Name), http.StatusAlreadyReported)
			return
		}
		if len(server.Peers()) >= server.MaxPeerCount() {
			http.Error(w, "Can't be joined", http.StatusNotAcceptable)
			return
		}
		if _, err := server.Do(command); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Handles incoming AppendEntries requests.
func (t *HTTPTransporter) appendEntriesHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceln(server.Name(), "RECV /appendEntries")

		req := &AppendEntriesRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.AppendEntries(req)
		if resp == nil {
			http.Error(w, "Failed creating response.", http.StatusInternalServerError)
			return
		}
		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

// Handles incoming RequestVote requests.
func (t *HTTPTransporter) requestVoteHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceln(server.Name(), "RECV /requestVote")

		req := &RequestVoteRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.RequestVote(req)
		if resp == nil {
			http.Error(w, "Failed creating response.", http.StatusInternalServerError)
			return
		}
		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

// Handles incoming Snapshot requests.
func (t *HTTPTransporter) snapshotHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceln(server.Name(), "RECV /snapshot")

		req := &SnapshotRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.RequestSnapshot(req)
		if resp == nil {
			http.Error(w, "Failed creating response.", http.StatusInternalServerError)
			return
		}
		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

// Handles incoming SnapshotRecovery requests.
func (t *HTTPTransporter) snapshotRecoveryHandler(server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceln(server.Name(), "RECV /snapshotRecovery")

		req := &SnapshotRecoveryRequest{}
		if _, err := req.Decode(r.Body); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := server.SnapshotRecoveryRequest(req)
		if resp == nil {
			http.Error(w, "Failed creating response.", http.StatusInternalServerError)
			return
		}
		if _, err := resp.Encode(w); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
