package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time" // NOVO: Pacote necessário para registrar o tempo
)

// (Nenhuma mudança nas constantes e variáveis globais)
const maxUploadSize = 10 * 1024 * 1024
const uploadPath = "./uploads"

var allowedExtensions = map[string]bool{
	"txt":  true,
	"adoc": true,
	"md":   true,
	"yaml": true,
	"yml":  true,
}

// NOVO: Middleware de Logging
// Esta função recebe um http.Handler e retorna um novo http.Handler.
// O novo handler primeiro registra os detalhes da requisição e depois chama o handler original.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Imprime o log de acesso no terminal
		log.Printf(
			"Requisição recebida: Método=%s, URI=%s, Remoto=%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		// Chama o próximo handler na cadeia (neste caso, o uploadFileHandler)
		next.ServeHTTP(w, r)
		
		// Loga o tempo que a requisição levou para ser processada
		log.Printf("Requisição concluída em %s", time.Since(start))
	})
}


// A função uploadFileHandler permanece exatamente a mesma
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. VALIDAÇÃO DO MÉTODO HTTP
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// 2. VALIDAÇÃO DO TAMANHO DO CONTEÚDO
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, fmt.Sprintf("O arquivo excede o limite de tamanho de %d MB", maxUploadSize/(1024*1024)), http.StatusBadRequest)
		return
	}

	// 3. OBTENÇÃO DO ARQUIVO
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo. Certifique-se que o campo se chama 'file'.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. VALIDAÇÃO DA EXTENSÃO DO ARQUIVO
	ext := filepath.Ext(handler.Filename)
	ext = strings.TrimPrefix(ext, ".")

	if !allowedExtensions[ext] {
		allowed := make([]string, 0, len(allowedExtensions))
		for k := range allowedExtensions {
			allowed = append(allowed, k)
		}
		errMsg := fmt.Sprintf("Tipo de arquivo não permitido. Extensões aceitas: %s", strings.Join(allowed, ", "))
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// 5. SALVANDO O ARQUIVO NO SERVIDOR
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		http.Error(w, "Não foi possível criar o diretório no servidor.", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filepath.Join(uploadPath, handler.Filename))
	if err != nil {
		http.Error(w, "Não foi possível salvar o arquivo.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Erro ao copiar o conteúdo do arquivo.", http.StatusInternalServerError)
		return
	}

	// 6. RESPOSTA DE SUCESSO
	fmt.Fprintf(w, "Upload do arquivo '%s' realizado com sucesso!", handler.Filename)
}

func main() {
	// NOVO: Criamos um http.Handler a partir da nossa função original.
	uploadHandler := http.HandlerFunc(uploadFileHandler)

	// NOVO: Em vez de registrar o handler diretamente, o "envolvemos" com o nosso middleware de logging.
	// A função http.Handle é usada aqui pois estamos trabalhando com http.Handler.
	http.Handle("/upload", loggingMiddleware(uploadHandler))

	port := "8080"
	log.Printf("Servidor iniciado em http://localhost:%s", port)
	log.Println("Endpoint de upload disponível em POST http://localhost:8080/upload")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}