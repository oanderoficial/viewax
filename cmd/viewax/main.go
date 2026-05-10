package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"viewax/internal/system"
	"viewax/internal/ui"
)

func main() {
	mode := flag.String("mode", "simple", "simple | pro | processes | json")
	watch := flag.Bool("watch", true, "atualiza em tempo real")
	help := flag.Bool("help", false, "mostra ajuda")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *mode == "json" {
		s, err := system.CollectSnapshot()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(s)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		s, err := system.CollectSnapshot()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print("\033[H\033[2J")

		switch *mode {
		case "pro":
			fmt.Println(ui.RenderPro(s))
		case "processes":
			fmt.Println(ui.RenderProcesses(s))
		case "memory":
			fmt.Println(ui.RenderMemory(s))
		default:
			fmt.Println(ui.RenderSimple(s))
		}

		if !*watch {
			return
		}

		select {
		case <-time.After(2 * time.Second):
		case <-stop:
			fmt.Println("\nSaindo do Viewax...")
			return
		}
	}
}

func printHelp() {
	fmt.Println(`
Viewax - monitor simples de saúde do sistema

Uso:
  viewax [opções]

Opções:
  --mode simple      Visão simples do sistema
  --mode pro         Visão completa com swap, load, uptime e processos
  --mode processes   Lista processos que mais consomem CPU
  --mode memory      Lista processos que mais consomem memória
  --mode json        Saída em JSON para automações
  --watch=true       Atualiza em tempo real
  --watch=false      Executa uma vez e encerra
  --help             Mostra esta ajuda

Exemplos:
  viewax
  viewax --mode pro
  viewax --mode processes
  viewax --mode memory
  viewax --mode json
  viewax --mode pro --watch=false

Atalhos:
  Ctrl+C             Sair do Viewax
`)
}