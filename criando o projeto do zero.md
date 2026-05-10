# Viewax

Viewax é um monitor simples de saúde do sistema para terminal, pensado para ser mais amigável que ferramentas como `top`, `htop` ou `nmon`.

Ele mostra CPU, memória, swap, disco, uptime, load average e processos de forma clara, com interpretação simples para usuários menos técnicos.

## Recursos atuais

- Visão simples do sistema
- Visão Pro com métricas detalhadas
- Top processos por CPU
- Top processos por memória
- Saída JSON
- Modo watch com atualização automática
- Interface colorida no terminal
- Compatível com macOS e Linux

## Requisitos

- Go instalado
- macOS ou Linux

Verifique se o Go está instalado:

```bash
go version
```

Caso não esteja instalado:

### macOS
Com Homebrew:

```bash
brew install go
```

### Linux 

Ubuntu/Debian:

```bash
sudo apt update
sudo apt install golang-go
```

Ou instale pela documentação oficial:

[https://go.dev/doc/install](https://go.dev/doc/install)

# Criando o projeto do zero

```bash 
mkdir viewax
cd viewax
go mod init viewax
```

### Crie a estrutura:

```bash 
mkdir -p cmd/viewax
mkdir -p internal/system
mkdir -p internal/ui
```
Arquivos principais:

```bash 
cmd/viewax/main.go
internal/system/snapshot.go
internal/ui/render.go
```

### Dependências

Instale as dependências:

```bash
go get github.com/shirou/gopsutil/v4/cpu
go get github.com/shirou/gopsutil/v4/disk
go get github.com/shirou/gopsutil/v4/host
go get github.com/shirou/gopsutil/v4/load
go get github.com/shirou/gopsutil/v4/mem
go get github.com/shirou/gopsutil/v4/process
go get github.com/charmbracelet/lipgloss
```
Depois rode:

```bash
go mod tidy
```

### Rodando o projeto

Visão simples:

```bash
go run ./cmd/viewax
```

Visão Pro:

```bash
go run ./cmd/viewax --mode pro
```

Processos por CPU:

```bash
go run ./cmd/viewax --mode processes
```
Processos por memória:

```bash
go run ./cmd/viewax --mode memory
```

Saída JSON:

```bash
go run ./cmd/viewax --mode json
```

Ajuda:

```bash
go run ./cmd/viewax --help
```

Rodar apenas uma vez, sem atualização contínua:

```bash
go run ./cmd/viewax --mode pro --watch=false
```

## Compilando

Gerar binário local:

```bash
go build -o viewax ./cmd/viewax
```

Executar:

```bash
./viewax
```

## Compilando para múltiplas plataformas

Crie a pasta de release:

```bash
mkdir -p release
```

macOS Apple Silicon:

```bash
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o release/viewax-macos-arm64 ./cmd/viewax
```

macOS Intel:

```bash
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o release/viewax-macos-amd64 ./cmd/viewax
```

Linux Intel/AMD:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/viewax-linux-amd64 ./cmd/viewax
```

Linux ARM64:

```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o release/viewax-linux-arm64 ./cmd/viewax
```

Verifique os binários:

```bash
ls -lh release
```

### Exemplo de uso após compilar

```bash
./viewax
./viewax --mode pro
./viewax --mode processes
./viewax --mode memory
./viewax --mode json
./viewax --help
```

### Estrutura do projeto

```bash
viewax/
├── cmd/
│   └── viewax/
│       └── main.go
├── internal/
│   ├── system/
│   │   └── snapshot.go
│   └── ui/
│       └── render.go
├── go.mod
├── go.sum
└── README.md
```

### Stack usada

* Go
* gopsutil
* Lip Gloss
