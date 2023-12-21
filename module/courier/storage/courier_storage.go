package storage

import (
	"context"
	"encoding/json"
	"errors"
	"geotask/module/courier/models"
	"github.com/redis/go-redis"
)

type CourierStorager interface {
	Save(ctx context.Context, courier models.Courier) error // сохранить курьера по ключу courier
	GetOne(ctx context.Context) (*models.Courier, error)    // получить курьера по ключу courier
}

type CourierStorage struct {
	storage *redis.Client
}

func NewCourierStorage(storage *redis.Client) CourierStorager {
	return &CourierStorage{storage: storage}
}

func (cs *CourierStorage) Save(ctx context.Context, courier models.Courier) error {
	// реализация сохранения в Redis
	courierJSON, err := json.Marshal(courier)
	if err != nil {
		return err
	}

	// Сохраняем в Redis
	err = cs.storage.Set(ctx, "courier", courierJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cs *CourierStorage) GetOne(ctx context.Context) (*models.Courier, error) {
	// реализация получения курьера из Redis
	// Получаем из Redis
	courierJSON, err := cs.storage.Get(ctx, "").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	// Преобразовываем JSON обратно в курьера
	var courier models.Courier
	err = json.Unmarshal([]byte(courierJSON), &courier)
	if err != nil {
		return nil, err
	}

	return &courier, nil
	//return nil, errors.New("not found")
}
