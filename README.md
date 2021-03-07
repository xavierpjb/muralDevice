# muralDevice
A go server which handles upload and serves images
https://blog.seriesci.com/how-to-measure-code-coverage-in-go/ for info on how to test
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
To build container run 
docker build -t mvral .
in same directory as dockerfile


To generate mock files
mockgen -source=<source> -destination=<destination>

To tar file 
docker save waduphaitian/mural_dev:latest | gzip > mural_dev.tar.gz