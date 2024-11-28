package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

// Verifica a conexão VPN
func checkVPNStatus() {
	cmd := exec.Command("nordvpn", "status")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Erro ao verificar status da VPN: %v", err)
	}
	fmt.Println("Status da VPN:")
	fmt.Println(string(output))
}

// Configura o proxy
func proxyHandler(w http.ResponseWriter, r *http.Request) {

	r.Header.Del("User-Agent")
	r.Header.Set("User-Agent", "AnonymousProxy/1.0")
	r.Header.Del("Referer")
	r.Header.Del("Cookie")

	// Encaminha a requisição para o destino original
	proxyURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
		return
	}

	// Copia os cabeçalhos originais
	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao encaminhar requisição", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copia a resposta para o cliente
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	// Verifica o status da VPN antes de iniciar o proxy
	checkVPNStatus()

	// Configura o servidor proxy
	http.HandleFunc("/", proxyHandler)
	fmt.Println("Proxy anônimo iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
