package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	requisitosQuizPrefix = "seed-eng-requisitos-"
	requisitosTitle      = "Quest 8 — Engenharia de Requisitos (ENADE)"
	requisitosDesc       = "20 questões de 1 minuto: o material da lanchonete do Zé e os itens comentados do Enade (2014, 2017 e 2021) sobre requisito, elicitação, especificação e métodos ágeis."
)

func requisitosPack() pack {
	return pack{
		idPrefix:  requisitosQuizPrefix,
		title:     requisitosTitle,
		desc:      requisitosDesc,
		questions: requisitosQuestions(),
	}
}

// EnsureRequisitosQuiz cria o quiz de engenharia de requisitos se o professor ainda não o tiver.
func EnsureRequisitosQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return requisitosPack().ensure(ctx, st, teacherID)
}

// IsRequisitosSeedQuiz reports whether a quiz was created by this seed pack.
func IsRequisitosSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, requisitosQuizPrefix)
}

func reqMC(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: timeLimitOneMin,
	}
}

func requisitosQuestions() []draftQuestion {
	return []draftQuestion{
		reqMC(
			"No material, o que é um requisito?",
			[]string{
				"Só o que o programador decide implementar depois que o código começa",
				"O que o sistema precisa fazer ou respeitar para o cliente ficar satisfeito — uma condição ou capacidade necessária",
				"O teste que o cliente faz no fim do projeto",
				"A linha de base congelada, que não pode mais mudar",
			},
			1,
		),
		reqMC(
			"Sobre o processo de Engenharia de Requisitos, o material afirma que:",
			[]string{
				"Elicitação, especificação e validação se repetem; os requisitos não são definidos uma única vez e congelados",
				"A especificação só acontece depois que o sistema está pronto",
				"A validação substitui a elicitação, então não é preciso conversar com o cliente",
				"O gerenciamento de mudanças proíbe qualquer alteração depois da primeira entrevista",
			},
			0,
		),
		reqMC(
			"No app da Lanchonete do Zé, “quero aumentar as vendas em 20% com pedidos pelo celular” é um requisito de:",
			[]string{
				"Sistema, porque já descreve a API de pagamento",
				"Usuário, porque está na linguagem do cliente que pede o lanche",
				"Negócio: o objetivo da organização, o porquê de construir o sistema",
				"Domínio técnico da cozinha, independente da empresa",
			},
			2,
		),
		reqMC(
			"[ENADE 2021 · Q09] Qual item é um requisito funcional?",
			[]string{
				"O software deve ser operacionalizado no Linux",
				"O tempo de desenvolvimento não deve ultrapassar seis meses",
				"O software deve emitir relatórios de compras a cada quinze dias",
				"O tempo de resposta não deve ultrapassar 30 segundos",
			},
			2,
		),
		reqMC(
			"[ENADE 2017 · Q12] Sistema acadêmico.\n\nR1: o professor lança notas.\nR2: o sistema deve poder ir para outro sistema operacional em no máximo 60 dias.\nR3: o estudante realiza a matrícula.\nR4: a nota aparece em até 2 segundos.\nR5: o auxiliar cadastra o estudante com no máximo 10 minutos de orientação.\n\nSão não funcionais apenas:",
			[]string{
				"R1, R2 e R3",
				"R1, R2 e R5",
				"R1, R3 e R4",
				"R2, R4 e R5",
			},
			3,
		),
		reqMC(
			"Por que “o sistema deve ser rápido” é um requisito ruim, no material?",
			[]string{
				"Porque requisito não funcional precisa ser mensurável: tempo, porcentagem ou quantidade, para poder ser testado",
				"Porque desempenho é sempre requisito funcional",
				"Porque rapidez só pode ser escrita como regra de negócio",
				"Porque a ISO 25010 não trata de desempenho",
			},
			0,
		),
		reqMC(
			"“Pedidos acima de R$ 80 têm entrega grátis.” No material, isso é:",
			[]string{
				"Requisito funcional: a função de calcular a taxa",
				"Requisito não funcional de desempenho",
				"Regra de negócio: política da lanchonete, que existiria mesmo sem o software",
				"História de usuário no formato Como / quero / para",
			},
			2,
		),
		reqMC(
			"[ENADE 2021 · Q14] A integração ERP–CRM atrasou e estourou o orçamento. A equipe se desentende, a qualidade não sabe o que avaliar, surgem solicitantes novos e o patrocinador demora a responder. O gerente deve priorizar o gerenciamento de:",
			[]string{
				"Custo, porque o orçamento estourou",
				"Tempo, porque o projeto atrasou",
				"Escopo, porque apareceram demandas novas",
				"Partes interessadas: os sintomas vêm de pessoas que não foram alinhadas",
			},
			3,
		),
		reqMC(
			"A cozinha da lanchonete funciona de um jeito que ninguém consegue explicar, “porque sempre foi assim”. A técnica de elicitação mais adequada é:",
			[]string{
				"Observação (etnografia): ver o trabalho real no ambiente do usuário",
				"Questionário, para obter números de muitos clientes",
				"Análise só de manuais antigos",
				"Brainstorming, para inventar promoções",
			},
			0,
		),
		reqMC(
			"O analista quer a opinião de 500 clientes espalhados sobre o cardápio. A técnica indicada no material é:",
			[]string{
				"Workshop JAD, reunindo os 500 no mesmo dia",
				"Prototipação de uma tela no Figma",
				"Entrevista aberta com cada um",
				"Questionário: muitas pessoas, dados quantitativos",
			},
			3,
		),
		reqMC(
			"O cliente “só sabe o que quer quando vê a tela”. A técnica de elicitação adequada é:",
			[]string{
				"Questionário fechado, sem mostrar interface",
				"Prototipação: uma versão preliminar para ele experimentar",
				"Linha de base, para congelar o que ainda não foi visto",
				"Teste de estresse da API de pagamento",
			},
			1,
		),
		reqMC(
			"No MoSCoW do app do Zé, “pedido por comando de voz” ficou de fora desta entrega. Essa classificação é:",
			[]string{
				"Must have: obrigatório para o MVP",
				"Should have: importante, mas não vital",
				"Won't have this time: fora do escopo desta entrega",
				"Could have: entra se sobrar tempo",
			},
			2,
		),
		reqMC(
			"Em “Fazer pedido”, calcular a taxa de entrega sempre ocorre. Aplicar cupom só ocorre se o cliente tiver cupom. Na UML isso é:",
			[]string{
				"«include» para calcular a taxa; «extend» para o cupom",
				"«extend» para os dois, porque os dois são funções",
				"«include» para os dois, inclusive o cupom",
				"Generalização entre Cliente e Entregador",
			},
			0,
		),
		reqMC(
			"Os três Cs de uma história de usuário (Ron Jeffries), no material, são:",
			[]string{
				"Código, compilação e configuração",
				"Cartão, conversa e confirmação (critérios de aceitação)",
				"Custo, cronograma e contrato",
				"Cliente, cozinheiro e contador",
			},
			1,
		),
		reqMC(
			"Boehm, como o material usa: “estamos construindo o produto certo?” corresponde a:",
			[]string{
				"Verificação: o teste confere se a taxa foi calculada como a especificação manda",
				"Rastreabilidade da matriz até o caso de teste",
				"Elicitação por questionário",
				"Validação: o sistema (ou o protótipo) é o que o cliente realmente quer",
			},
			3,
		),
		reqMC(
			"A matriz que liga RF02 à entrevista com o Zé, ao caso de uso e aos testes CT05 e CT06 ilustra:",
			[]string{
				"Rastreabilidade: da origem do requisito até a implementação e os testes",
				"MoSCoW, porque ordena Must e Should",
				"Definition of Done da sprint",
				"Teste de instalação no ambiente do cliente",
			},
			0,
		),
		reqMC(
			"[ENADE 2021 · Q12] Quem gerencia o product backlog no Scrum?",
			[]string{
				"O Scrum Master, que também escolhe as funcionalidades",
				"O quadro Kanban, automaticamente",
				"O Product Owner: a lista ordenada de funcionalidades desejadas pelo cliente",
				"O cliente só na retrospectiva, nunca durante a sprint",
			},
			2,
		),
		reqMC(
			"[ENADE 2021 · Q18] Viabilidade já aprovada. Para elicitação, especificação e validação, os artefatos adequados são:",
			[]string{
				"Entrevista com usuários; caso de uso dos requisitos funcionais; protótipo de telas",
				"Estudo de viabilidade; caso de uso dos requisitos funcionais; protótipo de telas",
				"Matriz de rastreabilidade; caso de uso dos requisitos não funcionais; protótipo",
				"Entrevista; caso de uso dos requisitos não funcionais; matriz de rastreabilidade",
			},
			0,
		),
		reqMC(
			"[ENADE 2017 · Q31] Qual associação de nomes está correta?",
			[]string{
				"Segurança = facilidade de entender e operar o sistema",
				"Tempo de resposta = desempenho; acessibilidade para pessoas com deficiência = usabilidade",
				"Confiabilidade = tempo de resposta do processamento",
				"Usabilidade = controle de quem pode acessar cada dado",
			},
			1,
		),
		reqMC(
			"[ENADE 2014 · Q25] São exemplos de requisitos não funcionais (qualidades), e não de tipos de teste:",
			[]string{
				"Segurança, desempenho, estresse e sistema",
				"Usabilidade, segurança, aceitação e confiabilidade",
				"Usabilidade, segurança, desempenho e confiabilidade",
				"Segurança, aceitação, testabilidade e confidencialidade",
			},
			2,
		),
	}
}
