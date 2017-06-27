package route

import (
	"fmt"
	"net/http"
	"regexp"
)

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

func (m *Mux) HandFunc(pattern string, hf func(*SessionCtx) Result) {
	m.Hand(pattern, HandleFunc(hf))
}

func (m *Mux) Filter(pattern string, handle Handle) {
	reg := regexp.MustCompile(pattern)
	m.filter = append(m.filter, &muxEntry{reg, handle})
}

func (m *Mux) FilterFunc(pattern string, hf func(*SessionCtx) Result) {
	m.Filter(pattern, HandleFunc(hf))
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fmt.Sprintf("\nserver panic err:\n%v", err)))
		}
	}()

	ctx := NewSessionCtx(w, r)
	ctx.Mux = m

	handle := m.matchMuxEntry(r.URL.Path)
	for _, h := range handle {
		if res := h.ServeCtx(ctx); res == R_RETURN {
			break
		}
	}

	if err := ctx.WriteResp(); err != nil {
		panic(err)
	}
}

func (m *Mux) matchMuxEntry(path string) []Handle {
	handle := make([]Handle, 0, 2)
	for _, me := range m.filter {
		if f, ok := me.match(path); ok {
			handle = append(handle, f)
		}
	}

	var H Handle
	for _, me := range m.handle {
		if h, ok := me.match(path); ok {
			H = h
			break
		}
	}
	if H == nil {
		H = HandleFunc(NotFound)
	}

	return append(handle, H)
}
