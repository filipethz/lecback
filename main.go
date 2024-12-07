package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Definir a estrutura para o banco de dados
type Purchase struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Data        string  `json:"data"` // Mudar para string temporariamente
	Cliente     string  `json:"cliente"`
	Dados       string  `json:"dados"`
	NomeProduto string  `json:"nome_produto"`
	Valor       float64 `json:"valor"`
	Pagamento   string  `json:"pagamento"`
	Obs         string  `json:"obs"`
}

// Função para conectar ao banco de dados
func connectDB() {
	// Ajuste a string de conexão conforme necessário
	dsn := "filipethomaz10:Gremio1903RocketLeague@tcp(mysql16-farm10.kinghost.net:3306)/filipethomaz10"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
}

// Função para migrar o banco de dados
func migrate() {
	// Realiza a migração para criar a tabela de compras se ela não existir
	db.AutoMigrate(&Purchase{})
}

// Função para converter string para time.Time
func parseDate(dateStr string) (time.Time, error) {
	// Formato esperado para a data (YYYY-MM-DD)
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

// Função para cadastrar uma nova compra
func createPurchase(c echo.Context) error {
	var purchase Purchase
	if err := c.Bind(&purchase); err != nil {
		log.Printf("Erro ao fazer bind dos dados: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
	}

	// Converte a data recebida de string para time.Time
	if purchase.Data != "" {
		parsedDate, err := parseDate(purchase.Data)
		if err != nil {
			log.Printf("Erro ao converter a data: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid date format"})
		}
		// Atribuir a data convertida ao campo Data
		purchase.Data = parsedDate.Format("2006-01-02")
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid or missing date"})
	}

	// Salvar a compra no banco de dados
	if err := db.Create(&purchase).Error; err != nil {
		log.Printf("Erro ao salvar compra no banco: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving purchase"})
	}

	return c.JSON(http.StatusOK, purchase)
}

// Função para listar as compras
func getPurchases(c echo.Context) error {
	var purchases []Purchase
	if err := db.Find(&purchases).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching purchases"})
	}
	fmt.Println(purchases) // Logar as compras no backend
	return c.JSON(http.StatusOK, purchases)
}

func main() {
	// Conectar ao banco de dados
	connectDB()

	// Realizar a migração para criar a tabela de compras
	migrate()

	// Criar uma nova instância do Echo
	e := echo.New()

	// Habilitar o middleware CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                              // Permitir todas as origens (modificar conforme necessário)
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},       // Métodos permitidos
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization}, // Cabeçalhos permitidos
	}))

	// Ativar o middleware de logger e recovery
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Definir as rotas
	e.POST("/purchase", createPurchase) // Rota para cadastrar compras
	e.GET("/purchases", getPurchases)   // Rota para listar compras

	// Iniciar o servidor na porta 8080
	e.Logger.Fatal(e.Start(":8080"))
}
