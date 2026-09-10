
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
run: build 
	air


# Extra: Corre los tests automatizados (haciendo todo el ciclo previo)
test: generate
	go test -v ./...

# Limpieza de archivos temporales
clean:
	rm -rf tmp