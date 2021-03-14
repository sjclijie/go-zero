package logic

import (
	"context"
	"bookstore/api/internal/svc"	"bookstore/api/internal/types"
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

func (l *BookLogic) Add(req types.AddReq) (*types.AddResp, error) {
	// todo: add your logic here and delete this line

	return &types.AddResp{}, nil
}

func (l *BookLogic) Check(req types.CheckReq) (*types.CheckResp, error) {
	// todo: add your logic here and delete this line

	return &types.CheckResp{}, nil
}

