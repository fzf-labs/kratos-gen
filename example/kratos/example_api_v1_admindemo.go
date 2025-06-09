package service

import (
	pb "github.com/fzf-labs/kratos-gen/example/api/v1"
	"github.com/go-kratos/kratos/v2/log"
)

func NewExampleApiV1AdminDemoService(
	logger log.Logger,
	adminDemoRepo *data.AdminDemoRepo,
) *ExampleApiV1AdminDemoService {
	l := log.NewHelper(log.With(logger, "module", "service/adminDemo"))
	return &ExampleApiV1AdminDemoService{
		log:           l,
		adminDemoRepo: adminDemoRepo,
	}
}

type ExampleApiV1AdminDemoService struct {
	pb.UnimplementedAdminDemoServer
	log           *log.Helper
	adminDemoRepo *data.AdminDemoRepo
}
