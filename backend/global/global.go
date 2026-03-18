package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	// Logger *logrus.Logger
	Db      *gorm.DB
	RedisDB *redis.Client
)
