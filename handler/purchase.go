package handler

import (
	"encoding/csv"
	"lechrome/model"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// Custom bind function to handle date parsing
func bindPurchaseData(c echo.Context, purchase *model.PurchaseData) error {
	type tempPurchaseData struct {
		User          string `json:"cliente"`
		UserData      string `json:"dados"`
		ProductName   string `json:"nome_produto"`
		Price         string `json:"valor"`
		PaymentMethod string `json:"pagamento"`
		UserInfo      string `json:"obs"`
		Data          string `json:"data"`
	}

	var tempData tempPurchaseData
	if err := c.Bind(&tempData); err != nil {
		return err
	}

	parsedDate, err := time.Parse("2006-01-02", tempData.Data)
	if err != nil {
		return err
	}

	purchase.Data = parsedDate
	purchase.User = tempData.User
	purchase.UserData = tempData.UserData
	purchase.ProductName = tempData.ProductName
	purchase.Price = tempData.Price
	purchase.PaymentMethod = tempData.PaymentMethod
	purchase.UserInfo = tempData.UserInfo

	return nil
}

// CreatePurchase handler to save purchase data
func CreatePurchase(c echo.Context) error {
	var purchase model.PurchaseData
	if err := bindPurchaseData(c, &purchase); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Formatar a data para salvar apenas dia, mês e ano
	formattedDate := purchase.Data.Format("2006-01-02")
	purchase.Data, _ = time.Parse("2006-01-02", formattedDate)

	file, err := os.OpenFile("purchases.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao abrir arquivo: " + err.Error()})
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	record := []string{
		formattedDate,
		purchase.User,
		purchase.UserData,
		purchase.ProductName,
		purchase.Price,
		purchase.PaymentMethod,
		purchase.UserInfo,
	}

	if err := writer.Write(record); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao escrever no arquivo: " + err.Error()})
	}

	writer.Flush()
	return c.JSON(http.StatusCreated, purchase)
}

// GetPurchases handler to retrieve all purchase data
func GetPurchases(c echo.Context) error {
	file, err := os.Open("purchases.csv")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao abrir arquivo: " + err.Error()})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao ler arquivo: " + err.Error()})
	}

	var purchases []model.PurchaseData
	for _, record := range records {
		data, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao parsear data: " + err.Error()})
		}

		purchase := model.PurchaseData{
			Data:          data,
			User:          record[1],
			UserData:      record[2],
			ProductName:   record[3],
			Price:         record[4],
			PaymentMethod: record[5],
			UserInfo:      record[6],
		}
		purchases = append(purchases, purchase)
	}

	return c.JSON(http.StatusOK, purchases)
}
