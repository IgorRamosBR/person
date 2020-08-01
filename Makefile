run-test:
	go test ./... -v

run-local:
	cd build/package; docker-compose up -d --build

stop-local:
	cd build/package; docker-compose down