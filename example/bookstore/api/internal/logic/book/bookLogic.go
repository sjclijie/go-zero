package logic

import (
	"bookstore/api/internal/svc"
	"bookstore/api/internal/types"
	"context"

	"github.com/sjclijie/go-zero/core/logx"
)

type BookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookLogic {
	return &BookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookLogic) Add(req types.AddReq) error {
	// todo: add your logic here and delete this line

	return nil
}
