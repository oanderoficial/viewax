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

## Como usar 

```bash
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
```
