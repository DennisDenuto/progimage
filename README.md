# Description


1. run http server

example curl requests:
```
curl -X POST http://localhost:8080/api/v1/image \
-F "file=@/tmp/image" \
-H "Content-Type: multipart/form-data"
```