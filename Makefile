setup:
	@docker rm mapsyncproxy -f || true
	@cd ./tools/infrastructure && docker build . -t mapsyncproxy:2.8
	@docker run --name mapsyncproxy --rm -p 8404:8404 -p 8888:8888 -p 8889:8889 -p 5555:5555 mapsyncproxy:2.8

push:
	gsutil cp ./tools/files/gcs.json gs://$(bucket)/gcs.json

analyze:
	curl -X POST http://localhost:8000/synchronize \
		-H 'Content-Type: application/json' \
		-d '{"map_name":"$(map_name)","bucket_name":"$(bucket)", "bucket_file_name":"gcs.json"}'