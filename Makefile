init:
	rm blockchain.db || true
	go run cmd/blockchain/main.go createblockchain --addr 0x001

balance:
	go run cmd/blockchain/main.go balance --addr 0x001

