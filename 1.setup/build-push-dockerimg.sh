docker buildx build --platform linux/amd64,linux/arm64 -t tomtris02/paas_image:testing  --push .
docker buildx build --platform linux/amd64,linux/arm64 -t tomtris02/go-app-backend:v1 --push .
docker buildx build --platform linux/amd64,linux/arm64 -t tomtris02/go-app-frontend:v2 --push .