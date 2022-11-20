demo: init balance payto

init:
	rm blockchain.db || true
	go run cmd/blockchain/main.go create-chain --addr 0x001

balance:
	go run cmd/blockchain/main.go balance --addr 0x001

payto:
	go run cmd/blockchain/main.go payto --from 0x001 --to 0x002 --amount 10

