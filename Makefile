.PHONY: test test-short build arwen arwendebug clean

clean:
	go clean -cache -testcache

build:
	go build ./...

arwen:
	go build -o ./cmd/arwen/arwen ./cmd/arwen
	cp ./cmd/arwen/arwen ./ipc/tests

arwendebug:
ifndef ARWENDEBUG_PATH
	$(error ARWENDEBUG_PATH is undefined)
endif
	go build -o ./cmd/arwendebug/arwendebug ./cmd/arwendebug
	cp ./cmd/arwendebug/arwendebug ${ARWENDEBUG_PATH}

test: clean arwen
	go test -count=1 ./...

test-short: arwen
	go test -short -count=1 ./...

build-c-contracts:
	moapy contract build ./test/contracts/erc20
	moapy contract build ./test/contracts/counter

	moapy contract build ./test/contracts/init-correct
	moapy contract build ./test/contracts/init-simple
	moapy contract build ./test/contracts/init-wrong
	moapy contract build ./test/contracts/misc
	moapy contract build ./test/contracts/signatures
	moapy contract build ./test/contracts/kalyan3104ei
	moapy contract build ./test/contracts/breakpoint

	moapy contract build ./test/contracts/exec-same-ctx-simple-parent
	moapy contract build ./test/contracts/exec-same-ctx-simple-child
	moapy contract build ./test/contracts/exec-same-ctx-child
	moapy contract build ./test/contracts/exec-same-ctx-parent
	moapy contract build ./test/contracts/exec-dest-ctx-parent
	moapy contract build ./test/contracts/exec-dest-ctx-child
	moapy contract build ./test/contracts/exec-same-ctx-recursive
	moapy contract build ./test/contracts/exec-same-ctx-recursive-parent
	moapy contract build ./test/contracts/exec-same-ctx-recursive-child
	moapy contract build ./test/contracts/exec-dest-ctx-recursive
	moapy contract build ./test/contracts/exec-dest-ctx-recursive-parent
	moapy contract build ./test/contracts/exec-dest-ctx-recursive-child
	moapy contract build ./test/contracts/async-call-parent
	moapy contract build ./test/contracts/async-call-child
	moapy contract build ./test/contracts/exec-same-ctx-builtin


build-delegation:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-delegation-rs
	git clone --depth=1 --branch=master https://github.com/kalyan3104/dme-delegation-rs.git ${SANDBOX}/sc-delegation-rs
	rm -rf ${SANDBOX}/sc-delegation-rs/.git
	moapy contract build ${SANDBOX}/sc-delegation-rs
	moapy contract test --directory="tests" ${SANDBOX}/sc-delegation-rs
	cp ${SANDBOX}/sc-delegation-rs/output/delegation.wasm ./test/delegation/delegation.wasm


build-dns:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/dme-dns-rs
	git clone --depth=1 --branch=master https://github.com/kalyan3104/dme-dns-rs.git ${SANDBOX}/dme-dns-rs
	rm -rf ${SANDBOX}/dme-dns-rs/.git
	moapy contract build ${SANDBOX}/dme-dns-rs
	moapy contract test --directory="tests" ${SANDBOX}/dme-dns-rs
	cp ${SANDBOX}/dme-dns-rs/output/dns.wasm ./test/dns/dns.wasm


build-sc-examples:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples

	moapy contract new --template=erc20-c --directory ${SANDBOX}/sc-examples erc20-c
	moapy contract build ${SANDBOX}/sc-examples/erc20-c
	cp ${SANDBOX}/sc-examples/erc20-c/output/wrc20_arwen.wasm ./test/erc20/contracts/erc20-c.wasm


build-sc-examples-rs:
ifndef SANDBOX
	$(error SANDBOX variable is undefined)
endif
	rm -rf ${SANDBOX}/sc-examples-rs
	
	moapy contract new --template=simple-coin --directory ${SANDBOX}/sc-examples-rs simple-coin
	moapy contract new --template=adder --directory ${SANDBOX}/sc-examples-rs adder
	moapy contract build ${SANDBOX}/sc-examples-rs/adder
	moapy contract build ${SANDBOX}/sc-examples-rs/simple-coin
	moapy contract test ${SANDBOX}/sc-examples-rs/adder
	moapy contract test ${SANDBOX}/sc-examples-rs/simple-coin
	cp ${SANDBOX}/sc-examples-rs/adder/output/adder.wasm ./test/adder/adder.wasm
	cp ${SANDBOX}/sc-examples-rs/simple-coin/output/simple-coin.wasm ./test/erc20/contracts/simple-coin.wasm
