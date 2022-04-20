package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/livekit/protocol/auth"
)

var (
	apiKey    string
	apiSecret string
	port      int64
)

type (
	TokenRequest struct {
		RoomName string `json:"roomName"`
		Username string `json:"userName"`
		Metadata string `json:"metadata"`
	}

	TokenResponse struct {
		Token string `json:"token"`
	}
)

func init() {
	flag.StringVar(&apiKey, "api_key", "", "api key")
	flag.StringVar(&apiSecret, "api_secret", "", "api secret")
	flag.Int64Var(&port, "port", 3000, "port listening on")
	flag.Parse()
}

func main() {
	app := fiber.New(fiber.Config{
		ReduceMemoryUsage:     true,
		DisableStartupMessage: true,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/tokens", func(c *fiber.Ctx) error {
		var req TokenRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "error decoding request")
		}

		token, err := getJoinToken(apiKey, apiSecret, req.RoomName, req.Username, req.Metadata)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error creating token: %v", err))
		}

		res := TokenResponse{
			Token: token,
		}

		return c.JSON(res)
	})

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func getJoinToken(apiKey, apiSecret, room, name, metadata string) (string, error) {
	canPublish := true
	canSubscribe := true

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		RoomCreate:   true,
		Room:         room,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	at.AddGrant(grant).
		SetIdentity(uuid.NewString()).
		SetName(name).
		SetMetadata(metadata).
		SetValidFor(time.Hour * 24)

	return at.ToJWT()
}
