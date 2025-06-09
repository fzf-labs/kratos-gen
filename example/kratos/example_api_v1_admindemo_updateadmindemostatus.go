package service

import (
	"context"

	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
)

// UpdateAdminDemoStatus 系统-用户-更新状态
func (e *ExampleApiV1AdminDemoService) UpdateAdminDemoStatus(ctx context.Context, req *pb.UpdateAdminDemoStatusReq) (*pb.UpdateAdminDemoStatusReply, error) {
	resp := &pb.UpdateAdminDemoStatusReply{}
	data, err := e.adminDemoRepo.FindOneCacheByID(ctx, req.GetId())
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	if data == nil || data.ID == "" {
		return nil, pb.ErrorReasonDataRecordNotFound()
	}
	oldData := e.adminDemoRepo.DeepCopy(data)
	data.Status = req.GetStatus()
	err = e.adminDemoRepo.UpdateOneCacheWithZero(ctx, data, oldData)
	if err != nil {
		return nil, pb.ErrorReasonDataSQLError(pb.WithError(err))
	}
	return resp, nil
}
