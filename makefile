
APP_NAME       = Tp2-Krohn-Trapani-Cicopiedi
CONTAINER_NAME = postgres-db

.PHONY: all run docker-up generate test build docker-down clean

# Al escribir 'make' o 'make run', ejecuta todo de principio a fin y corta solo
all: run

# 1. Levanta Docker
docker-up:
	@echo "==> 1. Levantando PostgreSQL en Docker..."
	@docker start $(CONTAINER_NAME) || docker-compose up -d
	@sleep 2

# 2. Genera código con sqlc
generate: docker-up
	@echo "==> 2. Generando código con sqlc..."
	@sqlc generate

# 3. Ejecuta las pruebas unitarias
test: generate
	@echo "==> 3. Ejecutando pruebas unitarias CRUD..."
	@go test -v ./...

# 4. Compila la aplicación
build: test
	@echo "==> 4. Compilando ejecutable..."
	@mkdir -p tmp
	@go build -o ./tmp/$(APP_NAME) .

# 5. Apaga el contenedor de Docker
docker-down:
	@echo "==> 5. Apagando contenedor de Docker..."
	@docker stop $(CONTAINER_NAME)

# Flujo principal: Ejecuta todo en orden y al final apaga Docker y cierra el proceso
run: build
	@$(MAKE) docker-down
	@echo "==> ¡Proceso finalizado con éxito!"


	---------------
	// con air
	APP_NAME       = Tp2-Krohn-Trapani-Cicopiedi
CONTAINER_NAME = postgres-db

.PHONY: all run docker-up generate build test clean

# Si solo escriben 'make' en la terminal, ejecuta 'run' y hace TODO en orden
all: run

# 1. Enciende el contenedor de Postgres en Docker
docker-up:
	docker start $(CONTAINER_NAME)

# 2. Genera el código tipado de sqlc (requiere que Docker esté prendido)
generate: docker-up
	sqlc generate

# 3. Compila el ejecutable en tmp/ (requiere que sqlc haya generado el código)
build: generate
	mkdir -p tmp
	echo $(APP_NAME) 
	go build -o ./tmp/$(APP_NAME) .

# 4. Ejecuta la aplicación cson Air (requiere la compilación previa)
run: test build
	air


# Extra: Corre los tests automatizados (haciendo todo el ciclo previo)
test: generate
	go test -v ./...

# Limpieza de archivos temporales
clean:
	@docker stop $(CONTAINER_NAME) || true
	@rm -rf tmp