package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"math/rand"
)

var codeTplId = "18888"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(codeRepository *repository.CodeRepository, service sms.Service) *CodeService {
	return &CodeService{
		repo:   codeRepository,
		smsSvc: service,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	//生成一个验证码
	//塞进redis
	//发出
	code := svc.generateCode()
	err := svc.repo.Store(ctx, code, biz, phone)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//if err != nil {
	//
	//}
	return nil
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}

func (svc *CodeService) VerifyV1(ctx context.Context, biz string, phone string, inputCode string) error {
	return nil
}
