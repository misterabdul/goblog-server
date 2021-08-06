/**
 * So many times wasted on this code.
 * I've been stuck with this stupid binding bugs by Gin Gonic.
 * I have to search manually for the error. Let alone how to solve it.
 * The conclusion is:
 *   - ShouldBind method is really bad and cause stupid unnecessary problems for JSON request
 *   - Always use ShouldBindBodyWith method with necessary request type
 *
 * Why JSON request bindings can't invoked multiple times via ShouldBind method ?
 */

package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func shouldBind(c *gin.Context, obj interface{}) error {
	contentType := c.ContentType()
	switch {
	case contentType == "application/json":
		return c.ShouldBindBodyWith(obj, binding.JSON)
	case contentType == "application/x-msgpack":
		fallthrough
	case contentType == "application/msgpack":
		return c.ShouldBindBodyWith(obj, binding.MsgPack)
	default:
		return c.ShouldBind(obj)
	}
}
