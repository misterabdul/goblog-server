package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	authenticationHandler "github.com/misterabdul/goblog-server/internal/http/handlers/authentications"
	categoryHandler "github.com/misterabdul/goblog-server/internal/http/handlers/categories"
	commentHandler "github.com/misterabdul/goblog-server/internal/http/handlers/comments"
	meHandler "github.com/misterabdul/goblog-server/internal/http/handlers/me"
	notificationHandler "github.com/misterabdul/goblog-server/internal/http/handlers/notifications"
	otherHandler "github.com/misterabdul/goblog-server/internal/http/handlers/others"
	postHandler "github.com/misterabdul/goblog-server/internal/http/handlers/posts"
	userHandler "github.com/misterabdul/goblog-server/internal/http/handlers/users"
	authenticateMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	authorizeMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authorize"
)

// Initialize all routes.
func InitRoutes(
	server *gin.Engine,
	dbConn *mongo.Database,
	maxCtxDuration time.Duration,
) {
	server.NoRoute(otherHandler.NotFound())

	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", otherHandler.Ping())

			v1.GET("/users", userHandler.GetPublicUsers(maxCtxDuration, dbConn))
			v1.GET("/user/:user", userHandler.GetPublicUser(maxCtxDuration, dbConn))

			v1.GET("/categories", categoryHandler.GetPublicCategories(maxCtxDuration, dbConn))
			v1.GET("/category/:category", categoryHandler.GetPublicCategory(maxCtxDuration, dbConn))
			v1.GET("/category/:category/posts", categoryHandler.GetPublicCategoryPosts(maxCtxDuration, dbConn))

			v1.GET("/posts", postHandler.GetPublicPosts(maxCtxDuration, dbConn))
			v1.GET("/post/search", postHandler.SearchPublicPosts(maxCtxDuration, dbConn))
			v1.GET("/post/:post", postHandler.GetPublicPost(maxCtxDuration, dbConn))
			v1.GET("/post/:post/comments", commentHandler.GetPublicPostComments(maxCtxDuration, dbConn))

			v1.GET("/comment/:comment", commentHandler.GetPublicComment(maxCtxDuration, dbConn))
			v1.GET("/comment/:comment/replies", commentHandler.GetPublicCommentReplies(maxCtxDuration, dbConn))
			v1.POST("/comment", commentHandler.CreatePublicPostComment(maxCtxDuration, dbConn))
			v1.POST("/comment/reply", commentHandler.CreatePublicCommentReply(maxCtxDuration, dbConn))

			v1.POST("/signin", authenticationHandler.SignIn(maxCtxDuration, dbConn))
			v1.POST("/signup", authenticationHandler.SignUp(maxCtxDuration, dbConn))

			refresh := v1.Group("/refresh")
			refresh.Use(authenticateMiddleware.AuthenticateRefresh(maxCtxDuration, dbConn))
			{
				refresh.POST("/signout", authenticationHandler.SignOut(maxCtxDuration, dbConn))
				refresh.POST("/", authenticationHandler.Refresh(maxCtxDuration, dbConn))
			}

			auth := v1.Group("/auth")
			auth.Use(authenticateMiddleware.Authenticate(maxCtxDuration, dbConn))
			{
				auth.GET("/me", meHandler.GetMe(maxCtxDuration, dbConn))

				verifyPassword := auth.Group("/me")
				verifyPassword.Use(authenticateMiddleware.VerifyPassword(maxCtxDuration, dbConn))
				{
					verifyPassword.PUT("/", meHandler.UpdateMe(maxCtxDuration, dbConn))
					verifyPassword.PATCH("/", meHandler.UpdateMe(maxCtxDuration, dbConn))
					verifyPassword.PUT("/password", meHandler.UpdateMePassword(maxCtxDuration, dbConn))
					verifyPassword.PATCH("/password", meHandler.UpdateMePassword(maxCtxDuration, dbConn))
				}

				auth.GET("/notifications", notificationHandler.GetNotifications(maxCtxDuration, dbConn))
				auth.GET("/notifications/listen", notificationHandler.ServeListenedNotifications(maxCtxDuration, dbConn))
				auth.GET("/notification/:notification", notificationHandler.GetNotification(maxCtxDuration, dbConn))
				auth.PUT("/notification/:notification", notificationHandler.ReadNotification(maxCtxDuration, dbConn))
				auth.PATCH("/notification/:notification", notificationHandler.ReadNotification(maxCtxDuration, dbConn))
				auth.DELETE("/notification/:notification", notificationHandler.DeleteNotification(maxCtxDuration, dbConn))

				writer := auth.Group("/writer")
				writer.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Writer"))
				{
					writer.GET("/posts", postHandler.GetMyPosts(maxCtxDuration, dbConn))
					writer.GET("/post/:post", postHandler.GetMyPost(maxCtxDuration, dbConn))
					writer.POST("/post", postHandler.CreatePost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post", postHandler.UpdateMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post", postHandler.UpdateMyPost(maxCtxDuration, dbConn))
					writer.DELETE("/post/:post", postHandler.TrashMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/detrash", postHandler.DetrashMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/detrash", postHandler.DetrashMyPost(maxCtxDuration, dbConn))
					writer.DELETE("/post/:post/permanent", postHandler.DeleteMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/publish", postHandler.PublishMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/publish", postHandler.PublishMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/depublish", postHandler.DepublishMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/depublish", postHandler.DepublishMyPost(maxCtxDuration, dbConn))
					writer.GET("/post/:post/comments", commentHandler.GetMyPostComments(maxCtxDuration, dbConn))

					writer.GET("/comments", commentHandler.GetMyComments(maxCtxDuration, dbConn))
					writer.GET("/comment/:comment", commentHandler.GetMyComment(maxCtxDuration, dbConn))
					writer.DELETE("/comment/:comment", commentHandler.TrashMyComment(maxCtxDuration, dbConn))
					writer.PUT("/comment/:comment/detrash", commentHandler.DetrashMyComment(maxCtxDuration, dbConn))
					writer.PATCH("/comment/:comment/detrash", commentHandler.DetrashMyComment(maxCtxDuration, dbConn))
					writer.DELETE("/comment/:comment/permanent", commentHandler.DeleteMyComment(maxCtxDuration, dbConn))
				}

				editor := auth.Group("/editor")
				editor.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Editor"))
				{
					editor.GET("/categories", categoryHandler.GetCategories(maxCtxDuration, dbConn))
					editor.GET("/category/:category", categoryHandler.GetCategory(maxCtxDuration, dbConn))
					editor.POST("/category", categoryHandler.CreateCategory(maxCtxDuration, dbConn))
					editor.PUT("/category/:category", categoryHandler.UpdateCategory(maxCtxDuration, dbConn))
					editor.PATCH("/category/:category", categoryHandler.UpdateCategory(maxCtxDuration, dbConn))
					editor.PUT("/category/:category/detrash", categoryHandler.DetrashCategory(maxCtxDuration, dbConn))
					editor.PATCH("/category/:category/detrash", categoryHandler.DetrashCategory(maxCtxDuration, dbConn))
					editor.DELETE("/category/:category", categoryHandler.TrashCategory(maxCtxDuration, dbConn))
					editor.DELETE("/category/:category/permanent", categoryHandler.DeleteCategory(maxCtxDuration, dbConn))

					editor.GET("/posts", postHandler.GetPosts(maxCtxDuration, dbConn))
					editor.GET("/post/:post", postHandler.GetPost(maxCtxDuration, dbConn))
					editor.PUT("/post/:post", postHandler.UpdatePost(maxCtxDuration, dbConn))
					editor.PATCH("/post/:post", postHandler.UpdatePost(maxCtxDuration, dbConn))
					editor.DELETE("/post/:post", postHandler.TrashPost(maxCtxDuration, dbConn))
					editor.DELETE("/post/:post/permanent", postHandler.DeletePost(maxCtxDuration, dbConn))
					editor.GET("/post/:post/comments", commentHandler.GetPostComments(maxCtxDuration, dbConn))

					editor.GET("/comments", commentHandler.GetComments(maxCtxDuration, dbConn))
					editor.GET("/comment/:comment", commentHandler.GetComment(maxCtxDuration, dbConn))
					editor.DELETE("/comment/:comment", commentHandler.TrashComment(maxCtxDuration, dbConn))
					editor.PUT("/comment/:comment/detrash", commentHandler.DetrashComment(maxCtxDuration, dbConn))
					editor.PATCH("/comment/:comment/detrash", commentHandler.DetrashComment(maxCtxDuration, dbConn))
					editor.DELETE("/comment/:comment/permanent", commentHandler.DeleteComment(maxCtxDuration, dbConn))
				}

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Admin"))
				{
					admin.GET("/users", userHandler.GetUsers(maxCtxDuration, dbConn))
					admin.GET("/user/:user", userHandler.GetUser(maxCtxDuration, dbConn))
					admin.POST("/user", userHandler.CreateUser(maxCtxDuration, dbConn))
					admin.PUT("/user/:user", userHandler.UpdateUser(maxCtxDuration, dbConn))
					admin.PATCH("/user/:user", userHandler.UpdateUser(maxCtxDuration, dbConn))
					admin.DELETE("/user/:user", userHandler.TrashUser(maxCtxDuration, dbConn))
					admin.PUT("/user/:user/detrash", userHandler.DetrashUser(maxCtxDuration, dbConn))
					admin.PATCH("/user/:user/detrash", userHandler.DetrashUser(maxCtxDuration, dbConn))
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "SuperAdmin"))
				{
					superadmin.PUT("/adminize/:user", userHandler.AdminizeUser(maxCtxDuration, dbConn))
					superadmin.PATCH("/adminize/:user", userHandler.AdminizeUser(maxCtxDuration, dbConn))
					superadmin.PUT("/deadminize/:user", userHandler.DeadminizeUser(maxCtxDuration, dbConn))
					superadmin.PATCH("/deadminize/:user", userHandler.DeadminizeUser(maxCtxDuration, dbConn))
					superadmin.DELETE("/user/:user/permanent", userHandler.DeleteUser(maxCtxDuration, dbConn))
				}
			}
		}
	}
}
