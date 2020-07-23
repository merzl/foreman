package execution

import (
	"context"
	"github.com/kopaygorodsky/brigadier/pkg/pubsub/endpoint"
	"github.com/kopaygorodsky/brigadier/pkg/pubsub/message"
	"github.com/kopaygorodsky/brigadier/pkg/pubsub/transport/pkg"
	"github.com/kopaygorodsky/brigadier/pkg/runtime/scheme"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	LogLevel = iota
)

type MessageExecutionCtx interface {
	Message() *message.Message
	Context() context.Context
	Valid() bool
	Send(message *message.Message, options ...endpoint.DeliveryOption) error
	Return(delay time.Duration) error
	LogMessage(msg string, level string)
}

type messageExecutionCtx struct {
	isValid bool
	ctx     context.Context
	inPkg   pkg.IncomingPkg
	message *message.Message
	router  endpoint.Router
	logger  *log.Logger
}

func (m messageExecutionCtx) Valid() bool {
	return m.isValid
}

func (m messageExecutionCtx) Context() context.Context {
	return m.ctx
}

func (m messageExecutionCtx) Send(message *message.Message, options ...endpoint.DeliveryOption) error {
	endpoints := m.router.Route(scheme.WithKey(message.Name))

	if len(endpoints) == 0 {
		m.logger.Warningf("No endpoints defined for message %s", message.Name)
		return nil
	}

	for _, endp := range endpoints {
		if err := endp.Send(m.ctx, message, message.Headers, options...); err != nil {
			m.logger.Errorf("Error sending message id %s", message.ID)
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m messageExecutionCtx) Return(delay time.Duration) error {
	for {
		select {
		case <-m.ctx.Done():
			m.logger.Infof("Context is closed, exiting without returning msg: %s, delay is too long", m.message.ID)
			return nil
		case <-time.After(delay):
			if err := m.Send(m.message); err != nil {
				m.logger.Errorf("error when returning a message %s", m.message.ID)
				return errors.Wrapf(err, "error when returning a message %s", m.message.ID)
			}
		}
	}
}

func (m messageExecutionCtx) Message() *message.Message {
	return m.message
}

func (m messageExecutionCtx) LogMessage(msg string, level string) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		//show msg on error lvl
		m.logger.Error(err, msg)
	}
	m.logger.Log(lvl, msg)
}

type MessageExecutionCtxFactory interface {
	CreateCtx(ctx context.Context, inPkg pkg.IncomingPkg, message *message.Message) MessageExecutionCtx
}

type messageExecutionCtxFactory struct {
	router endpoint.Router
	logger *log.Logger
}

func NewMessageExecutionCtxFactory(router endpoint.Router, logger *log.Logger) MessageExecutionCtxFactory {
	return &messageExecutionCtxFactory{router: router, logger: logger}
}

func (m messageExecutionCtxFactory) CreateCtx(ctx context.Context, inPkg pkg.IncomingPkg, message *message.Message) MessageExecutionCtx {
	return &messageExecutionCtx{ctx: ctx, inPkg: inPkg, message: message, router: m.router, logger: m.logger}
}

type NoDefinedEndpoints struct {
	error
}

func WithNoDefinedEndpoints(err error) error {
	return NoDefinedEndpoints{err}
}
