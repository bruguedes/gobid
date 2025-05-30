package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	// Carregar o arquivo .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
		return
	}

	// Verificar se o comando 'tern' está disponível
	_, err := exec.LookPath("tern")
	if err != nil {
		fmt.Println("O comando 'tern' não foi encontrado no PATH:", err)
		return
	}

	// Configurar e executar o comando
	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("E: %v\nSaída: %s\n", err, string(output))
		return
	}

	// Exibir a saída do comando
	fmt.Println("Saída do comando:", string(output))

	fmt.Println("Tern migrations executed successfully.")
}
