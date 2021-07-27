package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misterabdul/goblog-server/internal/controllers/categories"
	"github.com/misterabdul/goblog-server/internal/controllers/notifications"
	"github.com/misterabdul/goblog-server/internal/controllers/posts"
	"github.com/misterabdul/goblog-server/internal/controllers/users"
	"github.com/misterabdul/goblog-server/internal/middlewares/authenticate"
	"github.com/misterabdul/goblog-server/internal/middlewares/authorize"
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

			v1.GET("/users", users.GetPublicUsers)
			v1.GET("/user/:user", users.GetPublicUser)

			v1.GET("/categories", categories.GetPublicCategories)
			v1.GET("/category/:category", categories.GetPublicCategory)
			v1.GET("/category/slug/:category", categories.GetPublicCategorySlug)
			v1.GET("/category/:category/posts", categories.GetPublicCategoryPosts)
			v1.GET("/category/:category/post/:post", categories.GetPublicCategoryPost)

			v1.GET("/posts", posts.GetPublicPosts)
			v1.GET("/post/:post", posts.GetPublicPost)
			v1.GET("/post/slug/:post", posts.GetPublicPostSlug)
			v1.GET("/post/:post/comments", posts.GetPublicPostComments)
			v1.GET("/post/:post/comment/:comment", posts.GetPublicPostComment)
			v1.POST("/post/:post/comment", posts.CreatePublicPostComment)

			v1.POST("/signin", users.SignIn)
			v1.POST("/signin/refresh", users.SignInRefresh)
			v1.POST("/signup", users.SignUp)

			auth := v1.Group("/auth")
			auth.Use(authenticate.Authenticate())
			{
				auth.POST("/signout", users.SignOut)

				auth.GET("/me", users.GetMe)
				auth.PUT("/me", users.UpdateMe)
				auth.PATCH("/me", users.UpdateMe)
				auth.PUT("/me/password", users.UpdateMePassword)
				auth.PATCH("/me/password", users.UpdateMePassword)

				auth.GET("/notifications", notifications.GetNotifications)
				auth.GET("/notification/:notification", notifications.GetNotification)
				auth.PUT("/notification/:notification/read", notifications.ReadNotification)

				auth.GET("/posts", posts.GetMyPosts)
				auth.GET("/post/:post", posts.GetMyPost)
				auth.POST("/post", posts.CreatePost)
				auth.PUT("/post/:post", posts.UpdateMyPost)
				auth.PATCH("/post/:post", posts.UpdateMyPost)
				auth.DELETE("/post/:post", posts.TrashMyPost)
				auth.DELETE("/post/:post/permanent", posts.DeleteMyPost)
				auth.PUT("/post/:post/publish", posts.PublishMyPost)
				auth.PATCH("/post/:post/publish", posts.PublishMyPost)
				auth.PUT("/post/:post/depublish", posts.DepublishMyPost)
				auth.PATCH("/post/:post/depublish", posts.DepublishMyPost)
				auth.GET("/post/:post/comments", posts.GetMyPostComments)
				auth.GET("/post/:post/comment/:comment", posts.GetMyPostComment)
				auth.DELETE("/post/:post/comment/:comment", posts.TrashMyPostComment)
				auth.DELETE("/post/:post/comment/:comment/permanent", posts.DeleteMyPostComment)

				admin := auth.Group("/admin")
				admin.Use(authorize.Authorize("admin"))
				{
					admin.GET("/users", users.GetUsers)
					admin.GET("/user/:user", users.GetUser)
					admin.POST("/user", users.CreateUser)
					admin.PUT("/user/:user", users.UpdateUser)
					admin.PATCH("/user/:user", users.UpdateUser)
					admin.DELETE("/user/:user", users.TrashUser)
					admin.DELETE("/user/:user/permanent", users.DeleteUser)

					admin.GET("/categories", categories.GetCategories)
					admin.GET("/category/:category", categories.GetCategory)
					admin.POST("/category", categories.CreateCategory)
					admin.PUT("/category/:category", categories.UpdateCategory)
					admin.PATCH("/category/:category", categories.UpdateCategory)
					admin.DELETE("/category/:category", categories.TrashCategory)
					admin.DELETE("/category/:category/permanent", categories.DeleteCategory)

					admin.GET("/posts", posts.GetPosts)
					admin.GET("/post/:post", posts.GetPost)
					admin.PUT("/post/:post", posts.UpdatePost)
					admin.PATCH("/post/:post", posts.UpdatePost)
					admin.DELETE("/post/:post", posts.TrashPost)
					admin.DELETE("/post/:post/permanent", posts.DeletePost)
					admin.GET("/post/:post/comments", posts.GetPostComments)
					admin.GET("/post/:post/comment/:comment", posts.GetPostComment)
					admin.DELETE("/post/:post/comment/:comment", posts.TrashPostComment)
					admin.DELETE("/post/:post/comment/:comment/permanent", posts.DeletePostComment)
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorize.Authorize("superadmin"))
				{
					superadmin.PUT("/adminize/:user", users.AdminizeUser)
					superadmin.PATCH("/adminize/:user", users.AdminizeUser)
					superadmin.PUT("/deadminize/:user", users.DeadminizeUser)
					superadmin.PATCH("/deadminize/:user", users.DeadminizeUser)
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
