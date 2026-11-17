package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	enadeRevisaoQuizPrefix = "seed-enade-revisao-"
	enadeRevisaoTitle      = "Quest 12 — Revisão ENADE (conteúdos do semestre)"
	enadeRevisaoDesc       = "15 objetivas oficiais do ENADE ADS (2 min): requisitos, análise/arquitetura, banco de dados e gerência de configuração. Enunciados e gabarito INEP; 4 alternativas (limite do sistema). Nota só por acertos."
)

func enadeRevisaoPack() pack {
	return pack{
		idPrefix:  enadeRevisaoQuizPrefix,
		title:     enadeRevisaoTitle,
		desc:      enadeRevisaoDesc,
		questions: enadeRevisaoQuestions(),
	}
}

// EnsureEnadeRevisaoQuiz cria o quiz de revisão ENADE se o professor ainda não o tiver.
func EnsureEnadeRevisaoQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return enadeRevisaoPack().ensure(ctx, st, teacherID)
}

// IsEnadeRevisaoSeedQuiz reports whether a quiz was created by this seed pack.
func IsEnadeRevisaoSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, enadeRevisaoQuizPrefix)
}

func revisaoMC(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: timeLimitTwoMin,
	}
}

func revisaoSQL(text, code string, options []string, correct int) draftQuestion {
	q := revisaoMC(text, options, correct)
	q.code = strings.TrimRight(strings.Trim(code, "\r\n"), " \t")
	q.codeLanguage = "sql"
	return q
}

