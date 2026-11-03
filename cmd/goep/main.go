package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/internal/cli"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "goep:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return interactive()
	}

	switch args[0] {
	case "create-module", "module":
		return cli.CreateModule(args[1:])
	case "schema", "assemble-schema":
		if len(args) != 1 {
			return errors.New("schema does not accept arguments")
		}
		return cli.AssembleSchema()
	case "generate-module-registry":
		return cli.GenerateModuleRegistry(args[1:])
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		usage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func interactive() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("goEP Core")
	fmt.Println("¿Qué quieres hacer?")
	fmt.Println("  1) Crear un módulo")
	fmt.Println("  2) Ensamblar el esquema PostgreSQL")
	fmt.Println("  3) Regenerar los registros de módulos")
	fmt.Println("  q) Salir")
	fmt.Print("Selecciona una opción: ")

	choice, err := readLine(reader)
	if err != nil {
		return err
	}
	switch choice {
	case "1":
		fmt.Print("Nombre del módulo (ej. sales-reports): ")
		name, err := readLine(reader)
		if err != nil {
			return err
		}
		if name == "" {
			return errors.New("el nombre del módulo no puede estar vacío")
		}
		fmt.Print("Descripción (opcional): ")
		summary, err := readLine(reader)
		if err != nil {
			return err
		}
		args := []string{}
		if summary != "" {
			args = append(args, "--summary", summary)
		}
		return cli.CreateModule(append(args, name))
	case "2":
		return cli.AssembleSchema()
	case "3":
		return cli.GenerateModuleRegistry(nil)
	case "q", "Q":
		return nil
	default:
		return fmt.Errorf("opción no válida: %q", choice)
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func usage() {
	fmt.Println(`Uso:
  go run ./cmd/goep                         Abrir el menú interactivo
  go run ./cmd/goep create-module <nombre>  Crear un módulo
  go run ./cmd/goep schema                  Ensamblar schema.sql
  go run ./cmd/goep generate-module-registry [opciones]

Sin argumentos, el CLI pregunta qué acción quieres realizar.`)
}
