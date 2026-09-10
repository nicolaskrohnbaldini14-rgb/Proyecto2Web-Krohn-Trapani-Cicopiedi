package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUsuario_CRUD(t *testing.T) {
	// 1. Conexión a PostgreSQL en Docker
	connStr := "postgres://postgres:postgres@localhost:5432/mi_tp2_db?sslmode=disable"
	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		t.Fatalf("Base de datos inalcanzable. ¿Está corriendo Docker?: %v", err)
	}

	queries := New(dbConn)
	ctx := context.Background()

	// 2. Limpiar la tabla usuario antes de probar
	_, err = dbConn.Exec("TRUNCATE TABLE usuario CASCADE")
	if err != nil {
		t.Fatalf("Error al limpiar la tabla usuario: %v", err)
	}

	var createdID int32

	// --- A. CREAR (Create) ---
	t.Run("CreateUsuario", func(t *testing.T) {
		usr, err := queries.CreateUsuario(ctx, CreateUsuarioParams{
			Nombre:   "Nicolas",
			Apellido: "Krohn",
			Dni:      "40123456",
			Email:    "nicolas@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("Error al crear usuario: %v", err)
		}
		if usr.IDUsuario == 0 {
			t.Error("Se esperaba un ID de usuario mayor a 0")
		}
		createdID = usr.IDUsuario
	})

	// --- B. LEER / OBTENER (Read) ---
	t.Run("GetUsuario", func(t *testing.T) {
		usr, err := queries.GetUsuario(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al obtener usuario: %v", err)
		}
		if usr.Email != "nicolas@example.com" {
			t.Errorf("Email no coincide. Esperado 'nicolas@example.com', obtenido '%s'", usr.Email)
		}
	})

	// --- C. ACTUALIZAR (Update) ---
	t.Run("UpdateUsuario", func(t *testing.T) {
		err := queries.UpdateUsuario(ctx, UpdateUsuarioParams{
			IDUsuario: createdID,
			Nombre:    "Nicolas Modificado",
			Apellido:  "Krohn",
			Dni:       "40123456",
			Email:     "nicolas.updated@example.com",
			Password:  "newpassword123",
		})
		if err != nil {
			t.Fatalf("Error al actualizar usuario: %v", err)
		}

		// Verificar que el cambio se guardó
		usr, _ := queries.GetUsuario(ctx, createdID)
		if usr.Nombre != "Nicolas Modificado" {
			t.Errorf("Nombre no actualizado. Esperado 'Nicolas Modificado', obtenido '%s'", usr.Nombre)
		}
	})

	// --- D. LISTAR (List) ---
	t.Run("ListUsuarios", func(t *testing.T) {
		usuarios, err := queries.ListUsuarios(ctx)
		if err != nil {
			t.Fatalf("Error al listar usuarios: %v", err)
		}
		if len(usuarios) == 0 {
			t.Error("Se esperaba al menos 1 usuario en la lista")
		}
	})

	// --- E. ELIMINAR (Delete) ---
	t.Run("DeleteUsuario", func(t *testing.T) {
		err := queries.DeleteUsuario(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al eliminar usuario: %v", err)
		}

		// Verificar que el usuario fue borrado correctamente
		_, err = queries.GetUsuario(ctx, createdID)
		if err == nil {
			t.Error("Se esperaba un error al buscar un usuario eliminado")
		}
	})
}