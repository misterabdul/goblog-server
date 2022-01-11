package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	authenticationController "github.com/misterabdul/goblog-server/internal/http/controllers/authentications"
	categoryController "github.com/misterabdul/goblog-server/internal/http/controllers/categories"
	meController "github.com/misterabdul/goblog-server/internal/http/controllers/me"
	notificationController "github.com/misterabdul/goblog-server/internal/http/controllers/notifications"
	postController "github.com/misterabdul/goblog-server/internal/http/controllers/posts"
	userController "github.com/misterabdul/goblog-server/internal/http/controllers/users"
	authenticateMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	authorizeMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authorize"
)

// Initialize all routes.
func initRoute(server *gin.Engine, dbConn *mongo.Database) {
	maxCtxDuration := 10 * time.Second

	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong"})
			})

			v1.GET("/users", userController.GetPublicUsers(maxCtxDuration, dbConn))
			v1.GET("/user/:user", userController.GetPublicUser(maxCtxDuration, dbConn))

			v1.GET("/categories", categoryController.GetPublicCategories(maxCtxDuration, dbConn))
			v1.GET("/category/:category", categoryController.GetPublicCategory(maxCtxDuration, dbConn))
			v1.GET("/category/:category/posts", categoryController.GetPublicCategoryPosts(maxCtxDuration, dbConn))

			v1.GET("/posts", postController.GetPublicPosts(maxCtxDuration, dbConn))
			v1.GET("/post/search", postController.SearchPublicPosts(maxCtxDuration, dbConn))
			v1.GET("/post/:post", postController.GetPublicPost(maxCtxDuration, dbConn))
			v1.GET("/post/:post/comments", postController.GetPublicPostComments(maxCtxDuration, dbConn))
			v1.GET("/comment/:comment", postController.GetPublicPostComment(maxCtxDuration, dbConn))
			v1.POST("/comment", postController.CreatePublicPostComment(maxCtxDuration, dbConn))

			v1.POST("/signin", authenticationController.SignIn(maxCtxDuration, dbConn))
			v1.POST("/signup", authenticationController.SignUp(maxCtxDuration, dbConn))

			refresh := v1.Group("/refresh")
			refresh.Use(authenticateMiddleware.AuthenticateRefresh(maxCtxDuration, dbConn))
			{
				refresh.POST("/signout", authenticationController.SignOut(maxCtxDuration, dbConn))
				refresh.POST("/", authenticationController.Refresh(maxCtxDuration, dbConn))
			}

			auth := v1.Group("/auth")
			auth.Use(authenticateMiddleware.Authenticate(maxCtxDuration, dbConn))
			{
				auth.GET("/me", meController.GetMe(maxCtxDuration, dbConn))

				verifyPassword := auth.Group("/me")
				verifyPassword.Use(authenticateMiddleware.VerifyPassword(maxCtxDuration, dbConn))
				{
					verifyPassword.PUT("/", meController.UpdateMe(maxCtxDuration, dbConn))
					verifyPassword.PATCH("/", meController.UpdateMe(maxCtxDuration, dbConn))
					verifyPassword.PUT("/password", meController.UpdateMePassword(maxCtxDuration, dbConn))
					verifyPassword.PATCH("/password", meController.UpdateMePassword(maxCtxDuration, dbConn))
				}

				auth.GET("/notifications", notificationController.GetNotifications(maxCtxDuration, dbConn))
				auth.GET("/notifications/listen", notificationController.ServeListenedNotifications(maxCtxDuration, dbConn))
				auth.GET("/notification/:notification", notificationController.GetNotification(maxCtxDuration, dbConn))
				auth.PUT("/notification/:notification", notificationController.ReadNotification(maxCtxDuration, dbConn))
				auth.PATCH("/notification/:notification", notificationController.ReadNotification(maxCtxDuration, dbConn))
				auth.DELETE("/notification/:notification", notificationController.DeleteNotification(maxCtxDuration, dbConn))

				writer := auth.Group("/writer")
				writer.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Writer"))
				{
					writer.GET("/posts", postController.GetMyPosts(maxCtxDuration, dbConn))
					writer.GET("/post/:post", postController.GetMyPost(maxCtxDuration, dbConn))
					writer.POST("/post", postController.CreatePost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post", postController.UpdateMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post", postController.UpdateMyPost(maxCtxDuration, dbConn))
					writer.DELETE("/post/:post", postController.TrashMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/detrash", postController.DetrashMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/detrash", postController.DetrashMyPost(maxCtxDuration, dbConn))
					writer.DELETE("/post/:post/permanent", postController.DeleteMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/publish", postController.PublishMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/publish", postController.PublishMyPost(maxCtxDuration, dbConn))
					writer.PUT("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration, dbConn))
					writer.PATCH("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration, dbConn))
					writer.GET("/post/:post/comments", postController.GetMyPostComments(maxCtxDuration, dbConn))
					writer.GET("/comment/:comment", postController.GetMyPostComment(maxCtxDuration, dbConn))
					writer.DELETE("/comment/:comment", postController.TrashMyPostComment(maxCtxDuration, dbConn))
					writer.PUT("/comment/:comment/detrash", postController.DetrashMyPostComment(maxCtxDuration, dbConn))
					writer.PATCH("/comment/:comment/detrash", postController.DetrashMyPostComment(maxCtxDuration, dbConn))
					writer.DELETE("/comment/:comment/permanent", postController.DeleteMyPostComment(maxCtxDuration, dbConn))
				}

				editor := auth.Group("/editor")
				editor.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Editor"))
				{
					editor.GET("/categories", categoryController.GetCategories(maxCtxDuration, dbConn))
					editor.GET("/category/:category", categoryController.GetCategory(maxCtxDuration, dbConn))
					editor.POST("/category", categoryController.CreateCategory(maxCtxDuration, dbConn))
					editor.PUT("/category/:category", categoryController.UpdateCategory(maxCtxDuration, dbConn))
					editor.PATCH("/category/:category", categoryController.UpdateCategory(maxCtxDuration, dbConn))
					editor.PUT("/category/:category/detrash", categoryController.DetrashCategory(maxCtxDuration, dbConn))
					editor.PATCH("/category/:category/detrash", categoryController.DetrashCategory(maxCtxDuration, dbConn))
					editor.DELETE("/category/:category", categoryController.TrashCategory(maxCtxDuration, dbConn))
					editor.DELETE("/category/:category/permanent", categoryController.DeleteCategory(maxCtxDuration, dbConn))

					editor.GET("/posts", postController.GetPosts(maxCtxDuration, dbConn))
					editor.GET("/post/:post", postController.GetPost(maxCtxDuration, dbConn))
					editor.PUT("/post/:post", postController.UpdatePost(maxCtxDuration, dbConn))
					editor.PATCH("/post/:post", postController.UpdatePost(maxCtxDuration, dbConn))
					editor.DELETE("/post/:post", postController.TrashPost(maxCtxDuration, dbConn))
					editor.DELETE("/post/:post/permanent", postController.DeletePost(maxCtxDuration, dbConn))
					editor.GET("/post/:post/comments", postController.GetPostComments(maxCtxDuration, dbConn))
					editor.GET("/comment/:comment", postController.GetPostComment(maxCtxDuration, dbConn))
					editor.DELETE("/comment/:comment", postController.TrashPostComment(maxCtxDuration, dbConn))
					editor.PUT("/comment/:comment/detrash", postController.DetrashPostComment(maxCtxDuration, dbConn))
					editor.PATCH("/comment/:comment/detrash", postController.DetrashPostComment(maxCtxDuration, dbConn))
					editor.DELETE("/comment/:comment/permanent", postController.DeletePostComment(maxCtxDuration, dbConn))
				}

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "Admin"))
				{
					admin.GET("/users", userController.GetUsers(maxCtxDuration, dbConn))
					admin.GET("/user/:user", userController.GetUser(maxCtxDuration, dbConn))
					admin.POST("/user", userController.CreateUser(maxCtxDuration, dbConn))
					admin.PUT("/user/:user", userController.UpdateUser(maxCtxDuration, dbConn))
					admin.PATCH("/user/:user", userController.UpdateUser(maxCtxDuration, dbConn))
					admin.DELETE("/user/:user", userController.TrashUser(maxCtxDuration, dbConn))
					admin.DELETE("/user/:user/permanent", userController.DeleteUser(maxCtxDuration, dbConn))
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize(maxCtxDuration, dbConn, "SuperAdmin"))
				{
					superadmin.PUT("/adminize/:user", userController.AdminizeUser(maxCtxDuration, dbConn))
					superadmin.PATCH("/adminize/:user", userController.AdminizeUser(maxCtxDuration, dbConn))
					superadmin.PUT("/deadminize/:user", userController.DeadminizeUser(maxCtxDuration, dbConn))
					superadmin.PATCH("/deadminize/:user", userController.DeadminizeUser(maxCtxDuration, dbConn))
				}
			}
		}
	}

	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found."})
	})
}