func enadeRevisaoQuestions() []draftQuestion {
	return []draftQuestion{
		// ========== REQUISITOS ==========
		revisaoMC(
			"[ENADE 2021 · Q09] Qual dos itens a seguir é um requisito funcional?",
			[]string{
				"O software deve ser operacionalizado no sistema operacional Linux.",
				"O tempo de desenvolvimento do software não deve ultrapassar seis meses.",
				"O software deve emitir relatórios de compras a cada quinze dias.",
				"O tempo de resposta do software não deve ultrapassar 30 segundos.",
			},
			2, // oficial C
		),
		revisaoMC(
			"[ENADE 2017 · Q12] Sistema acadêmico de uma universidade:\n\nR1: o sistema deve permitir que cada professor realize o lançamento de notas das turmas nas quais lecionou.\nR2: o sistema deverá possibilitar seu transporte para outro sistema operacional em, no máximo, sessenta dias.\nR3: o sistema deve permitir que um estudante realize a sua matrícula nas disciplinas oferecidas em um semestre letivo.\nR4: o sistema atualiza a nota do estudante, permitindo sua visualização, em até dois segundos depois do registro do professor.\nR5: o sistema deve permitir que o auxiliar de serviços acadêmicos realize o cadastro de um estudante em não mais do que dez minutos de orientação.\n\nRepresentam descrições de requisitos não funcionais os requisitos",
			[]string{
				"R1, R2 e R3.",
				"R1, R2 e R5.",
				"R1, R3 e R4.",
				"R2, R4 e R5.",
			},
			3, // oficial D
		),
		revisaoMC(
			"[ENADE 2021 · Q18] A engenharia de requisitos inclui quatro subprocessos: estudo de viabilidade, elicitação e análise de requisitos, especificação e validação (Sommerville).\n\nUma equipe desenvolverá um software de suporte técnico. O estudo de viabilidade já foi realizado e aprovado. A equipe seguirá os subprocessos de elicitação e análise, especificação e validação.\n\nPara esses três subprocessos, quais artefatos podem ser utilizados?",
			[]string{
				"Documento de entrevista com usuários; modelo de caso de uso para os requisitos funcionais; prototipação de telas.",
				"Documento de estudo de viabilidade; modelo de caso de uso para os requisitos funcionais; prototipação de telas.",
				"Matriz de rastreabilidade; modelo de caso de uso para os requisitos não funcionais; prototipação de telas.",
				"Documento de entrevista com usuários; modelo de caso de uso para os requisitos não funcionais; matriz de rastreabilidade.",
			},
			0, // oficial A
		),
		revisaoMC(
			"[ENADE 2017 · Q28] Requisitos descrevem o que o sistema deve fazer, seus serviços e suas restrições. Descobrir, analisar, documentar e verificar esses serviços e restrições é o processo chamado de engenharia de requisitos (Sommerville).\n\nQuanto à etapa de especificação de requisitos, avalie:\nI. Identificar as expectativas e necessidades dos stakeholders.\nII. Distribuir os requisitos em categorias, explorar suas relações e classificar sua importância para os stakeholders.\nIII. Produzir um documento de especificação de requisitos compreensível para todos os stakeholders.\nIV. Examinar a especificação para assegurar que todos os requisitos foram definidos sem inconsistências.\n\nSão atribuições na etapa de especificação:",
			[]string{
				"I e III, apenas.",
				"I e IV, apenas.",
				"II e III, apenas.",
				"II e IV, apenas.",
			},
			2, // oficial C
		),
		// ========== ANÁLISE / ARQUITETURA ==========
		revisaoMC(
			"[ENADE 2021 · Q19] Arquitetura de software é uma representação que permite analisar a efetividade do projeto no entendimento dos requisitos declarados. Durante a fase de concepção da arquitetura, podem-se considerar alternativas em um estágio em que mudanças ainda são realizadas com menor esforço. (PRESSMAN, 2016, adaptado)\n\nI. Em arquiteturas orientadas a objetos, a comunicação entre os componentes do software é realizada por intermédio da troca de mensagens.\nII. As arquiteturas monolíticas consistem de um sistema dividido em pequenas partes, possibilitando que estas tenham sua manutenção, execução e evolução individual.\nIII. No padrão arquitetural Modelo-Visão-Controle (MVC), a camada de Modelo armazena as interações realizadas no Controle, podendo ser apresentados/manipulados posteriormente na Visão.\nIV. Nas arquiteturas microsserviços, o software possui componentes altamente acoplados, dificultando a manutenção.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e III.",
				"II e III.",
				"II e IV.",
				"I, III e IV.",
			},
			0, // oficial A (E era I,III,IV — removida D oficial; mantida A)
		),
		revisaoMC(
			"[ENADE 2011 · Q12] Analise as afirmações sobre a UML:\n\nI. A UML é uma metodologia para o desenvolvimento orientado a objetos, pois fornece representações gráficas e semântica para modelagem.\nII. O diagrama de casos de uso demonstra o comportamento externo do sistema, na perspectiva do usuário; é o mais abstrato, flexível e informal da UML.\nIII. Relacionamento de extensão de um caso de uso A para B significa que toda vez que A for executado ele incorporará o comportamento de B.\nIV. Os diagramas de comportamento da UML demonstram como ocorrem as trocas de mensagens entre os objetos para se atingir um objetivo.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"II e IV.",
				"III e IV.",
				"I, II e III.",
			},
			1, // oficial B
		),
		revisaoMC(
			"[ENADE 2014 · Q18] Nos anos 70, os sistemas executavam em mainframes com aplicativos escritos em linguagens estruturadas e com todas as funcionalidades em um único módulo. Posteriormente surgiram as aplicações cliente-servidor e, a partir dos anos 90, as aplicações em três camadas.\n\nConsiderando a evolução da arquitetura de software de sistemas de informação, avalie:\nI. A separação dos sistemas em três camadas lógicas torna os sistemas mais complexos, requerendo pessoal mais especializado.\nII. A separação dos sistemas em três camadas lógicas torna os sistemas mais flexíveis, permitindo que as partes possam ser alteradas de forma independente.\nIII. A separação dos sistemas em três camadas lógicas aumentou o acoplamento, dificultando a manutenção.\n\nÉ correto o que se afirma em",
			[]string{
				"II, apenas.",
				"III, apenas.",
				"I e II, apenas.",
				"I e III, apenas.",
			},
			0, // oficial A
		),
		revisaoMC(
			"[ENADE 2014 · Q17] Casos de uso podem ser organizados pela especificação de relacionamentos de generalização, inclusão e extensão. (JACOBSON; BOOCH; RUMBAUGH, 2006, adaptado)\n\nI. Para casos de uso, a generalização significa que o caso de uso filho herda o comportamento e o significado do caso de uso pai e o filho deverá acrescentar ou sobrescrever o comportamento de seu pai.\nII. Um relacionamento de inclusão significa que o caso de uso base incorpora explicitamente o comportamento de outro caso de uso. O caso de uso base poderá permanecer isolado, mas, sob certas condições, seu comportamento poderá ser incluído pelo comportamento de outro caso de uso.\nIII. Um relacionamento estendido significa que o caso de uso base incorpora implicitamente o comportamento de outro caso de uso. O caso de uso estendido nunca permanece isolado, mas é apenas instanciado como parte de alguma base maior que o estende.\nIV. Um relacionamento estendido é utilizado para modelar comportamento opcional do sistema e um subfluxo executado somente sob determinadas condições.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"I e IV.",
				"II e III.",
				"II, III e IV.",
			},
			1, // oficial B
		),

		// ========== BANCO DE DADOS ==========
		revisaoSQL(
			"[ENADE 2014 · Q11] Em uma biblioteca, um livro está emprestado quando houver registro em Empréstimo com data_dev nula. O comando SQL que lista os títulos dos livros disponíveis para empréstimo é:",
			`select titulo from livro
except
select l.titulo from emprestimo e
inner join livro l on e.livro_cod = l.liv_cod
where e.data_dev is null`,
			[]string{
				"EXCEPT com WHERE data_dev IS NULL (todos os livros menos os ainda emprestados).",
				"UNION com WHERE data_dev IS NULL.",
				"EXCEPT com WHERE data_dev IS NOT NULL.",
				"UNION com LEFT JOIN e data_dev IS NULL.",
			},
			0, // oficial A
		),
		revisaoMC(
			"[ENADE 2014 · Q22] O modelo relacional pode impor vários tipos de restrições. (ELMASRI; NAVATHE)\n\nRestrições que permitem controlar situações como “o salário de um empregado não deve exceder o salário do supervisor do empregado” e utilizam mecanismos chamados triggers (gatilhos) na sua implementação são do tipo:",
			[]string{
				"restrições de domínio.",
				"restrições de unicidade.",
				"restrições de integridade referencial.",
				"restrições de integridade semântica.",
			},
			3, // oficial E → aqui D
		),
		revisaoMC(
			"[ENADE 2017 · Q34] A partir do diálogo do Auto da Compadecida e do DER Pessoa–Cavalo (pessoa 0,1 — 1,1 cavalo; identificador de pessoa = codigo; rg é atributo descritivo), avalie:\n\nI. O Chicó e a velha poderão ser cadastrados na entidade pessoa.\nII. O Chicó e a velha poderão ter mais que um cavalo cadastrados.\nIII. O atributo rg da entidade pessoa pode ter a função de chave primária nessa entidade.\nIV. O cavalo deverá ter no mínimo uma pessoa e uma pessoa poderá ser cadastrada sem a necessidade de ter um cavalo.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e III.",
				"I e IV.",
				"II e III.",
				"I, II e IV.",
			},
			1, // oficial B
		),
		revisaoSQL(
			"[ENADE 2021 · Q17] Modelo bancário simplificado: CLIENTE (COD_CLIENTE, NOME), CONTA (NUM_CONTA, COD_CLIENTE, …), HISTORICO_MOVIMENTACAO (COD_HISTORICO, NUM_CONTA, VALOR, …).\n\nIndique a consulta que recupera o nome dos clientes e o valor médio movimentado por eles, exibindo primeiro os de maior média.",
			`SELECT CL.NOME, AVG(HM.VALOR) AS MEDIA
FROM CLIENTE CL
INNER JOIN CONTA CT ON CT.COD_CLIENTE = CL.COD_CLIENTE
INNER JOIN HISTORICO_MOVIMENTACAO HM ON HM.NUM_CONTA = CT.NUM_CONTA
GROUP BY CL.COD_CLIENTE, CL.NOME
ORDER BY MEDIA DESC;`,
			[]string{
				"Listar cada movimentação sem AVG e sem GROUP BY, ordenando do menor para o maior.",
				"Usar AVG(CL.NOME) sem GROUP BY.",
				"INNER JOIN pelas chaves corretas, GROUP BY COD_CLIENTE e NOME, ORDER BY AVG(...) DESC.",
				"Agrupar só por NOME e usar LEFT JOIN (mistura homônimos e inclui quem não movimentou).",
			},
			2, // oficial C
		),

		// ========== GERÊNCIA DE CONFIGURAÇÃO ==========
		revisaoMC(
			"[ENADE 2021 · Q21] No decorrer de um projeto de desenvolvimento de software é compreensível que mudanças venham a ocorrer. A Gestão de Configuração de Software estabelece um conjunto de atividades para gerenciar alterações através de todo o ciclo de vida de um software. (PRESSMAN, 2016, adaptado)\n\nAssinale a opção correta acerca das atividades de GCS.",
			[]string{
				"Identificação de itens na configuração, gerenciamento de alterações, validação de versões e controle de qualidade.",
				"Identificação de objetos na configuração, controle de versão, controle de alterações, auditoria de configuração e relatório de status.",
				"Gerência de projeto, gerência de configuração, ciclo de vida do projeto, processo de software e evolução das configurações.",
				"Construção de modelo de requisitos, métricas de projeto, métricas de controle de versão e relatório de manutenção.",
			},
			1, // oficial B
		),
		revisaoMC(
			"[ENADE 2021 · Q22] Sobre o processo de Gerência de Configuração de Software (GCS), avalie:\n\nI. O controle de versão possibilita compartilhamento e edição colaborativa, gestão de diferentes ramos e existência de diferentes versões de forma simultânea.\nII. A auditoria de configuração é uma atividade de garantia de qualidade do sistema que objetiva assegurar que a qualidade do software seja mantida quando feitas alterações requisitadas.\nIII. O relatório de status possui informações sobre o GCS, tem por objetivo manter o cliente informado e deve ser gerado ao final do processo.\nIV. A GCS proporciona atividades de acompanhamento e controle de mudanças iniciadas logo depois que o software for fornecido ao cliente e colocado em operação.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"II e III.",
				"III e IV.",
				"I, II e IV.",
			},
			0, // oficial A
		),
		revisaoMC(
			"[ENADE 2011 · Q24] No que diz respeito aos baselines (linhas-base) da gerência de configuração, avalie:\n\nI. As baselines representam conjuntos de itens de configuração formalmente aprovados que servem de base para as etapas seguintes de desenvolvimento.\nII. As baselines são definidas e podem ocorrer ao final de cada uma das fases do processo de desenvolvimento de software, ou de algum outro modo definido pela gerência.\nIII. Um item de configuração de software “baselined” pode ser alterado a qualquer momento durante o desenvolvimento independentemente de um procedimento formal.\nIV. Ao término de uma etapa do desenvolvimento, e após sua aceitação formal, a baseline na qual a etapa se baseou pode ser descartada pois já se encontra representada nos artefatos gerados.\nV. Uma das funções da tarefa “Identificação da Configuração” envolve a definição de uma nomenclatura que possibilite a identificação inequívoca dos itens de configuração, baselines e releases.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"III e V.",
				"III e IV.",
				"I, II e IV.",
				"I, II e V.",
			},
			3, // oficial E
		),
	}
}
