package route

import (
	"net/http"
	"regexp"
)

type Result int

const (
	R_RETURN = Result(iota)
)

type Handle interface {
	ServeCtx(*SessionCtx) Result
}

type HandleFunc func(*SessionCtx) Result

func (hf HandleFunc) ServeCtx(sc *SessionCtx) Result {
	return hf(sc)
}

type muxEntry struct {
	reg     *regexp.Regexp
	handler Handle
}

func (me *muxEntry) match(tpl string) (Handle, bool) {
	if me.reg.MatchString(tpl) {
		return me.handler, true
	}
	return nil, false
}

type Mux struct {
	filter []*muxEntry
	handle []*muxEntry
	Values map[string]interface{}
}

func NewMux() *Mux {
	return &Mux{
		filter: []*muxEntry{},
		handle: []*muxEntry{},
	}
}

func (m *Mux) Hand(pattern string, handle Handle) {
	reg := regexp.MustCompile(pattern)
	m.handle = append(m.handle, &muxEntry{reg, handle})
}

func (m *Mux) HandFunc(pattern string, handleFunc HandleFunc) {
	m.Hand(pattern, handleFunc)
}

func (m *Mux) Filter(pattern string, handle Handle) {
	reg := regexp.MustCompile(pattern)
	m.filter = append(m.filter, &muxEntry{reg, handle})
}

func (m *Mux) FilterFunc(pattern string, handleFunc HandleFunc) {
	m.Filter(pattern, handleFunc)
}

func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := NewSessionCtx(rw, req)
	ctx.Mux = m

	if f, ok := matchMuxEntry(req.URL.Path, m.filter); ok {
		if res := f.ServeCtx(ctx); res == R_RETURN {
			return
		}
	}

	handle, _ := matchMuxEntry(req.URL.Path, m.handle)
	if handle != nil {
		handle.ServeCtx(ctx)
	} else {
		http.NotFound(rw, req)
	}
}

func matchMuxEntry(path string, mes []*muxEntry) (Handle, bool) {
	if len(mes) == 0 {
		return nil, false
	}

	for _, m := range mes {
		if h, ok := m.match(path); ok {
			return h, true
		}
	}
	return nil, false
}
