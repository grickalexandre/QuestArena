package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	herancaQuizPrefix = "seed-heranca-conta-"
	herancaTitle      = "Quest 3 — Herança até a oficina 6.1"
	herancaDesc       = "15 certo ou errado: herança, classe abstrata, Laboratório A (Animal), Laboratório B (Conta) e construtor com : base. Sem public/private/protected. Tempo folgado — a nota vale pelo acerto, não pela corrida de XP."
	timeLimitHeranca  = 180
)

func herancaPack() pack {
	return pack{
		idPrefix:  herancaQuizPrefix,
		title:     herancaTitle,
		desc:      herancaDesc,
		questions: herancaQuestions(),
	}
}

// EnsureHerancaQuiz cria o quiz da aula (até a oficina 6.1) se o professor ainda não o tiver.
func EnsureHerancaQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return herancaPack().ensure(ctx, st, teacherID)
}

func vf(text string, certa bool) draftQuestion {
	opts := []string{"Certo", "Errado"}
	idx := 0
	if !certa {
		idx = 1
	}
	// Alterna a ordem dos botões para o colega do lado não copiar “o da esquerda”.
	if len(text)%2 == 0 {
		opts[0], opts[1] = opts[1], opts[0]
		idx = 1 - idx
	}
	return draftQuestion{
		text:         text,
		options:      opts,
		correctIndex: idx,
		timeLimitSec: timeLimitHeranca,
	}
}

func herancaQuestions() []draftQuestion {
	return []draftQuestion{
		vf("Na herança, a classe filha reaproveita o que o pai já tem (atributos e métodos) e só escreve o que é específico.", true),
		vf("Se Poção herda de Item, ainda é preciso copiar de novo Nome, Preco e Vender() dentro da classe Poção.", false),
		vf("A linha class Pocao : Item significa que Poção herda de Item.", true),
		vf("Polimorfismo e herança são a mesma coisa: os dois nomes descrevem exatamente o mesmo conceito.", false),
		vf("Polimorfismo, nesta aula, é a mesma ação (Usar) produzir efeitos diferentes conforme o objeto real: poção cura, espada ataca, armadura defende.", true),
		vf("Em Item item = new Pocao(...); item.Usar(); o C# sempre executa o Usar do Item, nunca o da Poção.", false),
		vf("Se Item for abstract, o comando new Item() cria um item genérico vazio e o programa segue.", false),
		vf("A variável pode ser do tipo Item (o molde), mas o objeto criado com new precisa ser de uma classe concreta, como Poção ou Espada.", true),
		vf("Um método abstract no pai (Usar, EmitirSom, Sacar) obriga cada filha concreta a escrever o próprio com override.", true),
		vf("No Laboratório A, new Animal() é válido, porque Animal é o tipo da lista.", false),
		vf("No foreach de List<Animal>, a.EmitirSom() late ou mia conforme o objeto real, sem precisar converter o tipo na hora de emitir o som.", true),
		vf("No Laboratório B, uma conta corrente com saldo 50 e limite 200 aceita sacar 200; o saldo fica −150.", true),
		vf("No Laboratório B, a poupança herda o limite da conta corrente; por isso, com saldo 50, ela também aceita sacar 200.", false),
		vf("O construtor tem o mesmo nome da classe, não declara tipo de retorno e roda automaticamente quando aparece o new.", true),
		vf("O : base(nome, preco) pode ser escrito como um método dentro das chaves do construtor da filha, na primeira linha: base(nome, preco);", false),
	}
}

// IsHerancaSeedQuiz reports whether a quiz was created by this seed pack.
func IsHerancaSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, herancaQuizPrefix)
}
