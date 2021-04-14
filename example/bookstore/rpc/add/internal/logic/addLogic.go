package logic

import (
	"context"
	"github.com/davecgh/go-spew/spew"

	"bookstore/rpc/add/add"
	"bookstore/rpc/add/internal/svc"

	"github.com/sjclijie/go-zero/core/logx"
)

type AddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddLogic) Add(in *add.AddReq) (*add.AddResp, error) {
	// todo: add your logic here and delete this line
	spew.Dump(in.GetBook(), in.GetPrice())
	return &add.AddResp{Ok: true}, nil
}
