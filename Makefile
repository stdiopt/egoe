
WASMFILES=docs/main.wasm docs/index.html docs/wasm_exec.js


wasm: $(WASMFILES)

$(WASMFILES):
	GOOS=js GOARCH=wasm go build -o docs/main.wasm ./cmd/egoe-wasm
	cp assets/wasm/index.html docs/index.html
	cp assets/wasm/wasm_exec.js docs/wasm_exec.js

clean:
	rm -f $(WASMFILES)


.PHONY: wasm clean
