package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
func initRoute(server *gin.Engine) {
	maxCtxDuration := 10 * time.Second

	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			v1.GET("/users", userController.GetPublicUsers(maxCtxDuration))
			v1.GET("/user/:user", userController.GetPublicUser(maxCtxDuration))

			v1.GET("/categories", categoryController.GetPublicCategories(maxCtxDuration))
			v1.GET("/category/:category", categoryController.GetPublicCategory(maxCtxDuration))
			v1.GET("/category/:category/posts", categoryController.GetPublicCategoryPosts(maxCtxDuration))
			v1.GET("/category/:category/post/:post", categoryController.GetPublicCategoryPost(maxCtxDuration))

			v1.GET("/posts", postController.GetPublicPosts(maxCtxDuration))
			v1.GET("/post/:post", postController.GetPublicPost(maxCtxDuration))
			v1.GET("/post/:post/comments", postController.GetPublicPostComments(maxCtxDuration))
			v1.GET("/comment/:comment", postController.GetPublicPostComment(maxCtxDuration))
			v1.POST("/comment", postController.CreatePublicPostComment(maxCtxDuration))

			v1.POST("/signin", authenticationController.SignIn(maxCtxDuration))
			v1.POST("/signup", authenticationController.SignUp(maxCtxDuration))

			refresh := v1.Group("/refresh")
			refresh.Use(authenticateMiddleware.AuthenticateRefresh())
			{
				refresh.POST("/", authenticationController.Refresh(maxCtxDuration))
			}

			auth := v1.Group("/auth")
			auth.Use(authenticateMiddleware.Authenticate(maxCtxDuration))
			{
				auth.POST("/signout", authenticationController.SignOut(maxCtxDuration))

				auth.GET("/me", meController.GetMe(maxCtxDuration))

				verifyPassword := auth.Group("/me")
				verifyPassword.Use(authenticateMiddleware.VerifyPassword(maxCtxDuration))
				{
					verifyPassword.PUT("/", meController.UpdateMe(maxCtxDuration))
					verifyPassword.PATCH("/", meController.UpdateMe(maxCtxDuration))
					verifyPassword.PUT("/password", meController.UpdateMePassword(maxCtxDuration))
					verifyPassword.PATCH("/password", meController.UpdateMePassword(maxCtxDuration))
				}

				auth.GET("/notifications", notificationController.GetNotifications(maxCtxDuration))
				auth.GET("/notification/:notification", notificationController.GetNotification(maxCtxDuration))
				auth.PUT("/notification/:notification/read", notificationController.ReadNotification(maxCtxDuration))

				writer := auth.Group("/writer")
				writer.Use(authorizeMiddleware.Authorize(maxCtxDuration, "Writer"))
				{
					writer.GET("/posts", postController.GetMyPosts(maxCtxDuration))
					writer.GET("/post/:post", postController.GetMyPost(maxCtxDuration))
					writer.POST("/post", postController.CreatePost(maxCtxDuration))
					writer.PUT("/post/:post", postController.UpdateMyPost(maxCtxDuration))
					writer.PATCH("/post/:post", postController.UpdateMyPost(maxCtxDuration))
					writer.DELETE("/post/:post", postController.TrashMyPost(maxCtxDuration))
					writer.PUT("/post/:post/detrash", postController.DetrashMyPost(maxCtxDuration))
					writer.DELETE("/post/:post/permanent", postController.DeleteMyPost(maxCtxDuration))
					writer.PUT("/post/:post/publish", postController.PublishMyPost(maxCtxDuration))
					writer.PATCH("/post/:post/publish", postController.PublishMyPost(maxCtxDuration))
					writer.PUT("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration))
					writer.PATCH("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration))
					writer.GET("/post/:post/comments", postController.GetMyPostComments(maxCtxDuration))
					writer.GET("/comment/:comment", postController.GetMyPostComment(maxCtxDuration))
					writer.DELETE("/comment/:comment", postController.TrashMyPostComment(maxCtxDuration))
					writer.PUT("/comment/:comment/detrash", postController.DetrashMyPostComment(maxCtxDuration))
					writer.PATCH("/comment/:comment/detrash", postController.DetrashMyPostComment(maxCtxDuration))
					writer.DELETE("/comment/:comment/permanent", postController.DeleteMyPostComment(maxCtxDuration))
				}

				editor := auth.Group("/editor")
				editor.Use(authorizeMiddleware.Authorize(maxCtxDuration, "Editor"))
				{
					editor.GET("/categories", categoryController.GetCategories(maxCtxDuration))
					editor.GET("/category/:category", categoryController.GetCategory(maxCtxDuration))
					editor.POST("/category", categoryController.CreateCategory(maxCtxDuration))
					editor.PUT("/category/:category", categoryController.UpdateCategory(maxCtxDuration))
					editor.PATCH("/category/:category", categoryController.UpdateCategory(maxCtxDuration))
					editor.PUT("/category/:category/detrash", categoryController.DetrashCategory(maxCtxDuration))
					editor.DELETE("/category/:category", categoryController.TrashCategory(maxCtxDuration))
					editor.DELETE("/category/:category/permanent", categoryController.DeleteCategory(maxCtxDuration))

					editor.GET("/posts", postController.GetPosts(maxCtxDuration))
					editor.GET("/post/:post", postController.GetPost(maxCtxDuration))
					editor.PUT("/post/:post", postController.UpdatePost(maxCtxDuration))
					editor.PATCH("/post/:post", postController.UpdatePost(maxCtxDuration))
					editor.DELETE("/post/:post", postController.TrashPost(maxCtxDuration))
					editor.DELETE("/post/:post/permanent", postController.DeletePost(maxCtxDuration))
					editor.GET("/post/:post/comments", postController.GetPostComments(maxCtxDuration))
					editor.GET("/comment/:comment", postController.GetPostComment(maxCtxDuration))
					editor.DELETE("/comment/:comment", postController.TrashPostComment(maxCtxDuration))
					editor.PUT("/comment/:comment/detrash", postController.DetrashPostComment(maxCtxDuration))
					editor.PATCH("/comment/:comment/detrash", postController.DetrashPostComment(maxCtxDuration))
					editor.DELETE("/comment/:comment/permanent", postController.DeletePostComment(maxCtxDuration))
				}

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize(maxCtxDuration, "Admin"))
				{
					admin.GET("/users", userController.GetUsers(maxCtxDuration))
					admin.GET("/user/:user", userController.GetUser(maxCtxDuration))
					admin.POST("/user", userController.CreateUser(maxCtxDuration))
					admin.PUT("/user/:user", userController.UpdateUser(maxCtxDuration))
					admin.PATCH("/user/:user", userController.UpdateUser(maxCtxDuration))
					admin.DELETE("/user/:user", userController.TrashUser(maxCtxDuration))
					admin.DELETE("/user/:user/permanent", userController.DeleteUser(maxCtxDuration))
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize(maxCtxDuration, "SuperAdmin"))
				{
					superadmin.PUT("/adminize/:user", userController.AdminizeUser(maxCtxDuration))
					superadmin.PATCH("/adminize/:user", userController.AdminizeUser(maxCtxDuration))
					superadmin.PUT("/deadminize/:user", userController.DeadminizeUser(maxCtxDuration))
					superadmin.PATCH("/deadminize/:user", userController.DeadminizeUser(maxCtxDuration))
				}
			}
		}
	}

	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found.",
		})
	})
}
