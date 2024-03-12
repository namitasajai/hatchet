package ingestor

import (
	"context"
	"fmt"

	"github.com/steebchen/prisma-client-go/runtime/types"

	"github.com/hatchet-dev/hatchet/internal/datautils"
	"github.com/hatchet-dev/hatchet/internal/msgqueue"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/prisma/db"
	"github.com/hatchet-dev/hatchet/internal/services/ingestor/contracts"
	"github.com/hatchet-dev/hatchet/internal/services/shared/tasktypes"
	"github.com/hatchet-dev/hatchet/internal/telemetry"
)

type Ingestor interface {
	contracts.EventsServiceServer
	IngestEvent(ctx context.Context, tenantId, eventName string, data any) (*db.EventModel, error)
	IngestReplayedEvent(ctx context.Context, tenantId string, replayedEvent *db.EventModel) (*db.EventModel, error)
}

type IngestorOptFunc func(*IngestorOpts)

type IngestorOpts struct {
	eventRepository repository.EventRepository
	logRepository   repository.LogsRepository
	mq              msgqueue.MessageQueue
}

func WithEventRepository(r repository.EventRepository) IngestorOptFunc {
	return func(opts *IngestorOpts) {
		opts.eventRepository = r
	}
}

func WithLogRepository(r repository.LogsRepository) IngestorOptFunc {
	return func(opts *IngestorOpts) {
		opts.logRepository = r
	}
}

func WithMessageQueue(mq msgqueue.MessageQueue) IngestorOptFunc {
	return func(opts *IngestorOpts) {
		opts.mq = mq
	}
}

func defaultIngestorOpts() *IngestorOpts {
	return &IngestorOpts{}
}

type IngestorImpl struct {
	contracts.UnimplementedEventsServiceServer

	eventRepository repository.EventRepository
	logRepository   repository.LogsRepository
	mq              msgqueue.MessageQueue
}

func NewIngestor(fs ...IngestorOptFunc) (Ingestor, error) {
	opts := defaultIngestorOpts()

	for _, f := range fs {
		f(opts)
	}

	if opts.eventRepository == nil {
		return nil, fmt.Errorf("event repository is required. use WithEventRepository")
	}

	if opts.logRepository == nil {
		return nil, fmt.Errorf("log repository is required. use WithLogRepository")
	}

	if opts.mq == nil {
		return nil, fmt.Errorf("task queue is required. use WithMessageQueue")
	}

	return &IngestorImpl{
		eventRepository: opts.eventRepository,
		logRepository:   opts.logRepository,
		mq:              opts.mq,
	}, nil
}

func (i *IngestorImpl) IngestEvent(ctx context.Context, tenantId, key string, data any) (*db.EventModel, error) {
	ctx, span := telemetry.NewSpan(ctx, "ingest-event")
	defer span.End()

	// transform data to a JSON object
	jsonType, err := datautils.ToJSONType(data)

	if err != nil {
		return nil, fmt.Errorf("could not convert event data to JSON: %w", err)
	}

	event, err := i.eventRepository.CreateEvent(ctx, &repository.CreateEventOpts{
		TenantId: tenantId,
		Key:      key,
		Data:     jsonType,
	})

	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	telemetry.WithAttributes(span, telemetry.AttributeKV{
		Key:   "event_id",
		Value: event.ID,
	})

	err = i.mq.AddMessage(context.Background(), msgqueue.EVENT_PROCESSING_QUEUE, eventToTask(event))

	if err != nil {
		return nil, fmt.Errorf("could not add event to task queue: %w", err)
	}

	return event, nil
}

func (i *IngestorImpl) IngestReplayedEvent(ctx context.Context, tenantId string, replayedEvent *db.EventModel) (*db.EventModel, error) {
	ctx, span := telemetry.NewSpan(ctx, "ingest-replayed-event")
	defer span.End()

	// transform data to a JSON object
	var data *types.JSON

	if jsonType, ok := replayedEvent.Data(); ok {
		data = &jsonType
	}

	event, err := i.eventRepository.CreateEvent(ctx, &repository.CreateEventOpts{
		TenantId:      tenantId,
		Key:           replayedEvent.Key,
		Data:          data,
		ReplayedEvent: &replayedEvent.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	err = i.mq.AddMessage(context.Background(), msgqueue.EVENT_PROCESSING_QUEUE, eventToTask(event))

	if err != nil {
		return nil, fmt.Errorf("could not add event to task queue: %w", err)
	}

	return event, nil
}

func eventToTask(e *db.EventModel) *msgqueue.Message {
	payload, _ := datautils.ToJSONMap(tasktypes.EventTaskPayload{
		EventId: e.ID,
	})

	metadata, _ := datautils.ToJSONMap(tasktypes.EventTaskMetadata{
		EventKey: e.Key,
		TenantId: e.TenantID,
	})

	return &msgqueue.Message{
		ID:       "event",
		Payload:  payload,
		Metadata: metadata,
		Retries:  3,
	}
}
