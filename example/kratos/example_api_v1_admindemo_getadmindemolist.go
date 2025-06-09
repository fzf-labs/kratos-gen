package service

import (
	"context"

	"github.com/fzf-labs/godb/orm/condition"
	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// GetAdminDemoList 系统-用户-列表数据查询
func (e *ExampleApiV1AdminDemoService) GetAdminDemoList(ctx context.Context, req *pb.GetAdminDemoListReq) (*pb.GetAdminDemoListReply, error) {
	resp := &pb.GetAdminDemoListReply{
		Total: 0,
		List:  []*pb.AdminDemoInfo{},
	}
	param := &condition.Req{
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Query:    []*condition.QueryParam{},
		Order: []*condition.OrderParam{
			{
				Field: "created_at",
				Order: condition.DESC,
			},
		},
	}
	list, p, err := e.adminDemoRepo.FindMultiCacheByCondition(ctx, param)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	resp.Total = int32(p.Total)
	if len(list) > 0 {
		for _, v := range list {
			resp.List = append(resp.List, &pb.AdminDemoInfo{
				Id: v.ID,
				// TODO
			})
		}
	}
	return resp, nil
}
