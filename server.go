package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"muraldevice/artifact"
	"muraldevice/mural"

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
	"github.com/spf13/afero"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	trial       = 10
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

func main() {
	fs := afero.NewOsFs()
	softJSON, err := fs.Open("containerFiles/software.json")
	if err != nil {
		log.Panic(err)
	}
	muralHandler := mural.New(softJSON)

	fmt.Println("will try to connect to db")
	client := artifact.Dbdriver()
	arh := artifact.NewARH(client)
	artifactHandler := artifact.New(fs, arh)

	http.HandleFunc("/artifact", artifactHandler.HandleArtifacts)
	http.HandleFunc("/muralInfo", muralHandler.GetSoftwareSummary)

	http.Handle("/image", &imageserver_http.Handler{
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
	})
	http.HandleFunc("/", getterPoster)
	//clean up connection on application shutdown
	defer closeCons(context.TODO(), *client)
	cleanupOnExit(context.TODO(), *client)

	fmt.Println("Done with setup")
	http.ListenAndServe(":42069", nil)
}

func getterPoster(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "artifact.html")
}

func cleanupOnExit(context context.Context, client mongo.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		closeCons(context, client)
		os.Exit(1)
	}()

}

func closeCons(context context.Context, client mongo.Client) {
	client.Disconnect(context)
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
