$goPath = "/go/src/app/vendor/github.com/FINTprosjektet/fint-model"
docker run -v ${PWD}:${goPath} -w $goPath -e GOOS=windows golang go build