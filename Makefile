GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
APP_NAME=mflow
MAIN=main.go
VERSION := $(shell git describe --tags)
NC=\e[00m
BO=\e[01m
RD=\e[31m
GR=\e[32m
YL=\e[33m
BL=\e[34m

all:
	test build
build: 
	@$(GOBUILD) -o $(APP_NAME) -v
test: 
	@$(GOTEST) -v -cover
clean: 
	@$(GOCLEAN)
	@rm -f $(APP_NAME)
run:
	@$(GORUN) $(MAIN)
deps:
	@echo -e "${GR}Obteniendo dependencias${NC}"
	@$(GOMOD) vendor
	@$(GOMOD) download
	@echo -e "${GR}Limpiando lo que no es necesario${NC}"
	@$(GOMOD) verify
	@$(GOMOD) tidy
	@echo -e "${YL}Dependencias:${NC}"
	@$(GOMOD) graph
	@echo -e "${YL}Finalizado${NC}"