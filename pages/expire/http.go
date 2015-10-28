package pages

import (
	"io"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/delay"
	"github.com/gaego/ds"

	"app"
)

func init() {
	app.Get("/_cron/expire", app.CreateHandler(expire))
}

func getExpiredQuery(expireTime time.Time) *datastore.Query {
	return datastore.NewQuery("Peer").Filter("Expires <=", expireTime)
}

func doExpire(c appengine.Context, expireTime time.Time, cursorString string) {
	query := getExpiredQuery(expireTime).KeysOnly()

	if len(cursorString) > 0 {
		if cursor, err := datastore.DecodeCursor(cursorString); err != nil {
			c.Errorf("Failed to decode cursor: %s", err)
			return
		} else {
			query = query.Start(cursor)
		}
	}

	for {
		toDelete := make([]*datastore.Key, 0, 100)

		for queryIterator := query.Run(c); ; {
			peerKey, err := queryIterator.Next(nil)
			if len(toDelete) >= 100 {
				break
			} else if err == datastore.Done {
				c.Infof("Done finding expired peers")
				break
			} else if err != nil {
				c.Criticalf("Failed to get next peer: %#v (%s)", err, err)
				return
			}

			toDelete = append(toDelete, peerKey)
		}

		deleted := len(toDelete)
		c.Infof("Deleting %d expired peers", deleted)
		if err := ds.DeleteMulti(c, toDelete); err != nil {
			c.Criticalf("Failed to delete peers: %#v (%s)", err, err)
			return
		}

		if deleted <= 0 {
			break
		}
	}

	c.Infof("Finished deleting expired peers")
}

var expireDelay = delay.Func("doExpire", doExpire)

func expire(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	expireDelay.Call(c, time.Now(), "")

	w.WriteHeader(200)
	io.WriteString(w, "Expiring peers...")

	c.Debugf("Successfully began expire")
}
