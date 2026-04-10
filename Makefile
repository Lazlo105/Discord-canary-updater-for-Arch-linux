.PHONY: build run clean test help install

# Имя бинарного файла
BINARY_NAME=updater

# Путь к основному файлу
MAIN_FILE=main.go

# Директория для бинарника
BUILD_DIR=bin

# Флаг сборки
LDFLAGS=-ldflags "-s -w"

# По умолчанию - сборка
all: build

# Сборка проекта
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Запуск проекта
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_FILE)

# Сборка и запуск
start: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Тестирование (если есть тесты)
test:
	@echo "Running tests..."
	go test -v ./...

# Запуск тестов с покрытием
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Очистка артефактов сборки
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Установка бинарника в систему (требует sudo)
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Удаление установленной версии
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Форматирование кода
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Линтинг кода (требуется golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Проверка зависимостей
deps:
	@echo "Checking dependencies..."
	go mod verify
	go mod tidy

# Вывод справки
help:
	@echo "Available targets:"
	@echo "  build          - Build the project"
	@echo "  run            - Run the project without building"
	@echo "  start          - Build and run the project"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binary to /usr/local/bin"
	@echo "  uninstall      - Remove binary from /usr/local/bin"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  deps           - Verify and tidy dependencies"
	@echo "  help           - Show this help message"
