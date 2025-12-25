.PHONY: build clean compile serve

build:
	@mkdir -p build
	@go build -o kaimon
	@echo "✓ Built to ./kaimon"

clean:
	@rm -rf build kaimon
	@echo "✓ Cleaned"

compile: build
	@./kaimon compile

serve: build compile
	@./kaimon serve
