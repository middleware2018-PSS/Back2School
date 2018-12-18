# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help rundb initdb resetdb stopdb psql up down

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

COLUMNS = $(shell tput cols)
LINES = $(shell tput lines)

# DATABASE TASKS
rundb: ## Run the database container
	sudo docker run -d --rm --name back2school-db -e POSTGRES_PASSWORD=postgres --network host postgres

stopdb: ## Stop the database container
	sudo docker stop back2school-db

initdb: ## Create the database, run migrations and seed initial data
	buffalo pop create
	buffalo pop migrate
	buffalo task db:seed

resetdb: ## Drop all tables, re-run migrations and re-seed initial data
	buffalo pop reset
	buffalo pop migrate
	buffalo task db:seed

psql: ## Run psql in a separate container
	sudo docker run --rm -it --network host -e COLUMNS=$(COLUMNS) -e LINES=$(LINES) postgres psql -h localhost -U postgres

doc: ## Run swaggo to recreate the swagger documentation
	swag init -g actions/app.go

up: ## Build containters and start docker compose
	docker-compose build && docker-compose up

down: ## Stop docker compose and delete containers
	docker-compose down
