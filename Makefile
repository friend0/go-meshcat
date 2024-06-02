# Makefile

# Targets
.PHONY: install up build_app

install:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	npx wabpack server

up:
	air
	pnx webpaxk serve


build_app:
	/bin/bash web/meshcat/rebuild.sh