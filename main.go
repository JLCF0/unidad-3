package main

import (
	"database/sql" // Interactuar con bases de datos
	"fmt"          // Imprimir mensajes
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Librería que nos permite conectar a MySQL
)

func InsertarLibro(c *gin.Context) {
	nombre := c.Query("nombre")
	descripcion := c.Query("descripcion")
	autor := c.Query("autor")
	editorial := c.Query("editorial")
	fecha := c.Query("fecha_publicacion")

	var libro Libro

	libro.Nombre = nombre
	libro.Descripcion = descripcion
	libro.Autor = autor
	libro.Editorial = editorial
	libro.Fecha_publicacion = fecha

	var resultado string
	err := insertar(libro)
	if err != nil {
		resultado = "Ocurrio un error al insertar"
	} else {
		resultado = "Insertado correctamente"
	}

	c.JSON(200, gin.H{
		"Resultado": resultado,
	})
}

func MostrarLibro(c *gin.Context) {
	idstr := c.Query("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		c.JSON(200, gin.H{
			"Resultado": "Ingrese un valor numerico para el ID",
		})
	} else {
		libro, err := obtenerLibro(id)
		if err != nil {
			c.JSON(200, gin.H{
				"Resultado": "Ocurrio un error al mostrar el libro",
			})
		} else {
			c.JSON(200, gin.H{
				"Libro": "ID: " + strconv.Itoa(libro.Id) + " Nombre: " + libro.Nombre +
					" Descripcion: " + libro.Descripcion + " Autor: " + libro.Autor + " Editoral: " + libro.Editorial + " Fecha editorial: " + libro.Fecha_publicacion,
			})
		}

	}
}

func MostrarTodosLibros(c *gin.Context) {

	libros, err := obtenerLibros()
	if err != nil {
		c.JSON(200, gin.H{
			"Resultado": "Ocurrio un error al mostrar todos",
		})
	} else {
		for _, libro := range libros {
			c.JSON(200, gin.H{
				"Libro": "ID: " + strconv.Itoa(libro.Id) + " Nombre: " + libro.Nombre +
					" Descripcion: " + libro.Descripcion + " Autor: " + libro.Autor + " Editoral: " + libro.Editorial + " Fecha editorial: " + libro.Fecha_publicacion,
			})
		}
	}

}

func ActualizarLibro(c *gin.Context) {
	idstr := c.Query("id")
	nombre := c.Query("nombre")
	descripcion := c.Query("descripcion")
	autor := c.Query("autor")
	editorial := c.Query("editorial")
	fecha := c.Query("fecha_publicacion")

	id, err := strconv.Atoi(idstr)

	var resultado string
	if err != nil {
		resultado = "Ingrese un valor numerico para el ID"
	} else {
		var libro Libro
		libro.Id = id
		libro.Nombre = nombre
		libro.Descripcion = descripcion
		libro.Autor = autor
		libro.Editorial = editorial
		libro.Fecha_publicacion = fecha

		err := actualizar(libro)
		if err != nil {
			resultado = "Ocurrio un error al actualizar"
		} else {
			resultado = "Actualizado correctamente"
		}

	}

	c.JSON(200, gin.H{
		"Resultado": resultado,
	})
}

func EliminarLibro(c *gin.Context) {
	idstr := c.Query("id")

	id, err := strconv.Atoi(idstr)

	var resultado string
	if err != nil {
		resultado = "Ingrese un valor numerico para el ID"
	} else {

		var libro Libro

		libro.Id = id

		err := eliminar(libro)
		if err != nil {
			resultado = "Ocurrio un error al eliminar"
		} else {
			resultado = "Eliminado correctamente"
		}
	}

	c.JSON(200, gin.H{
		"Resultado": resultado,
	})
}

type Libro struct {
	Nombre, Descripcion, Autor, Editorial, Fecha_publicacion string
	Id                                                       int
}

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "root"
	pass := ""
	host := "tcp(127.0.0.1:3306)"
	nombreBaseDeDatos := "biblioteca"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	r := gin.Default()
	r.GET("/Insertar", InsertarLibro)
	r.GET("/Mostrar", MostrarLibro)
	r.GET("/MostrarTodos", MostrarTodosLibros)
	r.GET("/Actualizar", ActualizarLibro)
	r.GET("/Eliminar", EliminarLibro)

	r.Run(":8080")
}

func eliminar(libro Libro) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("DELETE FROM libros WHERE id = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()

	_, err = sentenciaPreparada.Exec(libro.Id)
	if err != nil {
		return err
	}
	return nil
}

func insertar(libro Libro) (e error) {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	// Preparamos para prevenir inyecciones SQL
	sentenciaPreparada, err := db.Prepare("INSERT INTO libros (nombre, descripcion, autor, editorial, fecha_publicacion) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(libro.Nombre, libro.Descripcion, libro.Autor, libro.Editorial, libro.Fecha_publicacion)
	if err != nil {
		return err
	}
	return nil
}

func obtenerLibro(id int) (Libro, error) {
	libro := Libro{}
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return libro, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT id, nombre, descripcion, autor, editorial, fecha_publicacion FROM libros where id = " + strconv.Itoa(id))

	if err != nil {
		return libro, err
	}
	// Si llegamos aquí, significa que no ocurrió ningún error
	defer filas.Close()

	// Recorrer todas las filas
	for filas.Next() {
		err = filas.Scan(&libro.Id, &libro.Nombre, &libro.Descripcion, &libro.Autor, &libro.Editorial, &libro.Fecha_publicacion)

	}
	return libro, nil
}

func obtenerLibros() ([]Libro, error) {
	libros := []Libro{}
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT id, nombre, descripcion, autor, editorial, fecha_publicacion FROM libros")

	if err != nil {
		return nil, err
	}
	// Si llegamos aquí, significa que no ocurrió ningún error
	defer filas.Close()

	// Aquí vamos a "mapear" lo que traiga la consulta en el while de más abajo
	var libro Libro

	// Recorrer todas las filas, en un "while"
	for filas.Next() {
		err = filas.Scan(&libro.Id, &libro.Nombre, &libro.Descripcion, &libro.Autor, &libro.Editorial, &libro.Fecha_publicacion)
		// Al escanear puede haber un error
		if err != nil {
			return nil, err
		}
		// Y si no, entonces agregamos lo leído al arreglo
		libros = append(libros, libro)
	}
	// Vacío o no, regresamos el arreglo de contactos
	return libros, nil
}

func actualizar(libro Libro) error {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE libros SET nombre = ?, descripcion = ?, autor = ?, editorial = ?, fecha_publicacion = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	// Pasar argumentos en el mismo orden que la consulta
	_, err = sentenciaPreparada.Exec(libro.Nombre, libro.Descripcion, libro.Autor, libro.Editorial, libro.Fecha_publicacion, libro.Id)
	return err // Ya sea nil o sea un error, lo manejaremos desde donde hacemos la llamada
}
