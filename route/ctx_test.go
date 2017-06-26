package route

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruinxs/ts"
	"github.com/bruinxs/ts/th"
	"github.com/bruinxs/util/ut"
)

func TestSessionCtx(t *testing.T) {
	var sc *SessionCtx
	var rd *ResultData
	var key, val interface{}
	var call func(sc *SessionCtx)
	hts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc = NewSessionCtx(w, r.WithContext(context.WithValue(r.Context(), key, val)))
		if call != nil {
			call(sc)
		}
		if rd != nil {
			if rd.Code == 0 {
				sc.Success(rd.Data)
			} else {
				sc.Err(rd.Code, rd.Msg, errors.New(rd.Err))
			}
		}
	}))

	//1,get
	call = func(sc *SessionCtx) {
		sc.SetVal("id", "123")
	}

	sc = nil
	key, val = "key", "val"
	rd = &ResultData{Code: 10, Msg: "test msg", Err: "err"}

	s1 := "str"
	s2 := "s1,s2,s3"
	i1 := 10
	f1 := 3.14
	res, err := th.GP_M(hts.URL, "", ut.M{"s1": s1, "s2": s2, "i1": i1, "f1": f1}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	//resp
	if g, w := res.IntV("code"), rd.Code; g != w {
		t.Errorf("code(%v) != %v", g, w)
		return
	}
	if g, w := res.StrV("msg"), rd.Msg; g != w {
		t.Errorf("msg(%v) != %v", g, w)
		return
	}
	if g, w := res.StrV("err"), rd.Err; g != w {
		t.Errorf("err(%v) != %v", g, w)
		return
	}

	//sc
	if g, w := sc.StrVal("id"), "123"; g != w {
		t.Errorf("sc id(%v) != %v", g, w)
		return
	}
	if g, w := sc.StrVal("fake"), ""; g != w {
		t.Errorf("sc fake(%v) != %v", g, w)
		return
	}

	var (
		_id, _s1, _key string
		_s2            []string
		_i1            int
		_f1            float64
	)
	err = sc.FetchVal(`
		id,m,0;
		s1,m,0;
		key,m,0;
		s2,m,n;
		i1,m,0;
		f1,m,0;
	`, &_id, &_s1, &_key, &_s2, &_i1, &_f1)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("_id(%v), _s1(%v), _key(%v), _s2(%v), _i1(%v), _f1(%v)", _id, _s1, _key, _s2, _i1, _f1)

	if _id != "123" {
		t.Errorf("_id(%v) != %v", _id, "123")
		return
	}
	if _s1 != s1 {
		t.Errorf("_s1(%v) != %v", _s1, s1)
		return
	}
	if _i1 != i1 {
		t.Errorf("_i1(%v) != %v", _i1, i1)
		return
	}
	if _f1 != f1 {
		t.Errorf("_f1(%v) != %v", _f1, f1)
		return
	}
	if !ts.CmpStr_Strict(_s2, []string{"s1", "s2", "s3"}) {
		t.Errorf("_s2(%v) != %v", _s2, s2)
		return
	}

	//2,post
	call = func(sc *SessionCtx) {
		var m ut.M
		err := sc.UnmarshalJSON(&m)
		if err != nil {
			t.Error(err)
			return
		}

		if m.StrV("s3") != "str3" {
			t.Errorf("s3(%v) != %v", m.StrV("s3"), "str3")
			return
		}
	}

	rd = &ResultData{Code: 0, Data: ut.M{"msg": "success"}}
	_, err = th.GP_M(hts.URL, "", nil, ut.M{"s3": "str3"})
	if err != nil {
		t.Error(err)
		return
	}
}
