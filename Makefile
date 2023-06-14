# indece Monitor
# Copyright (C) 2023 indece UG (haftungsbeschränkt)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License or any
# later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

BUILD_VERSION ?= dev-$(shell git rev-parse --short HEAD)

all: frontend_dev backend_dev
release: frontend_prod backend_prod

frontend_dev: build_frontend copy_frontend
frontend_prod: build_frontend copy_frontend

backend_dev: build_backend
backend_prod: build_backend

dependencies: dependencies_frontend dependencies_backend

lint: lint_frontend lint_widget

dependencies_frontend:
	cd ./frontend && npm ci

dependencies_backend:
	cd ./backend && make --always-make dependencies

lint_frontend:
	cd ./frontend && npm run lint

build_frontend:
	cd ./frontend && BUILD_VERSION=${BUILD_VERSION} npm run build

copy_frontend:
	rm -rf ./backend/assets/www/*
	mkdir -p ./backend/assets/www
	cp -r ./frontend/build/* ./backend/assets/www/

build_backend:
	cd ./backend && BUILD_VERSION=${BUILD_VERSION} make --always-make

clean:
	rm -rf ./backend/dist
	rm -rf ./frontend/build
