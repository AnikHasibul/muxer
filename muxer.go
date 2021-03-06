package muxer

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

type route struct {
	r  map[string]Handler
	mu sync.Mutex
}
type CTX struct {
	Ref     string
	Sess    *sessions.Session
	CSRF    *sessions.Session
	ReqID   string
	buff    string
	Err     error
	StartAt time.Time
	W       http.ResponseWriter
	R       *http.Request
}
type Handler func(*CTX)

var routes = route{
	r: make(map[string]Handler),
}
var (
	Boot  = func(*CTX) {}
	Defer = func(*CTX) {}
)

func Root(w http.ResponseWriter, r *http.Request) {
	if f, ok := routes.r[r.URL.Path]; ok {
		c := &CTX{
			Ref:     r.URL.Path,
			ReqID:   Key(),
			StartAt: time.Now(),
			W:       w,
			R:       r,
		}
		Boot(c)
		if c.Err != nil {
			return
		}
		defer Defer(c)
		f(c)
	}
}

func Register(path string, h Handler) {
	routes.mu.Lock()
	routes.r[path] = h
	routes.mu.Unlock()
}

func Delete(path string) {
	routes.mu.Lock()
	delete(routes.r, path)
	routes.mu.Unlock()
}

func Key() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func (c *CTX) BuffAdd(buff string) {
	c.buff += buff
}
func (c *CTX) BuffSet(buff string) {
	c.buff = buff
}
func (c *CTX) BuffGet() string {
	return c.buff
}
func (c *CTX) BuffFlush(w io.Writer) {
	if c.buff != "" {
		fmt.Fprintln(w, c.buff)
	}
}
