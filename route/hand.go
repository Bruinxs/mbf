package route

import "net/http"

type Result int

const (
	R_RETURN = Result(iota)
	R_CONTINUE
)

type Handle interface {
	ServeCtx(ctx *SessionCtx) Result
}

type HandleFunc func(*SessionCtx) Result

func (hf HandleFunc) ServeCtx(sc *SessionCtx) Result {
	return hf(sc)
}

func NotFound(ctx *SessionCtx) Result {
	http.NotFound(ctx.Rw, ctx.Req)
	return R_RETURN
}
