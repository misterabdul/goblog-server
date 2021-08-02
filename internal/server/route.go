package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	categoryController "github.com/misterabdul/goblog-server/internal/http/controllers/categories"
	notificationController "github.com/misterabdul/goblog-server/internal/http/controllers/notifications"
	postController "github.com/misterabdul/goblog-server/internal/http/controllers/posts"
	userController "github.com/misterabdul/goblog-server/internal/http/controllers/users"
	authenticateMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	authorizeMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authorize"
)

// Initialize all routes.
func initRoute(server *gin.Engine) {
	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})

			v1.GET("/users", userController.GetPublicUsers)
			v1.GET("/user/:user", userController.GetPublicUser)

			v1.GET("/categories", categoryController.GetPublicCategories)
			v1.GET("/category/:category", categoryController.GetPublicCategory)
			v1.GET("/category/slug/:category", categoryController.GetPublicCategorySlug)
			v1.GET("/category/:category/posts", categoryController.GetPublicCategoryPosts)
			v1.GET("/category/:category/post/:post", categoryController.GetPublicCategoryPost)

			v1.GET("/posts", postController.GetPublicPosts)
			v1.GET("/post/:post", postController.GetPublicPost)
			v1.GET("/post/slug/:post", postController.GetPublicPostSlug)
			v1.GET("/post/:post/comments", postController.GetPublicPostComments)
			v1.GET("/post/:post/comment/:comment", postController.GetPublicPostComment)
			v1.POST("/post/:post/comment", postController.CreatePublicPostComment)

			v1.POST("/signin", userController.SignIn)
			v1.POST("/signin/refresh", userController.SignInRefresh)
			v1.POST("/signup", userController.SignUp)

			auth := v1.Group("/auth")
			auth.Use(authenticateMiddleware.Authenticate())
			{
				auth.POST("/signout", userController.SignOut)

				auth.GET("/me", userController.GetMe)
				auth.PUT("/me", userController.UpdateMe)
				auth.PATCH("/me", userController.UpdateMe)
				auth.PUT("/me/password", userController.UpdateMePassword)
				auth.PATCH("/me/password", userController.UpdateMePassword)

				auth.GET("/notifications", notificationController.GetNotifications)
				auth.GET("/notification/:notification", notificationController.GetNotification)
				auth.PUT("/notification/:notification/read", notificationController.ReadNotification)

				auth.GET("/posts", postController.GetMyPosts)
				auth.GET("/post/:post", postController.GetMyPost)
				auth.POST("/post", postController.CreatePost)
				auth.PUT("/post/:post", postController.UpdateMyPost)
				auth.PATCH("/post/:post", postController.UpdateMyPost)
				auth.DELETE("/post/:post", postController.TrashMyPost)
				auth.DELETE("/post/:post/permanent", postController.DeleteMyPost)
				auth.PUT("/post/:post/publish", postController.PublishMyPost)
				auth.PATCH("/post/:post/publish", postController.PublishMyPost)
				auth.PUT("/post/:post/depublish", postController.DepublishMyPost)
				auth.PATCH("/post/:post/depublish", postController.DepublishMyPost)
				auth.GET("/post/:post/comments", postController.GetMyPostComments)
				auth.GET("/post/:post/comment/:comment", postController.GetMyPostComment)
				auth.DELETE("/post/:post/comment/:comment", postController.TrashMyPostComment)
				auth.DELETE("/post/:post/comment/:comment/permanent", postController.DeleteMyPostComment)

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize("admin"))
				{
					admin.GET("/users", userController.GetUsers)
					admin.GET("/user/:user", userController.GetUser)
					admin.POST("/user", userController.CreateUser)
					admin.PUT("/user/:user", userController.UpdateUser)
					admin.PATCH("/user/:user", userController.UpdateUser)
					admin.DELETE("/user/:user", userController.TrashUser)
					admin.DELETE("/user/:user/permanent", userController.DeleteUser)

					admin.GET("/categories", categoryController.GetCategories)
					admin.GET("/category/:category", categoryController.GetCategory)
					admin.POST("/category", categoryController.CreateCategory)
					admin.PUT("/category/:category", categoryController.UpdateCategory)
					admin.PATCH("/category/:category", categoryController.UpdateCategory)
					admin.DELETE("/category/:category", categoryController.TrashCategory)
					admin.DELETE("/category/:category/permanent", categoryController.DeleteCategory)

					admin.GET("/posts", postController.GetPosts)
					admin.GET("/post/:post", postController.GetPost)
					admin.PUT("/post/:post", postController.UpdatePost)
					admin.PATCH("/post/:post", postController.UpdatePost)
					admin.DELETE("/post/:post", postController.TrashPost)
					admin.DELETE("/post/:post/permanent", postController.DeletePost)
					admin.GET("/post/:post/comments", postController.GetPostComments)
					admin.GET("/post/:post/comment/:comment", postController.GetPostComment)
					admin.DELETE("/post/:post/comment/:comment", postController.TrashPostComment)
					admin.DELETE("/post/:post/comment/:comment/permanent", postController.DeletePostComment)
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize("superadmin"))
				{
					superadmin.PUT("/adminize/:user", userController.AdminizeUser)
					superadmin.PATCH("/adminize/:user", userController.AdminizeUser)
					superadmin.PUT("/deadminize/:user", userController.DeadminizeUser)
					superadmin.PATCH("/deadminize/:user", userController.DeadminizeUser)
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
