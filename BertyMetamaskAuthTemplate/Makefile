BERTY_PORT ?= 4296

run: tidy
	go run . -berty-node-addr=127.0.0.1:$(BERTY_PORT) --debug=true 

rerun: tidy
	go run . -berty-node-addr=127.0.0.1:$(BERTY_PORT) --debug=true

start-mini-companion:
	berty daemon --node.listeners=/ip4/127.0.0.1/tcp/$(BERTY_PORT)/grpc --p2p.ipfs-api-listeners /ip4/127.0.0.1/tcp/1$(BERTY_PORT)

mini:
	berty mini

install: tidy fmt
	go install .

# utilities
tidy:; go mod tidy
fmt:;  gofumpt -w .
