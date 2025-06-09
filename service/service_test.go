package service

import "testing"

func TestService_Service(t *testing.T) {
	NewService("../example/api", "../example/kratos").Run()
}
