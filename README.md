# webapp-golang
Web project on Golang

# Запуск dgraph 
docker run -d -it --network web_go -p 5080:5080 -p 6080:6080 -p 8085:8080 -p 9080:9080 -p 8000:8000 -v ~/dgraph:/dgraph --name dgraph dgraph/standalone:v21.03.0

# Build web_go 
docker build -t kholmden/web_go .

# Push web_go
push kholmden/web_go:latest

# Run web_god
docker run  --name webapp --network web_go --name web_go -e DB_URL=dgraph:9080 -e ENABLE_MIGRATE=true -e DROP_FIRST=true  -v /golang/.env:/.env -p 8080:8088 kholmden/web_go


