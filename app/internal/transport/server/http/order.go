package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewGetOrder(r service.OrderReader) *getOrderHandler {
	return &getOrderHandler{
		reader: r,
	}
}

type getOrderHandler struct {
	reader service.OrderReader
}

func (h *getOrderHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		orders, err := h.reader.Get()
		if err != nil {
			return fail(ctx, err)
		}

		return json(ctx, orders)
	}
}

func NewPostOrder(c service.OrderCreator) *postOrderHandler {
	return &postOrderHandler{
		creator: c,
	}
}

type postOrderHandler struct {
	creator service.OrderCreator
}

func (h *postOrderHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		idempotenceKey, err := parseIdempotenceKey(ctx)
		if err != nil {
			return fail(ctx, err)
		}

		var data service.OrderData
		err = ctx.BodyParser(&data)
		if err != nil {
			return fail(ctx, err)
		}

		o := service.Order{
			Id:        uuid.New().String(),
			OrderData: data,
		}

		err = h.creator.Create(idempotenceKey, o.ToDomain())
		if err != nil {
			return fail(ctx, err)
		}

		return RespondOk(ctx)
	}
}

func parseIdempotenceKey(ctx *fiber.Ctx) (uuid.UUID, error) {
	headers := ctx.GetReqHeaders()
	if key, ok := headers["X-Idempotence-Key"]; ok {

		keyUuid, err := uuid.Parse(key)
		if err != nil {

		}

		return keyUuid, nil
	}

	return uuid.New(), errors.New("idempotence key not found")
}
