.PHONY: run build test test-coverage clean

# Запуск приложения
run:
	go run cmd/api/main.go

# Сборка бинарника
build:
	mkdir -p bin
	go build -o bin/api cmd/api/main.go

# Запуск тестов
test:
	go test -v ./...

# Запуск тестов с покрытием
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Очистка
clean:
	rm -rf bin
	rm -f coverage.out coverage.html

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Форматирование кода
fmt:
	go fmt ./...

# Проверка линтером (требуется установка golangci-lint)
lint:
	golangci-lint run

# Полная проверка перед коммитом
pre-commit: fmt test
	@echo "All checks passed!"
