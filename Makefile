# Complete Makefile for mbservices deployment and service management

# Application variables
APP_NAME := mbservices
GITHUB_REPO := https://github.com/nsvirk/mbservices.git
GITHUB_BRANCH := main
APP_DIR := /home/ec2-user/apps/mbservices

# Service variables
SERVICE_NAME := $(APP_NAME)
SERVICE_DESCRIPTION := mbservices Go Application
SERVICE_EXEC := $(APP_DIR)/$(APP_NAME)
SERVICE_USER := ec2-user
SERVICE_GROUP := ec2-user
SERVICE_FILE := /etc/systemd/system/$(SERVICE_NAME).service

.PHONY: all deploy update build create-service install enable start stop restart status logs

all: deploy

deploy: update build create-service install enable start

update:
	@echo "Updating or cloning mbservices from GitHub..."
	@if [ -d $(APP_DIR)/.git ]; then \
		cd $(APP_DIR) && git pull origin $(GITHUB_BRANCH); \
	else \
		rm -rf $(APP_DIR) && \
		git clone -b $(GITHUB_BRANCH) $(GITHUB_REPO) $(APP_DIR); \
	fi

build:
	@echo "Building mbservices application..."
	@cd $(APP_DIR) && go build -o $(APP_NAME)

create-service:
	@echo "Creating systemd service file..."
	@echo "[Unit]" > $(SERVICE_NAME).service
	@echo "Description=$(SERVICE_DESCRIPTION)" >> $(SERVICE_NAME).service
	@echo "" >> $(SERVICE_NAME).service
	@echo "[Service]" >> $(SERVICE_NAME).service
	@echo "ExecStart=$(SERVICE_EXEC)" >> $(SERVICE_NAME).service
	@echo "User=$(SERVICE_USER)" >> $(SERVICE_NAME).service
	@echo "Group=$(SERVICE_GROUP)" >> $(SERVICE_NAME).service
	@echo "Restart=always" >> $(SERVICE_NAME).service
	@echo "" >> $(SERVICE_NAME).service
	@echo "[Install]" >> $(SERVICE_NAME).service
	@echo "WantedBy=multi-user.target" >> $(SERVICE_NAME).service
	@echo "Systemd service file created: $(SERVICE_NAME).service"

install:
	@echo "Installing mbservices service file..."
	@sudo mv $(SERVICE_NAME).service $(SERVICE_FILE)
	@sudo systemctl daemon-reload

enable:
	@echo "Enabling mbservices to start on boot..."
	@sudo systemctl enable $(SERVICE_NAME)

start:
	@echo "Starting mbservices..."
	@sudo systemctl start $(SERVICE_NAME)

stop:
	@echo "Stopping mbservices..."
	@sudo systemctl stop $(SERVICE_NAME) || true

restart: stop start

status:
	@echo "Checking mbservices status..."
	@sudo systemctl status $(SERVICE_NAME)

logs:
	@echo "Showing mbservices logs..."
	@sudo journalctl -u $(SERVICE_NAME) -f