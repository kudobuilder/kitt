.PHONY: kitt
kitt:
	@./scripts/build.sh $@

.PHONY: test
test:
	@./scripts/test.sh

.PHONY: lint
lint:
	@./scripts/lint.sh

.PHONY: clean
clean:
	@rm -rf bin
