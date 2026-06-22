package controllers

import (
	"SungClip/internal/types"
	"context"
)

type IControllers interface {
	VideoIngestion(ctx context.Context, request *types.RequestVideoIngestion) (response *types.ResponseVideoIngestion, err error)
	VideoEditing(ctx context.Context, request *types.RequestVideoEditing) (response *types.ResponseVideoEditing, err error)
}