package service

import gen "{{package}}/thrift/{{owner}}/{{project}}"
import "context"

func New{{class}}Service() *{{class}}ServiceImpl {
	return &{{class}}ServiceImpl{}
}

type {{class}}ServiceImpl struct {
}

func (svc {{class}}ServiceImpl) Echo(ctx context.Context, request *gen.EchoRequest) (r *gen.EchoResponse, err error) {
	return &gen.EchoResponse{Content: request.Content}, nil
}

func (svc {{class}}ServiceImpl) IsHealthy(ctx context.Context) (r bool, err error) {
	return true, nil
}