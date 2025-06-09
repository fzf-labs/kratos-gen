package service

import (
	"context"

	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// DeleteAdminDemo 系统-用户-删除多条数据
func (e *ExampleApiV1AdminDemoService) DeleteAdminDemo(ctx context.Context, req *pb.DeleteAdminDemoReq) (*pb.DeleteAdminDemoReply, error) {
	resp := &pb.DeleteAdminDemoReply{}
	err := e.adminDemoRepo.DeleteOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	return resp, nil
}
