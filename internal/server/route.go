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
			v1.GET("/category/slug/:category", categoryController.GetPublicCategorySlug(maxCtxDuration))
			v1.GET("/category/:category/posts", categoryController.GetPublicCategoryPosts(maxCtxDuration))
			v1.GET("/category/:category/post/:post", categoryController.GetPublicCategoryPost(maxCtxDuration))

			v1.GET("/posts", postController.GetPublicPosts(maxCtxDuration))
			v1.GET("/post/:post", postController.GetPublicPost(maxCtxDuration))
			v1.GET("/post/slug/:post", postController.GetPublicPostSlug(maxCtxDuration))
			v1.GET("/post/:post/comments", postController.GetPublicPostComments(maxCtxDuration))
			v1.GET("/post/:post/comment/:comment", postController.GetPublicPostComment(maxCtxDuration))
			v1.POST("/post/:post/comment", postController.CreatePublicPostComment(maxCtxDuration))

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

				me := auth.Group("/me")
				me.Use(authenticateMiddleware.VerifyPassword(maxCtxDuration))
				{
					me.PUT("/", meController.UpdateMe(maxCtxDuration))
					me.PATCH("/", meController.UpdateMe(maxCtxDuration))
					me.PUT("/password", meController.UpdateMePassword(maxCtxDuration))
					me.PATCH("/password", meController.UpdateMePassword(maxCtxDuration))
				}

				auth.GET("/notifications", notificationController.GetNotifications(maxCtxDuration))
				auth.GET("/notification/:notification", notificationController.GetNotification(maxCtxDuration))
				auth.PUT("/notification/:notification/read", notificationController.ReadNotification(maxCtxDuration))

				auth.GET("/posts", postController.GetMyPosts(maxCtxDuration))
				auth.GET("/post/:post", postController.GetMyPost(maxCtxDuration))
				auth.POST("/post", postController.CreatePost(maxCtxDuration))
				auth.PUT("/post/:post", postController.UpdateMyPost(maxCtxDuration))
				auth.PATCH("/post/:post", postController.UpdateMyPost(maxCtxDuration))
				auth.DELETE("/post/:post", postController.TrashMyPost(maxCtxDuration))
				auth.DELETE("/post/:post/permanent", postController.DeleteMyPost(maxCtxDuration))
				auth.PUT("/post/:post/publish", postController.PublishMyPost(maxCtxDuration))
				auth.PATCH("/post/:post/publish", postController.PublishMyPost(maxCtxDuration))
				auth.PUT("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration))
				auth.PATCH("/post/:post/depublish", postController.DepublishMyPost(maxCtxDuration))
				auth.GET("/post/:post/comments", postController.GetMyPostComments(maxCtxDuration))
				auth.GET("/post/:post/comment/:comment", postController.GetMyPostComment(maxCtxDuration))
				auth.DELETE("/post/:post/comment/:comment", postController.TrashMyPostComment(maxCtxDuration))
				auth.DELETE("/post/:post/comment/:comment/permanent", postController.DeleteMyPostComment(maxCtxDuration))

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize("admin"))
				{
					admin.GET("/users", userController.GetUsers(maxCtxDuration))
					admin.GET("/user/:user", userController.GetUser(maxCtxDuration))
					admin.POST("/user", userController.CreateUser(maxCtxDuration))
					admin.PUT("/user/:user", userController.UpdateUser(maxCtxDuration))
					admin.PATCH("/user/:user", userController.UpdateUser(maxCtxDuration))
					admin.DELETE("/user/:user", userController.TrashUser(maxCtxDuration))
					admin.DELETE("/user/:user/permanent", userController.DeleteUser(maxCtxDuration))

					admin.GET("/categories", categoryController.GetCategories(maxCtxDuration))
					admin.GET("/category/:category", categoryController.GetCategory(maxCtxDuration))
					admin.POST("/category", categoryController.CreateCategory(maxCtxDuration))
					admin.PUT("/category/:category", categoryController.UpdateCategory(maxCtxDuration))
					admin.PATCH("/category/:category", categoryController.UpdateCategory(maxCtxDuration))
					admin.DELETE("/category/:category", categoryController.TrashCategory(maxCtxDuration))
					admin.DELETE("/category/:category/permanent", categoryController.DeleteCategory(maxCtxDuration))

					admin.GET("/posts", postController.GetPosts(maxCtxDuration))
					admin.GET("/post/:post", postController.GetPost(maxCtxDuration))
					admin.PUT("/post/:post", postController.UpdatePost(maxCtxDuration))
					admin.PATCH("/post/:post", postController.UpdatePost(maxCtxDuration))
					admin.DELETE("/post/:post", postController.TrashPost(maxCtxDuration))
					admin.DELETE("/post/:post/permanent", postController.DeletePost(maxCtxDuration))
					admin.GET("/post/:post/comments", postController.GetPostComments(maxCtxDuration))
					admin.GET("/post/:post/comment/:comment", postController.GetPostComment(maxCtxDuration))
					admin.DELETE("/post/:post/comment/:comment", postController.TrashPostComment(maxCtxDuration))
					admin.DELETE("/post/:post/comment/:comment/permanent", postController.DeletePostComment(maxCtxDuration))
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize("superadmin"))
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
