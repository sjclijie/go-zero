package handler

import (
	"net/http"

	"bookstore/api/internal/logic/book"
	"bookstore/api/internal/svc"
	"bookstore/api/internal/types"

	"github.com/sjclijie/go-zero/rest/httpx"
)

func AddHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewBookLogic(r.Context(), ctx)
		err := l.Add(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
