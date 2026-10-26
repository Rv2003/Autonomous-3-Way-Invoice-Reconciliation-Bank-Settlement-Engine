package pdf

import "context"

type Service interface {
	ListPdf(ctx context.Context)(error)
}


type svc struct{

}

func NewService() Service{
	return &svc{}
}

func(s *svc) ListPdf(ctx context.Context) error{
	return nil
}