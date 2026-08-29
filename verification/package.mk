GO ?= go

.PHONY: benchmark conformance docs fuzz interoperability

benchmark:
	$(GO) test -run '^$$' -bench . -benchmem ./...

conformance:
	./scripts/check-conformance.sh

docs:
	./scripts/check-docs.sh

fuzz:
	./scripts/check-fuzz.sh 10s 4

interoperability:
	python3 scripts/check_interoperability.py
