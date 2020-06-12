package eventhandler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/vardius/go-api-boilerplate/cmd/user/internal/domain/user"
	"github.com/vardius/go-api-boilerplate/pkg/auth/oauth2"
	"github.com/vardius/go-api-boilerplate/pkg/domain"
	"github.com/vardius/go-api-boilerplate/pkg/eventbus"
)

// WhenUserAccessTokenWasRequested handles event
func WhenUserAccessTokenWasRequested(tokenProvider oauth2.TokenProvider) eventbus.EventHandler {
	fn := func(ctx context.Context, event domain.Event) {
		// this goroutine runs independently to request's goroutine,
		// therefor recover middleware will not recover from panic to prevent crash
		defer recoverEventHandler()

		logger := GetLogger(ctx)
		logger.Info(ctx, "[EventHandler] %s\n", event.Payload)

		e := user.WasRegisteredWithEmail{}

		err := json.Unmarshal(event.Payload, &e)
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", err)
			return
		}

		token, err := tokenProvider.RetrieveToken(ctx, string(e.Email))
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", err)
			return
		}

		b, err := json.Marshal(token)
		if err != nil {
			logger.Error(ctx, "[EventHandler] Error: %v\n", err)
			return
		}

		// @TODO: send token with an email as magic link
		log.Printf("[EventHandler] Access Token: %s\n", string(b))
	}

	return fn
}
