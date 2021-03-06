package route

import (
	"encoding/json"
	"net/http"

	"github.com/bruinxs/util/ut"
)

const (
	CODE_NOT_HANDLE = -1
	CODE_SUCCESS    = 0
)

type ResultData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Err  string `json:"err,omitempty"`
	Data ut.M   `json:"data,omitempty"`
}

type SessionCtx struct {
	W      http.ResponseWriter
	R      *http.Request
	Mux    *Mux
	Values map[string]interface{}
	Resp   *ResultData
}

func NewSessionCtx(w http.ResponseWriter, r *http.Request) *SessionCtx {
	ctx := &SessionCtx{W: w, R: r, Values: make(map[string]interface{}), Resp: &ResultData{Code: CODE_NOT_HANDLE}}
	return ctx
}

func (sc *SessionCtx) WriteResp() error {
	if sc.Resp.Code == CODE_NOT_HANDLE {
		return nil
	}
	bytes, err := json.Marshal(sc.Resp)
	if err != nil {
		return err
	}

	sc.W.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, err = sc.W.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (sc *SessionCtx) Success(data ut.M) Result {
	sc.Resp.Code = CODE_SUCCESS
	if sc.Resp.Data == nil {
		sc.Resp.Data = data
	} else {
		for key, val := range data {
			sc.Resp.Data[key] = val
		}
	}
	return R_CONTINUE
}

func (sc *SessionCtx) Err(code int, msg string, err error) Result {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	sc.Resp.Code, sc.Resp.Msg, sc.Resp.Err = code, msg, errStr
	return R_RETURN
}

func (sc *SessionCtx) Redirect(urlStr string) {
	http.Redirect(sc.W, sc.R, urlStr, http.StatusMovedPermanently)
}
