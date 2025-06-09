package service

import (
	"context"

	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// GetAdminDemoInfo 系统-用户-单条数据查询
func (e *ExampleApiV1AdminDemoService) GetAdminDemoInfo(ctx context.Context, req *pb.GetAdminDemoInfoReq) (*pb.GetAdminDemoInfoReply, error) {
	resp := &pb.GetAdminDemoInfoReply{}
	data, err := e.adminDemoRepo.FindOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	resp.Info = &pb.AdminDemoInfo{
		Id: data.ID,
		// TODO
	}
	return resp, nil
}
