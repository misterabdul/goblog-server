package service

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/queue/client"
)

type Service struct {
	dbConn      *mongo.Database
	queueClient *client.QueueClient

	User         *user
	RevokedToken *revokedToken
	Category     *category
	Post         *post
	Comment      *comment
	Page         *page
	Notification *notification
}

func NewService(
	dbConn *mongo.Database,
	queueClient *client.QueueClient,
) (service *Service) {
	return &Service{
		dbConn:      dbConn,
		queueClient: queueClient,

		User:         newUserService(dbConn),
		RevokedToken: newRevokedTokenService(dbConn),
		Category:     newCategoryService(dbConn),
		Post:         newPostService(dbConn),
		Comment:      newCommentService(dbConn),
		Page:         newPageService(dbConn),
		Notification: newNotificationService(dbConn)}
}
