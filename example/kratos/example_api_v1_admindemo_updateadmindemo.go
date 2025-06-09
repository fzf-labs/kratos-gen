package service

import (
	"context"

	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// UpdateAdminDemo 系统-用户-更新一条数据
func (e *ExampleApiV1AdminDemoService) UpdateAdminDemo(ctx context.Context, req *pb.UpdateAdminDemoReq) (*pb.UpdateAdminDemoReply, error) {
	resp := &pb.UpdateAdminDemoReply{}
	data, err := e.adminDemoRepo.FindOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	if data == nil || data.ID == "" {
		return nil, pb.ErrorReasonDataRecordNotFound()
	}
	oldData := e.adminDemoRepo.DeepCopy(data)
	// TODO
	err = e.adminDemoRepo.UpdateOneCacheWithZero(ctx, data, oldData)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	return resp, nil
}
