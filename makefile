build:
	go build maing.go

test:
	go test -cover -v

mock-cache:
	cd cache && moq --out ../mocks/cache.go --pkg mocks . IBigCache IJobCache

mock-core:
	cd core && moq --out ../mocks/core.go --pkg mocks . IBlockChain 

mock-merkletree:
	cd core/merkletree && moq --out ../../mocks/merkletree.go --pkg mocks . IMerkleNode

mock-helpers:
	cd helpers && moq --out ../mocks/helpers.go --pkg mocks . ILogger

mock:
	mkdir -p mocks
	make mock-cache
	make mock-core
	make mock-merkletree
	make mock-helpers
