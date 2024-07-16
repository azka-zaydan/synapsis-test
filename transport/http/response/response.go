package response

import (
	"encoding/json"
	"net/http"

	"github.com/azka-zaydan/synapsis-test/shared/failure"
	"github.com/azka-zaydan/synapsis-test/shared/logger"
	"github.com/gofiber/fiber/v2"
)

type Base struct {
	Message  *string      `json:"message,omitempty"`
	Data     *interface{} `json:"data,omitempty"`
	Error    *string      `json:"error,omitempty"`
	Metadata *interface{} `json:"metadata,omitempty"`
}

func WithMetadata(c *fiber.Ctx, code int, data interface{}, metadata interface{}) error {
	err := respond(c, code, fiber.Map{"data": &data, "metadata": &metadata})
	return err
}

func WithMessage(c *fiber.Ctx, code int, message string) error {
	err := respond(c, code, fiber.Map{"message": &message})
	return err
}

func WithJSON(c *fiber.Ctx, code int, jsonPayload interface{}) error {
	err := respond(c, code, fiber.Map{"data": &jsonPayload})
	return err
}

func WithError(c *fiber.Ctx, err error) error {
	code := failure.GetCode(err)
	errMsg := err.Error()
	err = respond(c, code, fiber.Map{"error": &errMsg})
	return err
}

func WithPreparingShutdown(w http.ResponseWriter) {
	message := "SERVER PREPARING TO SHUT DOWN"
	respondWithWriter(w, http.StatusServiceUnavailable, fiber.Map{"message": &message})
}

func respond(c *fiber.Ctx, code int, payload interface{}) error {
	err := c.Status(code).JSON(payload)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func respondWithWriter(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		logger.ErrorWithStack(err)
	}
}
