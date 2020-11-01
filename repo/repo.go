package repo

import (
	"context"
	"dbcache/models"
)

type Cache interface {
	Get(ctx context.Context, key string) (*map[string]string, error)
	Set(ctx context.Context, key string, val map[string]interface{}) error
}

type Repository interface {
	GetItem(ctx context.Context, itemId string) (*models.Item, error)
	GetTransport(ctx context.Context, transportId string) (*models.Transport, error)
	GetTransportItemView(ctx context.Context, transportId string) (*models.TransportItemView, error)
	CreateAlotItems(ctx context.Context, itemsConf models.CreateAlotItems) error
}

type Wrapper interface {
	GetItem(ctx context.Context, itemId string) (*models.Item, error)
	GetTransport(ctx context.Context, transportId string) (*models.Transport, error)
	GetTransportItems(ctx context.Context, transportId string) (*models.TransportItemView, error)

	GetCache() Cache
	GetRepository() Repository
}
