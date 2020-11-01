package repo

import (
	"context"
	"dbcache/models"
	"encoding/json"
	"fmt"
	"log"

	redis "github.com/go-redis/redis/v8"
)

const redisItemTag = "i"
const redisTransportTag = "t"
const redisTransportItemTag = "ti"

type CacheWrapper struct {
	CacheV      Cache
	RepositoryV Repository
}

func NewCacheWrapper(c Cache, r Repository) Wrapper {
	return &CacheWrapper{CacheV: c, RepositoryV: r}
}

func (w *CacheWrapper) GetCache() Cache {
	return w.CacheV
}

func (w *CacheWrapper) GetRepository() Repository {
	return w.RepositoryV
}

func (w *CacheWrapper) GetItem(ctx context.Context, itemId string) (*models.Item, error) {
	key := fmt.Sprintf("%s:%s", redisItemTag, itemId)
	c, err := w.CacheV.Get(ctx, key)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	//If there is no data in redis
	if err == redis.Nil {
		log.Println("Not found in redis")
		item, err := w.RepositoryV.GetItem(ctx, itemId)

		if err != nil {
			log.Println("Not found in db")
			return nil, err
		}

		//add data to cache
		//Maby make in goroutine
		key := fmt.Sprintf("%s:%s", redisItemTag, item.Id)
		err = w.CacheV.Set(ctx, key, map[string]interface{}{redisTransportTag: item.TransportId, "n": item.Number})

		if err != nil {
			log.Println(err)
		}

		return item, nil
	}

	//If data have founded in redis cache
	log.Println("Data have founded: ", c, err)
	return &models.Item{Id: itemId, TransportId: (*c)[redisTransportTag], Number: (*c)["n"]}, nil
}

func (w *CacheWrapper) GetTransport(ctx context.Context, transportId string) (*models.Transport, error) {
	key := fmt.Sprintf("%s:%s", redisTransportTag, transportId)
	c, err := w.CacheV.Get(ctx, key)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	//If there is no data in redis
	if err == redis.Nil {
		log.Println("Not found in redis")
		transport, err := w.RepositoryV.GetTransport(ctx, transportId)

		if err != nil {
			log.Println("Not found in db")
			return nil, err
		}

		//add data to cache
		//Maby make in goroutine
		key := fmt.Sprintf("%s:%s", redisTransportTag, transport.Id)
		err = w.CacheV.Set(ctx, key, map[string]interface{}{"n": transport.Number})

		if err != nil {
			log.Println(err)
		}

		return transport, nil
	}

	//If data have founded in redis cache
	log.Println("Data have founded: ", c, err)
	return &models.Transport{Id: transportId, Number: (*c)["n"]}, nil
}

func (w *CacheWrapper) GetTransportItems(ctx context.Context, transportId string) (*models.TransportItemView, error) {
	key := fmt.Sprintf("%s:%s", redisTransportItemTag, transportId)
	c, err := w.CacheV.Get(ctx, key)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	//If there is no data in redis
	if err == redis.Nil {
		log.Println("Not found in redis")
		transportView, err := w.RepositoryV.GetTransportItemView(ctx, transportId)

		if err != nil {
			log.Println("Not found in db")
			return nil, err
		}

		//add data to cache
		//Maby make in goroutine
		viewKey := fmt.Sprintf("%s:%s", redisTransportItemTag, transportView.Transport.Id)
		err = w.CacheV.Set(ctx, viewKey, map[string]interface{}{"n": transportView.Transport.Number})

		if err != nil {
			log.Println(err)
		}

		for _, item := range transportView.Items {
			itemKey := fmt.Sprintf("%s:%s", redisItemTag, item.Id)

			var itemStruct models.Item = models.Item{Id: item.Id,
				Number: item.Number, TransportId: item.TransportId}

			itemVal, err := json.Marshal(itemStruct)

			if err != nil {
				return nil, err
			}

			viewKey = fmt.Sprintf("%s:%s", redisTransportItemTag, item.TransportId)
			err = w.CacheV.Set(ctx, viewKey, map[string]interface{}{itemKey: string(itemVal)})

			if err != nil {
				log.Println(err)
			}
		}

		return transportView, nil
	}

	//If data have founded in redis cache
	log.Println("Data have founded: ", c, err)
	var transportView models.TransportItemView
	var items = make([]models.Item, 0)
	for index, item := range *c {
		log.Println(index, item)
		if index == "n" {
			transportView.Transport.Id = transportId
			transportView.Transport.Number = item
		} else {
			var unmarshaledItem models.Item
			err := json.Unmarshal([]byte(item), &unmarshaledItem)

			if err != nil {
				return nil, err
			}

			items = append(items, unmarshaledItem)
		}
	}
	transportView.Items = items

	//return &models.Transport{Id: transportId, Number: (*c)["n"]}, nil
	return &transportView, nil
}
