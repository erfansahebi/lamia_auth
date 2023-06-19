package handler

import (
	"context"
	"github.com/erfansahebi/lamia_auth/di"
)

type Handler struct {
	AppCtx context.Context
	Di     di.DIContainerInterface
}
