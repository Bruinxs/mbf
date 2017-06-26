package route

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Bruinxs/util/ut"
	"github.com/bruinxs/util/uv"
)

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

		sv := sc.R.FormValue(strKey)
		if sv != "" {
			return sv
		}
		sv = sc.R.PostFormValue(strKey)
		if sv != "" {
			return sv
		}
	}

	if val := sc.R.Context().Value(key); val != nil {
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

func (sc *SessionCtx) FetchVal(format string, args ...interface{}) error {
	return uv.Fetch(sc, format, args...)
}

//UnmarshalJSON parse json data from request body
func (sc *SessionCtx) UnmarshalJSON(v interface{}) error {
	data, err := ioutil.ReadAll(sc.R.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func (sc *SessionCtx) SetVal(key string, val interface{}) {
	sc.Values[key] = val
}

//StrVal string value
func (sc *SessionCtx) StrVal(key string) string {
	if sv, ok := sc.Values[key].(string); ok {
		return sv
	}
	return ""
}

//BVal body json value
func (sc *SessionCtx) BVal() ut.M {
	if v := sc.Value("bodyJM"); v != nil {
		return v.(ut.M)
	}
	var m ut.M
	err := sc.UnmarshalJSON(&m)
	if err != nil {
		panic(err)
	}
	sc.SetVal("bodyJM", m)
	return m
}
