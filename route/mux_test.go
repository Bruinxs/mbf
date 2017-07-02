package route_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/bruinxs/mbf/route"
	"github.com/bruinxs/ts"
	"github.com/bruinxs/ts/th"
	"github.com/bruinxs/util/ut"
)

func TestMux(t *testing.T) {
	list := []string{}
	mux := NewMux()
	hs := httptest.NewServer(mux)

	//no found
	_, err := th.Get(hs.URL, "", nil)
	if err == nil {
		t.Error(err)
		return
	}

	//handle
	mux.HandFunc(".*", func(ctx *SessionCtx) Result {
		list = append(list, "handle")
		return ctx.Success(ut.M{"msg": "fake"})
	})
	_, err = th.Get(hs.URL, "test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !ts.CmpStr_Strict(list, []string{"handle"}) {
		t.Errorf("list(%v) ill", list)
		return
	}

	//filter continue
	list = list[:0]
	mux.FilterFunc("^.*$", func(ctx *SessionCtx) Result {
		list = append(list, "all")
		return R_CONTINUE
	})
	mux.FilterFunc("filter.?.?.?", func(ctx *SessionCtx) Result {
		list = append(list, "filter")
		return R_CONTINUE
	})
	_, err = th.Get(hs.URL, "filter", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !ts.CmpStr_Strict(list, []string{"all", "filter", "handle"}) {
		t.Errorf("list(%v) ill", list)
		return
	}

	//filter return
	list = list[:0]
	mux.FilterFunc("filterext", func(ctx *SessionCtx) Result {
		list = append(list, "ext")
		return R_RETURN
	})
	_, err = th.Get(hs.URL, "filterext", nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !ts.CmpStr_Strict(list, []string{"all", "filter", "ext"}) {
		t.Errorf("list(%v) ill", list)
		return
	}

	//hand continue
	mux.HandFunc("handle", func(ctx *SessionCtx) Result {
		return ctx.Success(ut.M{"k1": "v1"})
	})
	mux.HandFunc("handle", func(ctx *SessionCtx) Result {
		return ctx.Success(ut.M{"k2": 222})
	})
	res, err := th.Get(hs.URL, "handle", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("res-> ", res)
}
