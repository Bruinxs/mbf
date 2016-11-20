package route

import (
	"net/http/httptest"
	"testing"

	"github.com/Bruinxs/tu"
	"github.com/Bruinxs/tu/th"
)

func TestMux(t *testing.T) {
	list := []string{}
	mux := NewMux()
	hs := httptest.NewServer(mux)

	//no found
	_, err := th.GP_M(hs.URL, "", nil, nil)
	if err == nil {
		t.Error(err)
		return
	}

	//handle
	mux.HandFunc(".*", func(ctx *SessionCtx) Result {
		list = append(list, "handle")
		return ctx.Success("fake")
	})
	_, err = th.GP_M(hs.URL, "test", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !tu.CmpStr_Strict(list, []string{"handle"}) {
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
	_, err = th.GP_M(hs.URL, "filter", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !tu.CmpStr_Strict(list, []string{"all", "filter", "handle"}) {
		t.Errorf("list(%v) ill", list)
		return
	}

	//filter return
	list = list[:0]
	mux.FilterFunc("filterext", func(ctx *SessionCtx) Result {
		list = append(list, "ext")
		return R_RETURN
	})
	_, err = th.GP_M(hs.URL, "filterext", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !tu.CmpStr_Strict(list, []string{"all", "filter", "ext"}) {
		t.Errorf("list(%v) ill", list)
		return
	}
}
