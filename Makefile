VERSION = 0.1.0
BINARY_NAME = "wmg"
BUILD_DIR = "./dist"

build:
		@echo "Building..."
		@go build -v -o ${BUILD_DIR}/${BINARY_NAME} -ldflags "-X main.version=${VERSION}"

install:
		@echo "Installing..."
		@go install -v -ldflags "-X main.version=${VERSION}"

clean:
		@echo "Cleaning..."
		@rm -rf $(BUILD_DIR)
