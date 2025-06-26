package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time" // Adicionado para incluir um timestamp no status

	"github.com/gin-gonic/gin"
)

// Constantes da aplicação
const (
	maxUploadSize = 5 * 1024 * 1024 // 5 MB
	uploadPath    = "./uploads"
	appVersion    = "1.0.2" // NOVO: Versão da aplicação
)

// Extensões permitidas
var allowedExtensions = map[string]bool{
	"txt":  true,
	"adoc": true,
	"md":   true,
	"yaml": true,
	"yml":  true,
}

// NOVO: Handler para a rota raiz que exibe o status da API
func statusHandler(c *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		// Se houver um erro, logamos no servidor para depuração
		// e definimos um valor padrão para a resposta.
		log.Printf("Erro ao obter o hostname: %v", err)
		hostname = "indisponível"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"version":     appVersion,
		"server_time": time.Now().Format(time.RFC3339), // Adiciona a hora atual do servidor
		"hostname":    hostname, // NOVO: Campo adicionado à resposta JSON
	})
}

// Handler de upload (sem alterações)
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo enviado. O campo deve se chamar 'file'."})
		return
	}

	ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
	if !allowedExtensions[ext] {
		allowed := make([]string, 0, len(allowedExtensions))
		for k := range allowedExtensions {
			allowed = append(allowed, k)
		}
		errMsg := fmt.Sprintf("Tipo de arquivo não permitido. Extensões aceitas: %s", strings.Join(allowed, ", "))
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível criar o diretório no servidor."})
		return
	}

	destinationPath := filepath.Join(uploadPath, file.Filename)
	if err := c.SaveUploadedFile(file, destinationPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível salvar o arquivo."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Upload do arquivo realizado com sucesso!",
		"filename":   file.Filename,
		"size_bytes": file.Size,
	})
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = maxUploadSize
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// --- REGISTRO DAS ROTAS ---
	router.GET("/", statusHandler)
	router.POST("/upload", uploadHandler)

	router.GET("/favicon.ico", func(c *gin.Context) {
        c.Status(http.StatusNoContent)
    })

	port := "8080"
	log.Printf("Servidor Gin v%s iniciado em http://localhost:%s", appVersion, port)
	log.Println("Rota de status disponível em: GET /")
	log.Println("Rota de upload disponível em: POST /upload")

	router.Run(":" + port)
}