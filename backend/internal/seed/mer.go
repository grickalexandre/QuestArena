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
			expectedAnswer: "Generalização ou especialização pergunta se isto é um tipo daquilo: Aluno é uma Pessoa. Autorrelacionamento é a mesma entidade em papéis diferentes: Funcionário supervisiona Funcionário (chefe e subordinado).",
			expectedAnswers: []string{
				"Especialização: subtipo de um supertipo (Pessoa → Aluno e Professor). Autorrelacionamento: a entidade se relaciona com ela mesma, com papéis nomeados.",
				"É um tipo de? → especialização. Liga dois da mesma entidade? → autorrelacionamento.",
			},
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Escreva o modelo físico de Pessoa com especialização em Aluno e Professor (tabelas, PK, FK e a cardinalidade 1:1).",
			expectedAnswer: "Tabela pessoa com PK PessoaId e dados comuns. Tabelas aluno e professor usam PessoaId como PK e FK para pessoa, relacionamento 1:1. Não criar AlunoId. Aluno tem ra, curso; professor tem formacao.",
			expectedAnswers: []string{
				"pessoa (PessoaId PK). aluno (PessoaId PK+FK → pessoa, 1:1). professor (PessoaId PK+FK → pessoa, 1:1).",
				"PK PessoaId ---------- FK PessoaId em aluno e professor, relacionamento 1:1. A identidade do aluno é a da pessoa.",
			},
			code: `pessoa     (PessoaId PK, nome, dataNasc, email)
aluno      ( ? )
professor  ( ? )`,
			codeLanguage: "sql",
			threshold:    0.48,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Por que não criar a entidade Chefe? Como fica o físico do funcionário que supervisiona outro funcionário?",
			expectedAnswer: "O chefe também é funcionário, então não se cria entidade Chefe. É autorrelacionamento 1:N na tabela funcionario: PK FuncionarioId e FK FuncionarioChefeId apontando para a mesma tabela. NULL no diretor, que não tem chefe.",
			expectedAnswers: []string{
				"Uma tabela funcionario. PK FuncionarioId. FK FuncionarioChefeId → funcionario(FuncionarioId). Papel antes do Id porque FuncionarioId já é a PK.",
				"Autorrelacionamento 1:N: um chefe tem vários subordinados. Não especializar em Chefe e Subordinado.",
			},
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "“Pasta contém pasta” em um drive na nuvem. Isso é especialização ou autorrelacionamento? Qual a cardinalidade e o físico da tabela pasta?",
			expectedAnswer: "Autorrelacionamento em pasta, papéis pasta pai e pasta filha, cardinalidade 1:N. Não é especialização: a filha não é outro tipo de pasta. Tabela pasta com PK PastaId e FK PastaPaiId. A raiz tem PastaPaiId NULL.",
			expectedAnswers: []string{
				"É a mesma entidade em níveis da árvore. pasta (PastaId PK, PastaPaiId FK → pasta). Uma pasta contém várias; cada filha tem no máximo um pai.",
				"Autorrelacionamento 1:N, não especialização. Papel antes do Id: PastaPaiId.",
			},
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Uma peça é composta de outras peças, com quantidade. Escreva o conceitual (cardinalidade) e o modelo físico (incluindo a tabela associativa e os nomes das FKs).",
			expectedAnswer: "Autorrelacionamento N:N: uma peça entra em várias montagens e uma montagem tem várias peças. Tabela peca (PecaId PK). Tabela associativa pecacomposicao com PecaId (conjunto), PecaItemId (componente) e quantidade. As duas FKs apontam para peca.",
			expectedAnswers: []string{
				"N:N precisa de tabela associativa. pecacomposicao (PecaComposicaoId, PecaId, PecaItemId, quantidade). Uma FK copia a PK; a outra leva o papel.",
				"Quantidade é atributo do relacionamento e fica na associativa, não em peca.",
			},
			threshold:    0.48,
			timeLimitSec: timeLimitFiveMin,
		},
	}
}
