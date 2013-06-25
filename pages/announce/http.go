package index

import (
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"github.com/zeebo/bencode"

	"app"
	"models"
)

const (
	eventTypeStarted = iota
	eventTypeStopped = iota
)

var internalError string

func init() {
	var err error

	internalError, err = bencode.EncodeString(&errorResponse{
		MinimumInterval: 60,
		Interval:        300,
		Reason:          "Internal error. Please try again later.",
	})
	if err != nil {
		panic("Error setting up internalError: " + err.Error())
	}

	app.Get("/announce", app.CreateHandler(run))
}

type errorResponse struct {
	MinimumInterval uint64 `bencode:"min interval"`
	Interval        uint64 `bencode:"interval"`
	Reason          string `bencode:"failure reason"`
}

type peerResponse struct {
	MinimumInterval uint64         `bencode:"min interval"`
	Interval        uint64         `bencode:"interval"`
	Peers           []*models.Peer `bencode:"peers"`
}

func writeInternalError(w http.ResponseWriter) {
	w.WriteHeader(500)
	io.WriteString(w, internalError)
}

func writeEncoded(c appengine.Context, w http.ResponseWriter, data interface{}) {
	newEncoder := bencode.NewEncoder(w)
	if err := newEncoder.Encode(data); err != nil {
		c.Errorf("Failed to encode data: %s", err)
		writeInternalError(w)
	}
}

func run(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	q := r.URL.Query()

	rawHash := q.Get("info_hash")
	if len(rawHash) == 0 {
		c.Infof("Missing info hash")
		w.WriteHeader(400)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Missing info_hash",
		})
		return
	}

	hash := fmt.Sprintf("%x", rawHash)
	if len(hash) != 40 {
		c.Infof("Invalid info hash: %s", hash)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Invalid info_hash",
		})
		return
	}
	c.Debugf("Info hash: %s", hash)

	rawPeerID := q.Get("peer_id")
	if rawPeerIDLength := len(rawPeerID); rawPeerIDLength == 0 {
		c.Infof("Missing peer ID")
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Missing peer_id",
		})
		return
	} else if rawPeerIDLength != 20 {
		c.Infof("Invalid peer ID: %s", rawPeerID)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Invalid peer_id",
		})
		return
	}
	peerID := hex.EncodeToString([]byte(rawPeerID))
	c.Debugf("Peer ID: %s", peerID)

	/*rawIP := q.Get("ip")
	if len(rawIP) == 0 {
		rawIP = r.RemoteAddr
	}*/
	rawIP := r.RemoteAddr

	IP := net.ParseIP(rawIP)
	if IP == nil {
		c.Infof("Invalid IP: %s", rawIP)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Invalid IP",
		})
		return
	}
	c.Debugf("IP: %s", IP.String())

	rawPort := q.Get("port")
	if len(rawPort) == 0 {
		c.Infof("Missing port")
		w.WriteHeader(400)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Missing port",
		})
		return
	}

	port64, err := strconv.ParseUint(rawPort, 10, 16)
	if err != nil {
		c.Infof("Failed to parse port: %s", err)
		w.WriteHeader(400)
		writeEncoded(c, w, &errorResponse{
			MinimumInterval: 600,
			Interval:        3600,
			Reason:          "Invalid port",
		})
		return
	}
	port := uint16(port64)
	c.Debugf("Port: %d", port)

	eventRaw := strings.ToLower(q.Get("event"))
	event := eventTypeStarted
	if len(eventRaw) > 0 {
		if eventRaw == "started" || eventRaw == "completed" {
			event = eventTypeStarted
		} else if eventRaw == "stopped" {
			event = eventTypeStopped
		} else {
			c.Infof("Invalid event: %s", q.Get("event"))
			w.WriteHeader(400)
			writeEncoded(c, w, &errorResponse{
				MinimumInterval: 600,
				Interval:        3600,
				Reason:          "Invalid event",
			})
			return
		}
	}
	c.Debugf("Event type: %s (%d)", eventRaw, event)

	var peersWanted int64 = 50
	rawWanted := q.Get("numwant")
	if len(rawWanted) > 0 {
		if tempWanted, err := strconv.ParseInt(rawWanted, 10, 64); err == nil {
			peersWanted = int64(math.Max(0, math.Min(500, float64(tempWanted))))
		} else {
			c.Infof("Invalid number of peers wanted: %s", rawWanted)
			w.WriteHeader(400)
			writeEncoded(c, w, &errorResponse{
				MinimumInterval: 600,
				Interval:        3600,
				Reason:          "Invalid numwant",
			})
			return
		}
	}
	c.Debugf("Number of peers: %d", peersWanted)

	var eventFuture app.Future
	if event == eventTypeStarted {
		eventFuture = models.InsertPeer(c, hash, peerID, IP.String(), port)
	} else if event == eventTypeStopped {
		eventFuture = models.RemovePeer(c, hash, peerID)
	}

	// If we have an eventFuture, wait until it finishes before returning the request.
	if eventFuture != nil {
		defer eventFuture.Wait()
	}

	response := &peerResponse{
		MinimumInterval: uint64(models.PeerRefreshTime / time.Second),
		Interval:        uint64(models.PeerExpireTime / time.Second),
	}

	var peerFuture app.Future
	if peersWanted > 0 {
		peerFuture = models.GetPeers(c, hash, peersWanted, peerID)
	}

	if peerFuture != nil {
		rawPeers, err := peerFuture.GetResult()
		if err != nil {
			writeInternalError(w)
			return
		}

		if peers, ok := rawPeers.([]*models.Peer); ok {
			for _, peer := range peers {
				if idBytes, err := hex.DecodeString(peer.ID); err == nil {
					peer.ID = string(idBytes)
					response.Peers = append(response.Peers, peer)
				} else {
					c.Infof("Error decoding peer ID (hash: %s) (ID: %s): %s", peer.Hash, peer.ID, err)
				}
			}
		} else {
			c.Errorf("Result from GetPeers is invalid: %#v", rawPeers)
			writeInternalError(w)
			return
		}
	}

	w.WriteHeader(200)
	writeEncoded(c, w, response)
}
