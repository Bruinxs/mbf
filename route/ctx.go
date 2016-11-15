package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Bruinxs/log"
	"github.com/Bruinxs/util/uv"
)

type ResultData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Err  string      `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type SessionCtx struct {
	Rw     http.ResponseWriter
	Req    *http.Request
	Mux    *Mux
	Values map[string]interface{}
}

func NewSessionCtx(rw http.ResponseWriter, req *http.Request) *SessionCtx {
	ctx := &SessionCtx{Rw: rw, Req: req, Values: make(map[string]interface{})}
	return ctx
}

func (sc *SessionCtx) Value(key interface{}) interface{} {
	if key == nil {
		return nil
	}

	strKey, strOk := key.(string)
	if strOk {
		val := sc.Values[strKey]
		if val != nil {
			return val
		}

		sv := sc.Req.FormValue(strKey)
		if sv != "" {
			return sv
		}
		sv = sc.Req.PostFormValue(strKey)
		if sv != "" {
			return sv
		}
	}

	if val := sc.Req.Context().Value(key); val != nil {
		return val
	}

	if strOk && sc.Mux != nil && len(sc.Mux.Values) > 0 {
		val := sc.Mux.Values[strKey]
		if val != nil {
			return val
		}
	}

	return nil
}

//
func (sc *SessionCtx) SetVal(key string, val interface{}) {
	sc.Values[key] = val
}

func (sc *SessionCtx) StrVal(key string) string {
	if sv, ok := sc.Values[key].(string); ok {
		return sv
	}
	return ""
}

func (sc *SessionCtx) FetchVal(format string, args ...interface{}) error {
	return uv.Fetch(sc, format, args...)
}

func (sc *SessionCtx) UnmarshalJson(v interface{}) error {
	data, err := ioutil.ReadAll(sc.Req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

//
func (sc *SessionCtx) RData(data *ResultData) Result {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.E("[SessionCtx.RData] marshal data to json err(%v)", err)
		return R_RETURN
	}

	sc.Rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, err = sc.Rw.Write(bytes)
	if err != nil {
		log.E("[SessionCtx.RData] write json err(%v)", err)
	}
	return R_RETURN
}

func (sc *SessionCtx) Success(data interface{}) Result {
	return sc.RData(&ResultData{0, "", "", data})
}

func (sc *SessionCtx) Err(code int, msg string, err error) Result {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return sc.RData(&ResultData{code, msg, errStr, nil})
}
