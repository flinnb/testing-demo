package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/disintegration/gift"
	"github.com/gin-gonic/gin"
)

func postImage(c *gin.Context) {
	imgFile, err := c.FormFile("image")
	if err != nil {
		c.Error(err)
		return
	}
	src, err := imgFile.Open()
	if err != nil {
		c.Error(err)
		return
	}
	jpgIn, err := jpeg.Decode(src)
	if err != nil {
		c.Error(err)
		return
	}
	filter := gift.ResizeToFit(256, 256, gift.LanczosResampling)
	g := gift.New(filter)
	pngOut := image.NewRGBA(g.Bounds(jpgIn.Bounds()))
	g.Draw(pngOut, jpgIn)

	var b bytes.Buffer
	err = png.Encode(&b, pngOut)
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, "image/png", b.Bytes())
}

func RegisterHandlers(group *gin.RouterGroup) {
	images := group.Group("images")
	images.POST("", postImage)
}
