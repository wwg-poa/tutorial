package main

import "fmt"
import "os"
import "math/rand"
import "time"

type Entrada int

// Possiveis entradas do usuário
const (
  Nada = iota
  ParaCima
  ParaBaixo
  ParaEsquerda
  ParaDireita
  SairDoJogo // Tecla ESC
)

type PacGo struct {
  posicao    Posicao
  figura     string
}

type Fantasma struct {
	posicao Posicao
	figura  string
}

var pacgo PacGo
var fantasmas []*Fantasma

func criarPacGo(pos Posicao, fig string) {
  pacgo = PacGo{
    posicao: pos,
    figura: fig,
  }
}

func criarFantasma(posicao Posicao, figura string) {
	fantasma := &Fantasma{posicao: posicao, figura: figura}
	fantasmas = append(fantasmas, fantasma)
}

func leEntradaDoUsuario(canal chan<- Entrada) {
	array := make([]byte, 10)

	for {
		lido, _ := os.Stdin.Read(array)

		if lido == 1 && array[0] == 0x1b {
			canal <- SairDoJogo
		} else if lido == 3 {
			if array[0] == 0x1b && array[1] == '[' {
				switch array[2] {
				case 'A': canal <- ParaCima
				case 'B':	canal <- ParaBaixo
				case 'C':	canal <- ParaDireita
				case 'D': canal <- ParaEsquerda
				}
			}
		}
	}
}

type Labirinto struct {
  largura       int
  altura        int
  mapa          []string
}

var labirinto Labirinto

func inicializarLabirinto() {
  labirinto = Labirinto{
    largura: 20,
    altura : 10,
    mapa   : []string{
      "####################",
      "#                 F#",
      "#                  #",
      "#                  #",
      "#                  #",
      "#                  #",
      "#                  #",
      "#                  #",
      "#G                F#",
      "####################",
    },
  }

  // Processa caracteres especiais
  for linha, linhaMapa := range labirinto.mapa {
    for coluna, caractere := range linhaMapa {
      switch( caractere ) {
        case 'G': { criarPacGo(Posicao{linha, coluna}, "G") }
        case 'F': { criarFantasma(Posicao{linha, coluna}, "F") }
      }
    }
  }
}

func desenhaTela() {
  LimpaTela()

  // Imprime mapa
  for _, linha := range labirinto.mapa {
    for _, char := range linha {
      switch char {
        case '#': fmt.Print("#")
        default:  fmt.Print(" ")
      }
    }
    fmt.Println("")
  }

  // Imprime PacGo
  MoveCursor(pacgo.posicao)
  fmt.Printf("%s", pacgo.figura)

  // Imprime fantasmas
	for _, fantasma := range fantasmas {
		MoveCursor(fantasma.posicao)
		fmt.Printf("%s", fantasma.figura)
	}

  // Move cursor para fora do labirinto
  MoveCursor(Posicao{labirinto.altura + 2, 0})
}

func moverPacGo(m Entrada) {
  var novaLinha = pacgo.posicao.linha
  var novaColuna = pacgo.posicao.coluna

  switch m {
    case ParaCima:
      novaLinha--
      if novaLinha < 0 {
        novaLinha = labirinto.altura - 1
      }
    case ParaBaixo:
      novaLinha++
      if novaLinha >= labirinto.altura {
        novaLinha = 0
      }
    case ParaDireita:
      novaColuna++
      if novaColuna >= labirinto.largura {
        novaColuna = 0
      }
    case ParaEsquerda:
      novaColuna--
      if novaColuna < 0 {
        novaColuna = labirinto.largura - 1
      }
  }

  conteudoDoMapa := labirinto.mapa[novaLinha][novaColuna]
  if conteudoDoMapa != '#' {
    pacgo.posicao.linha = novaLinha
    pacgo.posicao.coluna = novaColuna
  }
}

func moverFantasmas() {
  for _, fantasma := range fantasmas {
    // gera um número aleatório entre 0 e 4 (Direita = 3)
    var direcao = rand.Intn(ParaDireita+1)

    var novaPosicao = fantasma.posicao

    // Atualiza posição testando os limites do mapa
    switch direcao {
    case ParaCima:
      novaPosicao.linha--
      if novaPosicao.linha < 0 { novaPosicao.linha = labirinto.altura - 1 }
    case ParaBaixo:
      novaPosicao.linha++
      if novaPosicao.linha > labirinto.altura - 1 { novaPosicao.linha = 0 }
    case ParaEsquerda:
      novaPosicao.coluna--
      if novaPosicao.coluna < 0 { novaPosicao.coluna = labirinto.largura - 1 }
    case ParaDireita:
      novaPosicao.coluna++
      if novaPosicao.coluna > labirinto.largura - 1 { novaPosicao.coluna = 0 }
    }

    // Verifica se a posição nova do mapa é válida
    conteudoMapa := labirinto.mapa[novaPosicao.linha][novaPosicao.coluna]
    if conteudoMapa != '#' { fantasma.posicao = novaPosicao }
  }
}

func dorme(milisegundos time.Duration) {
  time.Sleep(time.Millisecond * milisegundos)
}

func main() {
  // Inicializar terminal
  Inicializa()
  defer Finaliza()

  // Inicializar labirinto
  inicializarLabirinto()

  // Cria rotina para ler entradas
  canal := make(chan Entrada, 10)
	go leEntradaDoUsuario(canal)

  // Loop principal
  for {
    // Desenha tela
    desenhaTela()

    // Processa entrada do jogador
    var tecla Entrada
		select {
		case tecla = <-canal:
		default:
		}

    if tecla == SairDoJogo {
      break
    } else {
      moverPacGo(tecla)
    }

    // Processa movimento dos fantasmas
    moverFantasmas()

    // Processa colisões

    // Dorme
    dorme(100)
  }
}
