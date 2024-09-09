package errors

import (
	"context"

	"github.com/melisource/fury_cards-go-toolkit/pkg/log/v2"
)

// LogService
type mockLogService struct {
	log.LogService
	ErrorlnStub func(ctx context.Context, msg string, fields ...log.Field)
	InfolnStub  func(ctx context.Context, msg string, fields ...log.Field)
	WarnStub    func(c context.Context, msg string, fields ...log.Field)
	SetStacktub func(stacktrace string) log.Field
	AnyStub     func(key string, value interface{}) log.Field
}

func (m mockLogService) Errorln(c context.Context, msg string, fields ...log.Field) {
	m.ErrorlnStub(c, msg, fields...)
}

func (m mockLogService) Warn(c context.Context, msg string, fields ...log.Field) {
	m.WarnStub(c, msg, fields...)
}

func (m mockLogService) SetStack(stacktrace string) log.Field {
	return m.SetStacktub(stacktrace)
}

func (m mockLogService) Infoln(c context.Context, msg string, fields ...log.Field) {
	m.InfolnStub(c, msg, fields...)
}

func (m mockLogService) Any(key string, value interface{}) log.Field {
	return m.AnyStub(key, value)
}
