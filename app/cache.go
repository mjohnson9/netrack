package app

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CacheType string

const (
	CachePublic          CacheType = "public"
	CachePrivate                   = "private"
	CacheNoCache                   = "no-cache"
	CacheNoStore                   = "no-store"
	CacheMustRevalidate            = "must-revalidate"
	CacheProxyRevalidate           = "proxy-revalidate"
)

func SetCacheControl(w http.ResponseWriter, cacheType CacheType, maxAge time.Duration) {
	controlString := string(cacheType)

	if maxAge >= 0 {
		controlString += ", max-age=" + strconv.FormatInt(int64(maxAge/time.Second), 10)
	}

	headers := w.Header()

	headers.Set("Cache-Control", controlString)
	headers.Set("Pragma", strings.Title(string(cacheType)))
}
