package service

import (
	"context"

	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// CreateAdminDemo 系统-用户-创建一条数据
func (e *ExampleApiV1AdminDemoService) CreateAdminDemo(ctx context.Context, req *pb.CreateAdminDemoReq) (*pb.CreateAdminDemoReply, error) {
	resp := &pb.CreateAdminDemoReply{}
	data := e.adminDemoRepo.NewData()
	// TODO
	err := e.adminDemoRepo.CreateOneCache(ctx, data)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	resp.Id = data.ID
	return resp, nil
}
