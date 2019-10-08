GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
APP_NAME=mflow
MAIN=main.go
VERSION := $(shell git describe --abbrev=0 --tags)
NC=\e[00m
BO=\e[01m
RD=\e[31m
GR=\e[32m
YL=\e[33m
BL=\e[34m

all:
	test build
build:
	mkdir -p release/mflow
	cp config.json release/mflow/
	cp oracle.json release/mflow/
	@$(GOBUILD) -o release/mflow/$(APP_NAME) -v
	cd release
	tar -czvf "${APP_NAME}-${VERSION}.tar.gz" ./mflow
	rm -fr mflow
	cd ..
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
