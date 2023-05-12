# Description

### Quick start
1. run http server `go run main.go -basePath=/tmp/ -port=8080`
1. Uploading an image:
example curl request:
```
curl -X POST http://localhost:8080/api/v1/image \
-F "file=@/tmp/image" \
-H "Content-Type: multipart/form-data"
```
1. Downloading an image:
example curl request:
```
curl -v curl -X GET http://127.0.0.1:8080/api/v1/image/1  
```


### Run tests:
1. ./hack/test.sh