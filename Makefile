setup:
	@docker rm mapsyncproxy -f || true
	@cd ./tools/infrastructure && docker build . -t mapsyncproxy:2.8
	@docker run --name mapsyncproxy --rm -p 8404:8404 -p 8888:8888 -p 8889:8889 -p 5555:5555 mapsyncproxy:2.8

run:
	go run main.go

push:
	gsutil cp ./tools/files/gcs.json gs://$(bucket)/gcs.json && \
	gsutil cp ./tools/files/gcs2.json gs://$(bucket)/gcs2.json

synchronize:
	curl -X POST http://localhost:8080/v1/map/$(map_name)/synchronize \
		-H 'Content-Type: application/json' \
		-d '{"bucket_name":"$(bucket)", "bucket_file_name":"$(bucket_file_name)"}'

generate:
	curl -X GET http://localhost:8080/v1/map/$(map_name)/generate

swagger:
	swag init