package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/models"
	"github.com/questarena/questarena/internal/store"
)

const (
	herancaQuizPrefix = "seed-heranca-conta-"
	herancaTitle      = "Herança C# — até Conta Corrente"
	herancaDesc       = "15 questões de 1 minuto: herança, virtual/override, polimorfismo, classe abstrata, Animal e Conta Corrente."
	timeLimitOneMin   = 60
)

type draftQuestion struct {
	text         string
	options      []string
	correctIndex int
}

func herancaQuizID(teacherID string) string {
	return herancaQuizPrefix + teacherID
}

// EnsureHerancaQuiz cria o quiz da aula (até Conta Corrente) se o professor ainda não o tiver.
func EnsureHerancaQuiz(ctx context.Context, st store.Store, teacherID string) error {
	if teacherID == "" {
		return nil
	}
	id := herancaQuizID(teacherID)
	existing, err := st.GetQuiz(ctx, id)
	if err != nil || existing == nil {
		q := &models.Quiz{
			ID:          id,
			TeacherID:   teacherID,
			Title:       herancaTitle,
			Description: herancaDesc,
		}
		if err := st.CreateQuiz(ctx, q); err != nil {
			return err
		}
	}
	qs, err := st.ListQuestions(ctx, id)
	if err != nil {
		return err
	}
	if len(qs) >= len(herancaQuestions()) {
		return nil
	}
	for i, d := range herancaQuestions() {
		already := false
		for _, q := range qs {
			if q.Order == i {
				already = true
				break
			}
		}
		if already {
			continue
		}
		item := &models.Question{
			QuizID:       id,
			Type:         models.QuestionMultipleChoice,
			Text:         d.text,
			Options:      append([]string{}, d.options...),
			CorrectIndex: d.correctIndex,
			Weight:       1,
			TimeLimitSec: timeLimitOneMin,
			Order:        i,
		}
		if err := st.CreateQuestion(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

func herancaQuestions() []draftQuestion {
	return []draftQuestion{
		{
			text: "O que é herança em POO, nesta aula?",
			options: []string{
				"Copiar o código do pai em outro arquivo, sem relação",
				"A classe filha recebe o que o pai já tem e pode acrescentar ou especializar",
				"Proibir o filho de ter métodos novos",
				"Criar um objeto sem usar classe",
			},
			correctIndex: 1,
		},
		{
			text: "Em C#, como se escreve “Poção herda de Item”?",
			options: []string{
				"class Pocao : Item",
				"class Pocao extends Item",
				"class Pocao implements Item",
				"class Pocao -> Item",
			},
			correctIndex: 0,
		},
		{
			text: "No RPG, o que o pai Item já traz para Poção, Espada e Armadura?",
			options: []string{
				"Só Curativa e Dano",
				"Apenas o método Main",
				"Nome, Preço, Usar() e Vender()",
				"Só o preço da loja",
			},
			correctIndex: 2,
		},
		{
			text: "Por que a herança sozinha não resolve o Usar() do RPG?",
			options: []string{
				"Todas as filhas herdariam a mesma mensagem “você usou…”",
				"Porque C# não permite herança",
				"Porque Item não pode ter Nome",
				"Porque lista não existe em C#",
			},
			correctIndex: 0,
		},
		{
			text: "O que significa virtual no método do pai?",
			options: []string{
				"O método fica privado",
				"O pai autoriza os filhos a mudarem este método",
				"A classe vira abstrata automaticamente",
				"O método não pode ser chamado",
			},
			correctIndex: 1,
		},
		{
			text: "O que significa override na classe filha?",
			options: []string{
				"A filha está reescrevendo o método do pai",
				"A filha apaga a classe pai",
				"Cria um novo objeto",
				"Impede a herança",
			},
			correctIndex: 0,
		},
		{
			text: "Qual par permite reescrever o método do pai corretamente?",
			options: []string{
				"static e final",
				"private e public",
				"virtual e override",
				"new e class",
			},
			correctIndex: 2,
		},
		{
			text: "O que faz base.Usar() dentro do override?",
			options: []string{
				"Executa também a versão do método do pai",
				"Cria um novo Item",
				"Impede o override",
				"Apaga a classe pai",
			},
			correctIndex: 0,
		},
		{
			text: "O que é polimorfismo nesta aula?",
			options: []string{
				"Ter várias classes sem nenhuma relação",
				"Mesma ação (Usar) com efeitos diferentes conforme o objeto",
				"Só usar variáveis public",
				"Proibir herança",
			},
			correctIndex: 1,
		},
		{
			text: "Em Item item = new Pocao(); item.Usar(); qual Usar() roda?",
			options: []string{
				"Sempre o Usar do Item (rótulo da variável)",
				"O Usar da Poção, porque o objeto real é Poção",
				"Dá erro de compilação",
				"Os dois automaticamente, sem usar base",
			},
			correctIndex: 1,
		},
		{
			text: "Se Item for abstract, o que acontece com new Item()?",
			options: []string{
				"Cria um Item genérico normalmente",
				"Erro: classe abstrata não vira objeto sozinha",
				"Cria um Item vazio em silêncio",
				"Só funciona no Programiz",
			},
			correctIndex: 1,
		},
		{
			text: "Por que Animal é abstract no laboratório?",
			options: []string{
				"Não existe “animal genérico”: só cachorro, gato, papagaio…",
				"C# exige Main em toda classe",
				"EmitirSom não pode ter override",
				"List não funciona com classes concretas",
			},
			correctIndex: 0,
		},
		{
			text: "No foreach de List<Animal>, a.EmitirSom() late ou mia?",
			options: []string{
				"Sempre late, porque a lista é de Animal",
				"Depende do objeto real: Cachorro late, Gato mia…",
				"Dá erro porque Animal é abstract",
				"Só funciona se fizer cast em todo mundo",
			},
			correctIndex: 1,
		},
		{
			text: "Ana tem conta corrente (saldo 50, limite 200). Sacar 200. O que acontece?",
			options: []string{
				"Negado, porque 200 > 50",
				"Aceito: saldo fica −150 (usou o cheque especial)",
				"Aceito, mas o limite some para sempre",
				"A conta vira poupança automaticamente",
			},
			correctIndex: 1,
		},
		{
			text: "Por que ((ContaCorrente)contaAna).Limite precisa de cast?",
			options: []string{
				"Porque Limite é abstract",
				"Porque a variável é Conta e Limite só existe na corrente",
				"Porque poupança também tem Limite",
				"Porque decimal exige cast sempre",
			},
			correctIndex: 1,
		},
	}
}

// IsHerancaSeedQuiz reports whether a quiz was created by this seed pack.
func IsHerancaSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, herancaQuizPrefix)
}
