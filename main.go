package main

import (
	"net/http" // para manejar respuestas http
	"strconv"  // para convertir strings a num y viceversa

	"github.com/gin-gonic/gin" // uso gin para manejar el servidor y las rutas

	"demob/src/application"     
	"demob/src/domain"         
	"demob/src/infraestructure" 
)

func main() {
	mysql := infraestructure.NewMysql()

	myGin := gin.Default()

	createProduct := application.NewCreateUseCase(mysql)       
	getAll := application.NewViewAllUseCase(mysql)             
	getById := application.NewViewByIdProductUseCase(mysql)     
	updateProduct := application.NewUpdateProductUseCase(mysql) 
	deleteProduct := application.NewDeleteProductUseCase(mysql) 


	myGin.GET("/products", func(c *gin.Context) {
		products, err := getAll.Run() 
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	myGin.GET("/products/:id", func(c *gin.Context) {
		// convierto el parmetro id a un num entero
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id no es correcto"})
			return
		}
		product, err := getById.Run(int32(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	myGin.POST("/products", func(c *gin.Context) {
		var newProduct domain.Product // creo una variable para el nuevo producto

		// intento leer el json del body y lo paso a la variable newProduct
		if err := c.ShouldBindJSON(&newProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input invalido: " + err.Error()})
			return
		}

		// corro el caso de uso para crear el producto
		if err := createProduct.Run(newProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el producto"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "producto creado exitosamente"})
	})

	// esta ruta actualiza un producto por su id
	myGin.PUT("/products/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
			return
		}

		var updatedProduct domain.Product // creo una variable para el producto actualizado
		// intento leer el json del body y lo paso a la variable updatedProduct
		if err := c.ShouldBindJSON(&updatedProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// corro el caso de uso para actualizar el producto
		if err := updateProduct.Run(int32(id), updatedProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el producto"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "producto actualizado exitosamente"})
	})

	myGin.DELETE("/products/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
			return
		}

		if err := deleteProduct.Run(int32(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo borrar el producto"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "producto borrado exitosamente"})
	})

	myGin.Run()
}
