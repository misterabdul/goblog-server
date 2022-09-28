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
	pageHandler "github.com/misterabdul/goblog-server/internal/http/handlers/pages"
	postHandler "github.com/misterabdul/goblog-server/internal/http/handlers/posts"
	userHandler "github.com/misterabdul/goblog-server/internal/http/handlers/users"
	authenticateMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authenticate"
	authorizeMiddleware "github.com/misterabdul/goblog-server/internal/http/middlewares/authorize"
	"github.com/misterabdul/goblog-server/internal/queue/client"
	"github.com/misterabdul/goblog-server/internal/service"
)

// Initialize all routes.
func InitRoutes(
	server *gin.Engine,
	dbConn *mongo.Database,
	queueClient *client.QueueClient,
	maxCtxDuration time.Duration,
) {
	svc := service.NewService(dbConn, queueClient)

	server.NoRoute(otherHandler.NotFound())

	api := server.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", otherHandler.Ping())

			v1.GET("/users", userHandler.GetPublicUsers(maxCtxDuration, svc))
			v1.GET("/user/:user", userHandler.GetPublicUser(maxCtxDuration, svc))

			v1.GET("/categories", categoryHandler.GetPublicCategories(maxCtxDuration, svc))
			v1.GET("/category/:category", categoryHandler.GetPublicCategory(maxCtxDuration, svc))
			v1.GET("/category/:category/posts", categoryHandler.GetPublicCategoryPosts(maxCtxDuration, svc))

			v1.GET("/posts", postHandler.GetPublicPosts(maxCtxDuration, svc))
			v1.GET("/post/search", postHandler.SearchPublicPosts(maxCtxDuration, svc))
			v1.GET("/post/:post", postHandler.GetPublicPost(maxCtxDuration, svc))
			v1.GET("/post/:post/comments", commentHandler.GetPublicPostComments(maxCtxDuration, svc))

			v1.GET("/pages", pageHandler.GetPublicPages(maxCtxDuration, svc))
			v1.GET("/page/search", pageHandler.SearchPublicPages(maxCtxDuration, svc))
			v1.GET("/page/:page", pageHandler.GetPublicPage(maxCtxDuration, svc))
			v1.GET("/page/slug", pageHandler.GetPublicPageBySlug(maxCtxDuration, svc))

			v1.GET("/comment/:comment", commentHandler.GetPublicComment(maxCtxDuration, svc))
			v1.GET("/comment/:comment/replies", commentHandler.GetPublicCommentReplies(maxCtxDuration, svc))
			v1.POST("/comment", commentHandler.CreatePublicPostComment(maxCtxDuration, svc))
			v1.POST("/comment/reply", commentHandler.CreatePublicCommentReply(maxCtxDuration, svc))

			v1.POST("/signin", authenticationHandler.SignIn(maxCtxDuration, svc))
			v1.POST("/signup", authenticationHandler.SignUp(maxCtxDuration, svc))

			refresh := v1.Group("/refresh")
			refresh.Use(authenticateMiddleware.AuthenticateRefresh(maxCtxDuration, svc))
			{
				refresh.POST("/signout", authenticationHandler.SignOut(maxCtxDuration, svc))
				refresh.POST("", authenticationHandler.Refresh(maxCtxDuration, svc))
			}

			auth := v1.Group("/auth")
			auth.Use(authenticateMiddleware.Authenticate(maxCtxDuration, svc))
			{
				auth.GET("/me", meHandler.GetMe(maxCtxDuration, svc))

				verifyPassword := auth.Group("/me")
				verifyPassword.Use(authenticateMiddleware.VerifyPassword(maxCtxDuration, svc))
				{
					verifyPassword.PUT("/", meHandler.UpdateMe(maxCtxDuration, svc))
					verifyPassword.PATCH("/", meHandler.UpdateMe(maxCtxDuration, svc))
					verifyPassword.PUT("/password", meHandler.UpdateMePassword(maxCtxDuration, svc))
					verifyPassword.PATCH("/password", meHandler.UpdateMePassword(maxCtxDuration, svc))
				}

				auth.GET("/notifications", notificationHandler.GetNotifications(maxCtxDuration, svc))
				auth.GET("/notifications/listen", notificationHandler.ServeListenedNotifications(maxCtxDuration, svc))
				auth.GET("/notification/:notification", notificationHandler.GetNotification(maxCtxDuration, svc))
				auth.PUT("/notification/:notification", notificationHandler.ReadNotification(maxCtxDuration, svc))
				auth.PATCH("/notification/:notification", notificationHandler.ReadNotification(maxCtxDuration, svc))
				auth.DELETE("/notification/:notification", notificationHandler.DeleteNotification(maxCtxDuration, svc))

				writer := auth.Group("/writer")
				writer.Use(authorizeMiddleware.Authorize(maxCtxDuration, svc, "Writer"))
				{
					writer.GET("/posts", postHandler.GetMyPosts(maxCtxDuration, svc))
					writer.GET("/posts/stats", postHandler.GetMyPostsStats(maxCtxDuration, svc))
					writer.GET("/post/:post", postHandler.GetMyPost(maxCtxDuration, svc))
					writer.POST("/post", postHandler.CreatePost(maxCtxDuration, svc))
					writer.PUT("/post/:post", postHandler.UpdateMyPost(maxCtxDuration, svc))
					writer.PATCH("/post/:post", postHandler.UpdateMyPost(maxCtxDuration, svc))
					writer.DELETE("/post/:post", postHandler.TrashMyPost(maxCtxDuration, svc))
					writer.PUT("/post/:post/detrash", postHandler.DetrashMyPost(maxCtxDuration, svc))
					writer.PATCH("/post/:post/detrash", postHandler.DetrashMyPost(maxCtxDuration, svc))
					writer.DELETE("/post/:post/permanent", postHandler.DeleteMyPost(maxCtxDuration, svc))
					writer.PUT("/post/:post/publish", postHandler.PublishMyPost(maxCtxDuration, svc))
					writer.PATCH("/post/:post/publish", postHandler.PublishMyPost(maxCtxDuration, svc))
					writer.PUT("/post/:post/depublish", postHandler.DepublishMyPost(maxCtxDuration, svc))
					writer.PATCH("/post/:post/depublish", postHandler.DepublishMyPost(maxCtxDuration, svc))
					writer.GET("/post/:post/comments", commentHandler.GetMyPostComments(maxCtxDuration, svc))
					writer.GET("/post/:post/comments/stats", commentHandler.GetMyPostCommentsStats(maxCtxDuration, svc))

					writer.GET("/comments", commentHandler.GetMyComments(maxCtxDuration, svc))
					writer.GET("/comments/stats", commentHandler.GetMyCommentsStats(maxCtxDuration, svc))
					writer.GET("/comment/:comment", commentHandler.GetMyComment(maxCtxDuration, svc))
					writer.DELETE("/comment/:comment", commentHandler.TrashMyComment(maxCtxDuration, svc))
					writer.PUT("/comment/:comment/detrash", commentHandler.DetrashMyComment(maxCtxDuration, svc))
					writer.PATCH("/comment/:comment/detrash", commentHandler.DetrashMyComment(maxCtxDuration, svc))
					writer.DELETE("/comment/:comment/permanent", commentHandler.DeleteMyComment(maxCtxDuration, svc))
				}

				editor := auth.Group("/editor")
				editor.Use(authorizeMiddleware.Authorize(maxCtxDuration, svc, "Editor"))
				{
					editor.GET("/categories", categoryHandler.GetCategories(maxCtxDuration, svc))
					editor.GET("/categories/stats", categoryHandler.GetCategoriesStats(maxCtxDuration, svc))
					editor.GET("/category/:category", categoryHandler.GetCategory(maxCtxDuration, svc))
					editor.POST("/category", categoryHandler.CreateCategory(maxCtxDuration, svc))
					editor.PUT("/category/:category", categoryHandler.UpdateCategory(maxCtxDuration, svc))
					editor.PATCH("/category/:category", categoryHandler.UpdateCategory(maxCtxDuration, svc))
					editor.PUT("/category/:category/detrash", categoryHandler.DetrashCategory(maxCtxDuration, svc))
					editor.PATCH("/category/:category/detrash", categoryHandler.DetrashCategory(maxCtxDuration, svc))
					editor.DELETE("/category/:category", categoryHandler.TrashCategory(maxCtxDuration, svc))
					editor.DELETE("/category/:category/permanent", categoryHandler.DeleteCategory(maxCtxDuration, svc))

					editor.GET("/posts", postHandler.GetPosts(maxCtxDuration, svc))
					editor.GET("/posts/stats", postHandler.GetPostsStats(maxCtxDuration, svc))
					editor.GET("/post/:post", postHandler.GetPost(maxCtxDuration, svc))
					editor.POST("/post", postHandler.CreatePost(maxCtxDuration, svc))
					editor.PUT("/post/:post", postHandler.UpdatePost(maxCtxDuration, svc))
					editor.PATCH("/post/:post", postHandler.UpdatePost(maxCtxDuration, svc))
					editor.DELETE("/post/:post", postHandler.TrashPost(maxCtxDuration, svc))
					editor.DELETE("/post/:post/permanent", postHandler.DeletePost(maxCtxDuration, svc))
					editor.PUT("/post/:post/publish", postHandler.PublishPost(maxCtxDuration, svc))
					editor.PATCH("/post/:post/publish", postHandler.PublishPost(maxCtxDuration, svc))
					editor.PUT("/post/:post/depublish", postHandler.DepublishPost(maxCtxDuration, svc))
					editor.PATCH("/post/:post/depublish", postHandler.DepublishPost(maxCtxDuration, svc))
					editor.GET("/post/:post/comments", commentHandler.GetPostComments(maxCtxDuration, svc))
					editor.GET("/post/:post/comments/stats", commentHandler.GetPostCommentsStats(maxCtxDuration, svc))

					editor.GET("/comments", commentHandler.GetComments(maxCtxDuration, svc))
					editor.GET("/comments/stats", commentHandler.GetCommentsStats(maxCtxDuration, svc))
					editor.GET("/comment/:comment", commentHandler.GetComment(maxCtxDuration, svc))
					editor.DELETE("/comment/:comment", commentHandler.TrashComment(maxCtxDuration, svc))
					editor.PUT("/comment/:comment/detrash", commentHandler.DetrashComment(maxCtxDuration, svc))
					editor.PATCH("/comment/:comment/detrash", commentHandler.DetrashComment(maxCtxDuration, svc))
					editor.DELETE("/comment/:comment/permanent", commentHandler.DeleteComment(maxCtxDuration, svc))

					editor.GET("/pages", pageHandler.GetPages(maxCtxDuration, svc))
					editor.GET("/pages/stats", pageHandler.GetPagesStats(maxCtxDuration, svc))
					editor.GET("/page/:page", pageHandler.GetPage(maxCtxDuration, svc))
					editor.POST("/page", pageHandler.CreatePage(maxCtxDuration, svc))
					editor.PUT("/page/:page", pageHandler.UpdatePage(maxCtxDuration, svc))
					editor.PATCH("/page/:page", pageHandler.UpdatePage(maxCtxDuration, svc))
					editor.DELETE("/page/:page", pageHandler.TrashPage(maxCtxDuration, svc))
					editor.PUT("/page/:page/publish", pageHandler.PublishPage(maxCtxDuration, svc))
					editor.PATCH("/page/:page/publish", pageHandler.PublishPage(maxCtxDuration, svc))
					editor.PUT("/page/:page/depublish", pageHandler.DepublishPage(maxCtxDuration, svc))
					editor.PATCH("/page/:page/depublish", pageHandler.DepublishPage(maxCtxDuration, svc))
					editor.PUT("/page/:page/detrash", pageHandler.DetrashPage(maxCtxDuration, svc))
					editor.PATCH("/page/:page/detrash", pageHandler.DetrashPage(maxCtxDuration, svc))
					editor.DELETE("/page/:page/permanent", pageHandler.DeletePage(maxCtxDuration, svc))
				}

				admin := auth.Group("/admin")
				admin.Use(authorizeMiddleware.Authorize(maxCtxDuration, svc, "Admin"))
				{
					admin.GET("/users", userHandler.GetUsers(maxCtxDuration, svc))
					admin.GET("/users/stats", userHandler.GetUsersStats(maxCtxDuration, svc))
					admin.GET("/user/:user", userHandler.GetUser(maxCtxDuration, svc))
					admin.POST("/user", userHandler.CreateUser(maxCtxDuration, svc))
					admin.PUT("/user/:user", userHandler.UpdateUser(maxCtxDuration, svc))
					admin.PATCH("/user/:user", userHandler.UpdateUser(maxCtxDuration, svc))
					admin.DELETE("/user/:user", userHandler.TrashUser(maxCtxDuration, svc))
					admin.PUT("/user/:user/detrash", userHandler.DetrashUser(maxCtxDuration, svc))
					admin.PATCH("/user/:user/detrash", userHandler.DetrashUser(maxCtxDuration, svc))
				}

				superadmin := auth.Group("/superadmin")
				superadmin.Use(authorizeMiddleware.Authorize(maxCtxDuration, svc, "SuperAdmin"))
				{
					superadmin.PUT("/adminize/:user", userHandler.AdminizeUser(maxCtxDuration, svc))
					superadmin.PATCH("/adminize/:user", userHandler.AdminizeUser(maxCtxDuration, svc))
					superadmin.PUT("/deadminize/:user", userHandler.DeadminizeUser(maxCtxDuration, svc))
					superadmin.PATCH("/deadminize/:user", userHandler.DeadminizeUser(maxCtxDuration, svc))
					superadmin.DELETE("/user/:user/permanent", userHandler.DeleteUser(maxCtxDuration, svc))
				}
			}
		}
	}
}
