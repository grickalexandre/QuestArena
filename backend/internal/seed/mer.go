package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	merQuizPrefix = "seed-mer-gen-auto-"
	merTitle      = "Quest 5 — MER: generalização e autorrelacionamento"
	merDesc       = "15 questões do material MER: 10 objetivas (1 min) e 5 dissertativas (5 min) sobre especialização, papéis e modelo físico."
)

func merPack() pack {
	return pack{
		idPrefix:  merQuizPrefix,
		title:     merTitle,
		desc:      merDesc,
		questions: merQuestions(),
	}
}

// EnsureMerQuiz cria o quiz de generalização e autorrelacionamento se o professor ainda não o tiver.
func EnsureMerQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return merPack().ensure(ctx, st, teacherID)
}

// IsMerSeedQuiz reports whether a quiz was created by this seed pack.
func IsMerSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, merQuizPrefix)
}

func merQuestions() []draftQuestion {
	return []draftQuestion{
		{
			text: "No MER desta aula, qual pergunta identifica generalização / especialização?",
			options: []string{
				"Isto se liga a outro da mesma entidade?",
				"Isto é um tipo daquilo?",
				"O relacionamento tem atributo quantidade?",
				"A tabela precisa de hífen no nome?",
			},
			correctIndex: 1,
		},
		{
			text: "Aluno e Professor herdam atributos de Pessoa (nome, CPF, e-mail). Isso é:",
			options: []string{
				"Autorrelacionamento 1:N",
				"Especialização: Aluno e Professor são subtipos de Pessoa",
				"Relacionamento ternário",
				"Entidade associativa",
			},
			correctIndex: 1,
		},
		{
			text: "“Um funcionário é chefe de vários funcionários.” O modelo correto é:",
			options: []string{
				"Especialização: Funcionário → Chefe e Subordinado, sem relacionamento",
				"Duas entidades distintas: Chefe e Funcionário, sempre",
				"Autorrelacionamento em Funcionário, papéis chefe e subordinado, em geral 1:N",
				"Generalização de Departamento",
			},
			correctIndex: 2,
		},
		{
			text: "Conta → Corrente e Poupança é:",
			options: []string{
				"Autorrelacionamento, porque uma conta transfere para outra",
				"Especialização: corrente e poupança são contas, cada uma com dados extras",
				"O mesmo que “conta transfere para conta”",
				"Um relacionamento entre Conta e Agência",
			},
			correctIndex: 1,
		},
		{
			text: "Especialização total significa que:",
			options: []string{
				"Toda instância do supertipo pertence a algum subtipo (não existe pessoa “solta”)",
				"Pode existir pessoa que não é aluno nem professor",
				"Só existe um subtipo",
				"O relacionamento é obrigatoriamente N:N",
			},
			correctIndex: 0,
		},
		{
			text: "Aluno cursa Disciplina. Por que isso NÃO é especialização?",
			options: []string{
				"Porque toda especialização é 1:N",
				"Porque são entidades diferentes ligadas por um fato; aluno não é um tipo de disciplina",
				"Porque disciplina não pode ter PK",
				"Porque cursar é sempre autorrelacionamento",
			},
			correctIndex: 1,
		},
		{
			text: "Na tabela aluno (especialização de pessoa), a PK correta no padrão da aula é:",
			options: []string{
				"PessoaId (PK e FK, relacionamento 1:1)",
				"aluno-id",
				"id",
				"AlunoId (PK nova, além da FK PessoaId)",
			},
			correctIndex: 0,
		},
		{
			text: "Na tabela carro, especialização de veiculo, o correto é:",
			options: []string{
				"PK CarroId e FK VeiculoId (duas colunas)",
				"VeiculoId é PK e FK (1:1): a identidade do carro é a do veículo",
				"PK id e FK veiculo",
				"Não precisa de FK; carro e veiculo são a mesma tabela",
			},
			correctIndex: 1,
		},
		{
			text: "Tabela menu, PK MenuId. A FK do menu pai deve se chamar:",
			options: []string{
				"MenuId (o mesmo nome da PK, na mesma tabela)",
				"idPai",
				"MenuIdPai (Id no meio)",
				"MenuPaiId (papel antes do Id; sempre termina em Id)",
			},
			correctIndex: 3,
		},
		{
			text: "Peça composta de Peça, com quantidade, é autorrelacionamento:",
			options: []string{
				"1:1, sem atributo",
				"somente especialização",
				"1:N, sem atributo quantidade",
				"N:N, com tabela associativa e atributo quantidade",
			},
			correctIndex: 3,
		},
		{
			text: "Explique a diferença entre generalização/especialização e autorrelacionamento. Dê um exemplo de cada.",
			expectedAnswer: "Especialização: Aluno é um tipo de Pessoa. Autorrelacionamento: Funcionário supervisiona Funcionário.",
			expectedAnswers: []string{
				"É um tipo de? → especialização. Liga dois da mesma entidade? → autorrelacionamento.",
				"Generalização junta tipos parecidos. Autorrelacionamento é a mesma entidade em papéis diferentes (chefe e subordinado).",
			},
			keyTerms:     []string{"especialização", "pessoa", "autorrelacionamento", "funcionário"},
			threshold:    0.45,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Escreva o modelo físico de Pessoa com especialização em Aluno e Professor (tabelas, PK, FK e a cardinalidade 1:1).",
			expectedAnswer: "pessoa (PessoaId PK). aluno e professor com PessoaId PK+FK para pessoa, relacionamento 1:1. Não criar AlunoId.",
			expectedAnswers: []string{
				"PK PessoaId ---------- FK PessoaId em aluno e professor, 1:1.",
				"Três tabelas. A filha reusa PessoaId como PK e FK.",
			},
			keyTerms:     []string{"pessoa", "PessoaId", "aluno", "professor", "1:1"},
			code: `pessoa     (PessoaId PK, nome, dataNasc, email)
aluno      ( ? )
professor  ( ? )`,
			codeLanguage: "sql",
			threshold:    0.45,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Por que não criar a entidade Chefe? Como fica o físico do funcionário que supervisiona outro funcionário?",
			expectedAnswer: "Chefe também é funcionário. Autorrelacionamento 1:N: PK FuncionarioId e FK FuncionarioChefeId na mesma tabela. NULL no diretor.",
			expectedAnswers: []string{
				"Uma tabela funcionario. FK FuncionarioChefeId aponta para FuncionarioId.",
				"Não especializar em Chefe e Subordinado. Papel antes do Id.",
			},
			keyTerms:     []string{"funcionário", "chefe", "autorrelacionamento", "FuncionarioChefeId", "1:N"},
			threshold:    0.45,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "“Pasta contém pasta” em um drive na nuvem. Isso é especialização ou autorrelacionamento? Qual a cardinalidade e o físico da tabela pasta?",
			expectedAnswer: "Autorrelacionamento 1:N em pasta. PK PastaId, FK PastaPaiId. Raiz com pai NULL. Não é especialização.",
			expectedAnswers: []string{
				"pasta (PastaId PK, PastaPaiId FK). Uma pasta contém várias; a filha tem no máximo um pai.",
				"Mesma entidade em árvore, papéis pai e filha.",
			},
			keyTerms:     []string{"autorrelacionamento", "1:N", "PastaId", "PastaPaiId"},
			threshold:    0.45,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Uma peça é composta de outras peças, com quantidade. Escreva o conceitual (cardinalidade) e o modelo físico (incluindo a tabela associativa e os nomes das FKs).",
			expectedAnswer: "Autorrelacionamento N:N. peca (PecaId). pecacomposicao com PecaId, PecaItemId e quantidade.",
			expectedAnswers: []string{
				"N:N precisa de tabela associativa. Uma FK copia a PK; a outra leva o papel.",
				"Quantidade fica na associativa, não em peca.",
			},
			keyTerms:     []string{"N:N", "associativa", "PecaId", "PecaItemId", "quantidade"},
			threshold:    0.45,
			timeLimitSec: timeLimitFiveMin,
		},
	}
}
