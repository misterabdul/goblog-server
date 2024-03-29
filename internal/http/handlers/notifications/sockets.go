package notifications

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/misterabdul/goblog-server/internal/database/models"
	"github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/http/responses"
	internalGin "github.com/misterabdul/goblog-server/internal/pkg/gin"
	"github.com/misterabdul/goblog-server/internal/service"
)

func ServeListenedNotifications(
	maxCtxDuration time.Duration,
	svc *service.Service,
) (handler gin.HandlerFunc) {

	return func(c *gin.Context) {
		var (
			me          *models.UserModel
			messageChan = make(chan string)
			upgrader    websocket.Upgrader
			wsConn      *websocket.Conn
			err         error
		)

		if me, err = authenticate.GetAuthenticatedUser(c); err != nil {
			responses.Unauthenticated(c, err)
			return
		}
		upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		if wsConn, err = upgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
			responses.InternalServerError(c, err)
			return
		}
		defer wsConn.Close()

		go checkNotifications(svc, me, &messageChan)
		for {
			wsConn.WriteMessage(websocket.TextMessage, []byte(<-messageChan))
		}
	}
}

func checkNotifications(
	svc *service.Service,
	me *models.UserModel,
	messageChan *chan string,
) {
	var (
		ctx           = context.TODO()
		notifications []*models.NotificationModel
		latestCheck   = time.Now()
		messageBuff   string
		err           error
	)

	for {
		time.Sleep(3 * time.Second)
		if notifications, err = svc.Notification.GetMany(ctx, bson.M{
			"$and": []bson.M{
				{"owner.username": me.Username},
				{"createdat": bson.M{"$gt": primitive.NewDateTimeFromTime(latestCheck)}}}},
			internalGin.CreateFindOptions(25, 1, "createdat", false),
		); err != nil {
			continue
		}
		messageBuff = fmt.Sprintf("There is %d new notification(s)", len(notifications))
		(*messageChan) <- messageBuff
	}
}
