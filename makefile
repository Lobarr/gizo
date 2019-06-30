build:
	go build maing.go

test:
	go test -cover -v

mock-cache:
	cd cache && moq --out ../mocks/cache.go --pkg mocks . IBigCache IJobCache

mock-core:
	cd core && moq --out ../mocks/core.go --pkg mocks . IBlockChain

mock-helpers:
	cd helpers && moq --out ../mocks/helpers.go --pkg mocks . ILogger

mock:
	make mock-cache
	make mock-core
	make mock-helpers
