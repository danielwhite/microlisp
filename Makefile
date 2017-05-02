# Each suite is the name of a Lisp file contained in testdata/ that
# should be loaded as part of system tests.
SUITES=eval

BINARY=microlisp
TEST_BINARY=$(BINARY).test

SOURCES=$(shell find $(CURDIR) -name '*.go')

COVER_DIR=_coverage

.PHONY: all
all: test $(BINARY)

.PHONY: test test-unit test-system

test: test-unit test-system

test-unit:
	go test ./...

test-system: $(addprefix $(COVER_DIR)/,$(addsuffix .html,$(SUITES)))

.PHONY: clean
clean:
	rm -rf $(BINARY) $(TEST_BINARY) $(COVER_DIR)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)

$(TEST_BINARY): $(SOURCES)
	go test -c -o $@ -covermode=count -coverpkg ./...

$(COVER_DIR)/%.cov: testdata/%.lisp $(TEST_BINARY)
	mkdir -p $(@D)
	./microlisp.test -test.system -test.coverprofile $@ -test.run Main -load $<

$(COVER_DIR)/%.html: $(COVER_DIR)/%.cov
	go tool cover -html=$< -o $@
