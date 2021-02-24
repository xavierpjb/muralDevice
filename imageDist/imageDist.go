package imageDist

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_gift "github.com/pierrre/imageserver/http/gift"
	imageserver_http_image "github.com/pierrre/imageserver/http/image"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/gif"
	imageserver_image_gift "github.com/pierrre/imageserver/image/gift"
	_ "github.com/pierrre/imageserver/image/jpeg"
	_ "github.com/pierrre/imageserver/image/png"
	imageserver_source "github.com/pierrre/imageserver/source"
)

var (
	ServerLocal = imageserver.Server(imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		source, err := params.GetString(imageserver_source.Param)
		if err != nil {
			return nil, err
		}
		im, err := Get(source)
		if err != nil {
			return nil, &imageserver.ParamError{Param: imageserver_source.Param, Message: err.Error()}
		}
		return im, nil
	}))
)

func ImageDistributor() *imageserver_http.Handler {

	dist := imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			&imageserver_http.SourceParser{},
			&imageserver_http_gift.ResizeParser{},
			&imageserver_http_image.FormatParser{},
			&imageserver_http_image.QualityParser{},
		}),
		Server: &imageserver.HandlerServer{
			Server: ServerLocal,
			Handler: &imageserver_image.Handler{
				Processor: &imageserver_image_gift.ResizeProcessor{},
			},
		},
	}
	return &dist
}

func Get(name string) (*imageserver.Image, error) {
	return loadImage(name, "jpeg"), nil
}
func loadImage(filename string, format string) *imageserver.Image {
	filePath := filepath.Join("containerFiles/artifacts", filename)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	im := &imageserver.Image{
		Format: format,
		Data:   data,
	}
	return im
}
