package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	c      *gin.Context
	ctx    context.Context
	dbConn *mongo.Database
}

func New(
	c *gin.Context,
	ctx context.Context,
	dbConn *mongo.Database,
) *Service {
	return &Service{
		c:      c,
		ctx:    ctx,
		dbConn: dbConn,
	}
}
