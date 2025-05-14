# webapp-golang
Web project on Golang

# Запуск dgraph 
docker run -d -it --network web_go -p 5080:5080 -p 6080:6080 -p 8085:8080 -p 9080:9080 -p 8000:8000 -v ~/dgraph:/dgraph --name dgraph dgraph/standalone:v21.03.0

# Build web_go 
docker build -t kholmden/web_go .

# Push web_go
docker push kholmden/web_go:latest

# Run web_god
docker run -d  --network web_go --name web_go  -v ~/env/.env:/.env -p 8088:8080 kholmden/web_go

# good luck


