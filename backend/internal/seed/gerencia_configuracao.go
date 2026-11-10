package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	gcsQuizPrefix = "seed-gerencia-configuracao-"
	gcsTitle      = "Quest 11 — Gerência de Configuração (ENADE)"
	gcsDesc       = "15 questões de 1 minuto: fundamentos de GCS (Pressman, Sommerville, IEEE 828) e itens oficiais do ENADE ADS 2011, 2017 e 2021 (baseline, versionamento, atividades e processo)."
)

func gcsPack() pack {
	return pack{
		idPrefix:  gcsQuizPrefix,
		title:     gcsTitle,
		desc:      gcsDesc,
		questions: gcsQuestions(),
	}
}

// EnsureGerenciaConfiguracaoQuiz cria o quiz de gerência de configuração se o professor ainda não o tiver.
func EnsureGerenciaConfiguracaoQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return gcsPack().ensure(ctx, st, teacherID)
}

// IsGerenciaConfiguracaoSeedQuiz reports whether a quiz was created by this seed pack.
func IsGerenciaConfiguracaoSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, gcsQuizPrefix)
}

func gcsMC(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: timeLimitOneMin,
	}
}

func gcsQuestions() []draftQuestion {
	return []draftQuestion{
		// --- ENADE oficiais (adaptadas a 4 alternativas; gabarito INEP preservado) ---
		gcsMC(
			"[ENADE 2011 · Q24] Sobre baselines na gerência de configuração, avalie:\n\nI. Representam conjuntos de itens de configuração formalmente aprovados que servem de base para as etapas seguintes.\nII. Podem ocorrer ao final de cada fase do desenvolvimento ou de outro modo definido pela gerência.\nIII. Um item “baselined” pode ser alterado a qualquer momento, sem procedimento formal.\nIV. Após a aceitação de uma etapa, a baseline pode ser descartada, pois já está nos artefatos gerados.\nV. A identificação da configuração inclui nomenclatura que permita identificar de forma inequívoca ICs, baselines e releases.\n\nÉ correto apenas o que se afirma em:",
			[]string{
				"III e V",
				"III e IV",
				"I, II e IV",
				"I, II e V",
			},
			3, // oficial E
		),
		gcsMC(
			"[ENADE 2017 · Q22] Em uma empresa, o time explica a gerência de configuração ao novo desenvolvedor. Avalie:\n\nI. Para ter acesso à base, o profissional deve fazer checkout e baixar os arquivos do servidor.\nII. Quando dois desenvolvedores alteram o mesmo documento, o conflito no commit se resolve com uma operação de tag (release).\nIII. No desenvolvimento de um novo caso de uso com vários arquivos, recomenda-se criar uma ramificação (branch).\n\nÉ correto apenas o que se afirma em:",
			[]string{
				"I e II",
				"I e III",
				"II e III",
				"apenas II",
			},
			1, // oficial B
		),
		gcsMC(
			"[ENADE 2021 · Q21] A Gestão de Configuração de Software estabelece atividades para gerenciar alterações ao longo de todo o ciclo de vida. Assinale a opção correta acerca das atividades de GCS (Pressman):",
			[]string{
				"Identificação de itens, gerenciamento de alterações, validação de versões e controle de qualidade",
				"Identificação de objetos na configuração, controle de versão, controle de alterações, auditoria de configuração e relatório de status",
				"Gerência de projeto, gerência de configuração, ciclo de vida, processo de software e evolução das configurações",
				"Construção de modelo de requisitos, métricas de projeto, métricas de controle de versão e relatório de manutenção",
			},
			1, // oficial B
		),
		gcsMC(
			"[ENADE 2021 · Q22] Sobre o processo de Gerência de Configuração de Software (GCS), avalie:\n\nI. O controle de versão permite compartilhamento colaborativo, gestão de ramos e existência simultânea de diferentes versões.\nII. A auditoria de configuração é atividade de garantia de qualidade que visa manter a qualidade do software quando há alterações.\nIII. O relatório de status serve principalmente para manter o cliente informado e deve ser gerado só ao final do processo.\nIV. A GCS só inicia depois que o software é entregue ao cliente e colocado em operação.\n\nÉ correto apenas o que se afirma em:",
			[]string{
				"I e II",
				"II e III",
				"III e IV",
				"I, II e IV",
			},
			0, // oficial A
		),

		// --- Fundamentos do material ---
		gcsMC(
			"Segundo o material (Pressman/Sommerville), a Gerência de Configuração de Software (GCS) é, principalmente:",
			[]string{
				"A gestão de prazo, custo e escopo do projeto",
				"O conjunto de atividades que identifica, controla e comunica mudanças e o estado da configuração ao longo do ciclo de vida",
				"Apenas os testes de aceitação feitos pelo cliente na entrega",
				"A elicitação de requisitos funcionais junto aos usuários",
			},
			1,
		),
		gcsMC(
			"Um Item de Configuração (IC / SCI) é:",
			[]string{
				"Qualquer rascunho informal da equipe, inclusive anotações descartáveis",
				"Somente o código-fonte compilado em produção",
				"Um artefato tratado como unidade de controle (requisitos, código, script, teste, manual, build, release etc.)",
				"O cronograma de reuniões do Change Control Board",
			},
			2,
		),
		gcsMC(
			"A identificação da configuração tem como objetivo principal:",
			[]string{
				"Definir o orçamento do projeto e o prazo de cada sprint",
				"Permitir a identificação inequívoca de qualquer versão de qualquer artefato (nomenclatura, relações, baselines e releases)",
				"Substituir a auditoria funcional pelo teste unitário",
				"Gerar automaticamente o relatório de status só no encerramento",
			},
			1,
		),
		gcsMC(
			"No controle de versão, qual operação marca uma versão estável/entregável (não resolve conflito)?",
			[]string{
				"Checkout / clone",
				"Commit / check-in",
				"Merge / update",
				"Tag / release",
			},
			3,
		),
		gcsMC(
			"Quando dois desenvolvedores alteram o mesmo trecho e há conflito no commit, a operação adequada é:",
			[]string{
				"Tag / release, para “colar” as versões",
				"Update / merge (sincronizar, integrar e resolver o conflito)",
				"Descartar a baseline imediatamente",
				"Criar um novo IC sem identificação",
			},
			1,
		),
		gcsMC(
			"Depois que um IC está sob baseline, uma alteração “ad hoc” (sem pedido, análise e aprovação) viola principalmente:",
			[]string{
				"O controle de alterações (change control)",
				"A elicitação de requisitos",
				"A multiprogramação do sistema operacional",
				"A conversão entre bases numéricas",
			},
			0,
		),
		gcsMC(
			"Sobre baseline (linha-base), assinale a afirmação correta:",
			[]string{
				"É um conjunto de ICs formalmente revisado e aprovado que serve de referência; não se descarta após a aceitação da etapa",
				"É o mesmo que tag Git usada apenas para resolver conflitos",
				"É um relatório gerado só para o cliente no final do projeto",
				"Permite alterar qualquer IC baselined sem procedimento formal",
			},
			0,
		),
		gcsMC(
			"Na auditoria de configuração, a distinção clássica é entre auditoria:",
			[]string{
				"Financeira e contábil",
				"Funcional (o item atende à especificação?) e física (o entregue corresponde ao documentado no repositório?)",
				"De prazo e de escopo",
				"De elicitação e de validação de requisitos",
			},
			1,
		),
		gcsMC(
			"O relatório de status (configuration status accounting):",
			[]string{
				"É um boletim único, gerado apenas na festa de encerramento do projeto, só para o cliente",
				"Registra e reporta continuamente o estado dos ICs, solicitações de mudança, baselines e releases",
				"Substitui o controle de versão e elimina a necessidade de commit",
				"Só existe depois que o software entra em manutenção corretiva",
			},
			1,
		),
		gcsMC(
			"Em relação ao ciclo de vida, a GCS:",
			[]string{
				"Começa somente depois que o software é fornecido ao cliente e colocado em operação",
				"Percorre todo o ciclo de vida, desde o desenvolvimento — não é atividade exclusiva de “pós-venda”",
				"Substitui a gerência de projetos em prazo e custo",
				"Só se aplica a manuais de usuário, nunca ao código-fonte",
			},
			1,
		),
		gcsMC(
			"Pressman distingue controle de versão e controle de alterações. A diferença correta é:",
			[]string{
				"Versão = governança formal da mudança (SCR, CCB); alterações = armazenamento no repositório",
				"Versão = técnica de armazenamento/recuperação (histórico, branches); alterações = processo de governança da mudança (pedido, impacto, aprovação)",
				"São sinônimos absolutos: qualquer commit já é uma baseline",
				"Versão só existe em produção; alterações só existem em requisitos",
			},
			1,
		),
	}
}
