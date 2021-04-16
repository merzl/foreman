package handlers

import (
	"context"
	"fmt"
	log "github.com/go-foreman/foreman/log"
	"github.com/go-foreman/foreman/pubsub/message"
	"github.com/go-foreman/foreman/pubsub/message/execution"
	"github.com/go-foreman/foreman/runtime/scheme"
	sagaPkg "github.com/go-foreman/foreman/saga"
	"github.com/go-foreman/foreman/saga/contracts"
	"github.com/go-foreman/foreman/saga/mutex"
	"github.com/pkg/errors"
	"time"
)

func NewSagaControlHandler(sagaStore sagaPkg.Store, mutex mutex.Mutex, sagaRegistry scheme.KnownTypesRegistry, logger log.Logger) *SagaControlHandler {
	return &SagaControlHandler{typesRegistry: sagaRegistry, store: sagaStore, mutex: mutex, logger: logger}
}

type SagaControlHandler struct {
	typesRegistry scheme.KnownTypesRegistry
	store         sagaPkg.Store
	mutex         mutex.Mutex
	logger        log.Logger
}

func (h SagaControlHandler) Handle(execCtx execution.MessageExecutionCtx) error {
	var (
		sagaInstance sagaPkg.Instance
		sagaCtx      sagaPkg.SagaContext
		err          error
	)

	ctx := execCtx.Context()
	msg := execCtx.Message()

	switch cmd := msg.Payload().(type) {
	case *contracts.StartSagaCommand:
		sagaInstance, err = h.createSaga(cmd)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := h.store.Create(ctx, sagaInstance); err != nil {
			return errors.Wrapf(err, "error  saving created saga `%s` with id %s to store", "", cmd.SagaId)
		}

		sagaCtx = sagaPkg.NewSagaCtx(execCtx, sagaInstance)

		if err := sagaInstance.Start(sagaCtx); err != nil {
			return errors.Wrapf(err, "error starting saga `%s`", sagaInstance.UID())
		}

	case *contracts.RecoverSagaCommand:
		if err := h.mutex.Lock(ctx, cmd.SagaId); err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			if err := h.mutex.Release(ctx, cmd.SagaId); err != nil {
				h.logger.Log(log.ErrorLevel, err)
			}
		}()

		sagaInstance, err = h.fetchSaga(ctx, cmd.SagaId)

		if err != nil {
			return errors.WithStack(err)
		}

		if !sagaInstance.Status().Failed() || sagaInstance.Status().Completed() || sagaInstance.Status().Recovering() || sagaInstance.Status().Compensating() {
			h.logger.Logf(log.InfoLevel, "Saga `%s` has status %s, you can't start recovering the process", sagaInstance.Status(), sagaInstance.UID())
			return nil
		}

		sagaCtx = sagaPkg.NewSagaCtx(execCtx, sagaInstance)

		if err := sagaInstance.Recover(sagaCtx); err != nil {
			return errors.Wrapf(err, "error recovering saga `%s`", sagaInstance.UID())
		}

	case *contracts.CompensateSagaCommand:
		if err := h.mutex.Lock(ctx, cmd.SagaId); err != nil {
			return errors.WithStack(err)
		}

		defer func() {
			if err := h.mutex.Release(ctx, cmd.SagaId); err != nil {
				h.logger.Log(log.ErrorLevel, err)
			}
		}()

		sagaInstance, err = h.fetchSaga(ctx, cmd.SagaId)

		if err != nil {
			return errors.WithStack(err)
		}

		if !sagaInstance.Status().Failed() || sagaInstance.Status().Compensating() {
			h.logger.Logf(log.InfoLevel, "Saga `%s` has status `%s`, you can't compensate the process", sagaInstance.UID(), sagaInstance.Status())
			return nil
		}

		sagaCtx = sagaPkg.NewSagaCtx(execCtx, sagaInstance)

		if err := sagaInstance.Compensate(sagaCtx); err != nil {
			return errors.Wrapf(err, "error compensating saga `%s`", sagaInstance.UID())
		}

	default:
		return errors.Errorf("unknown command type `%s` for SagaControlHandler. Supported: StartSagaCommand, RecoverSagaCommand, CompensateSagaCommand", msg.Payload().GroupKind().String())
	}

	for _, delivery := range sagaCtx.Deliveries() {
		outcomingMessage := message.NewOutcomingMessage(delivery.Payload, message.WithHeaders(msg.Headers()))
		if err := execCtx.Send(outcomingMessage, delivery.Options...); err != nil {
			execCtx.LogMessage(log.ErrorLevel, fmt.Sprintf("error sending delivery for saga %s. Delivery: (%v). %s", sagaCtx.SagaInstance().UID(), delivery, err))
			return errors.Wrapf(err, "error sending delivery for saga %s. Delivery: (%v)", sagaCtx.SagaInstance().UID(), delivery)
		}
		sagaCtx.SagaInstance().AttachEvent(sagaPkg.HistoryEvent{UID: msg.UID(), Payload: delivery.Payload, CreatedAt: time.Now(), SagaStatus: sagaInstance.Status().String()})
	}

	return h.store.Update(ctx, sagaInstance)
}

//saga is map[string]interface{} on this step
func (h SagaControlHandler) createSaga(startCmd *contracts.StartSagaCommand) (sagaPkg.Instance, error) {
	if startCmd.SagaId == "" {
		return nil, errors.Errorf("sagaId is empty")
	}

	if startCmd.Saga == nil {
		return nil, errors.Errorf("saga payload is nil")
	}

	saga, ok := startCmd.Saga.(sagaPkg.Saga)

	if !ok {
		return nil, errors.Errorf("error asserting that startCmd.Saga is Saga type")
	}

	return sagaPkg.NewSagaInstance(startCmd.SagaId, startCmd.ParentId, saga), nil
}

func (h SagaControlHandler) fetchSaga(ctx context.Context, sagaId string) (sagaPkg.Instance, error) {
	sagaInstance, err := h.store.GetById(ctx, sagaId)

	if err != nil {
		return nil, errors.Wrapf(err, "error fetching saga instance `%s` from store", sagaId)
	}

	if sagaInstance == nil {
		return nil, errors.Errorf("Saga instance `%s` not found", sagaId)
	}

	return sagaInstance, nil
}
