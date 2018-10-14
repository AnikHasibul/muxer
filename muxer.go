package muxer

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type route struct {
	r  map[string]Handler
	mu sync.Mutex
}
type CTX struct {
	ReqID   string
	buff    string
	StartAt time.Time
	W       http.ResponseWriter
	R       *http.Request
}
type Handler func(*CTX)

var routes = route{
	r: make(map[string]Handler),
}

func Root(w http.ResponseWriter, r *http.Request) {
	if f, ok := routes.r[r.URL.Path]; ok {
		c := &CTX{
			ReqID:   Key(),
			StartAt: time.Now(),
			W:       w,
			R:       r,
		}
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
	fmt.Fprintln(w, c.buff)
}
func (c *CTX) BuffFlushByte(w io.Writer) {
	w.Write([]byte(c.buff))
}
