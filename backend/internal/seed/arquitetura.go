package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	arquiteturaQuizPrefix = "seed-analise-arq-"
	arquiteturaTitle      = "Quest 6 — Análise e Arquitetura (ENADE)"
	arquiteturaDesc       = "30 questões só de alternativas (A–D). Nenhuma dissertativa. 6 oficiais do Enade (ADS) e 24 no mesmo formato, incluindo itens que costumam cair."
)

func arquiteturaPack() pack {
	return pack{
		idPrefix:  arquiteturaQuizPrefix,
		title:     arquiteturaTitle,
		desc:      arquiteturaDesc,
		questions: arquiteturaQuestions(),
	}
}

// EnsureArquiteturaQuiz cria o quiz de análise e arquitetura se o professor ainda não o tiver.
func EnsureArquiteturaQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return arquiteturaPack().ensure(ctx, st, teacherID)
}

// IsArquiteturaSeedQuiz reports whether a quiz was created by this seed pack.
func IsArquiteturaSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, arquiteturaQuizPrefix)
}

func longMC(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: 90,
	}
}

func arquiteturaQuestions() []draftQuestion {
	return []draftQuestion{
		longMC(
			"[ENADE 2021 · Q18] A engenharia de requisitos inclui quatro subprocessos: estudo de viabilidade, elicitação, especificação e validação (Sommerville).\n\nUma equipe desenvolverá software de suporte técnico. O estudo de viabilidade já foi realizado e aprovado. A equipe seguirá elicitação, especificação e validação.\n\nPara esses três subprocessos, quais artefatos podem ser utilizados?",
			[]string{
				"Documento de entrevista com usuários; modelo de caso de uso para os requisitos funcionais; prototipação de telas.",
				"Documento de estudo de viabilidade; modelo de caso de uso para os requisitos funcionais; prototipação de telas.",
				"Matriz de rastreabilidade; modelo de caso de uso para os requisitos não funcionais; prototipação de telas.",
				"Documento de entrevista com usuários; modelo de caso de uso para os requisitos não funcionais; matriz de rastreabilidade.",
			},
			0,
		),
		longMC(
			"[ENADE 2021 · Q19] Arquitetura permite analisar cedo se o projeto atende aos requisitos, quando mudar ainda custa pouco (Pressman).\n\nI. Em arquiteturas orientadas a objetos, a comunicação entre componentes é por troca de mensagens.\nII. Arquiteturas monolíticas consistem em um sistema dividido em pequenas partes, possibilitando manutenção, execução e evolução individual.\nIII. No MVC, a camada de Modelo armazena as interações realizadas no Controle, podendo ser apresentados/manipulados posteriormente na Visão.\nIV. Nas arquiteturas microsserviços, o software possui componentes altamente acoplados, dificultando a manutenção.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e III.",
				"II e III.",
				"II e IV.",
				"I, II e IV.",
			},
			0,
		),
		longMC(
			"[ENADE 2017] Sistema acadêmico de uma universidade:\n\nR1: o sistema deve permitir que cada professor realize o lançamento de notas das turmas nas quais lecionou.\nR2: o sistema deverá possibilitar seu transporte para outro sistema operacional em, no máximo, sessenta dias.\nR3: o sistema deve permitir que um estudante realize a sua matrícula nas disciplinas oferecidas em um semestre letivo.\nR4: o sistema atualiza a nota do estudante, permitindo sua visualização, em até dois segundos depois do registro do professor.\nR5: o sistema deve permitir que o auxiliar de serviços acadêmicos realize o cadastro de um estudante em não mais do que dez minutos de orientação.\n\nRepresentam descrições de requisitos não funcionais os requisitos",
			[]string{
				"R1, R2 e R3.",
				"R1, R2 e R5.",
				"R1, R3 e R4.",
				"R2, R4 e R5.",
			},
			3,
		),
		longMC(
			"[ENADE 2017 · Q28] Requisitos descrevem o que o sistema deve fazer, serviços e restrições. Descobrir, analisar, documentar e verificar é engenharia de requisitos (Sommerville).\n\nQuanto à etapa de especificação:\nI. Identificar expectativas e necessidades dos stakeholders.\nII. Distribuir os requisitos em categorias, explorar relações e classificar importância para os stakeholders.\nIII. Produzir um documento de especificação compreensível para todos os stakeholders.\nIV. Examinar a especificação para assegurar que todos os requisitos foram definidos sem inconsistências.\n\nSão atribuições na etapa de especificação:",
			[]string{
				"I e III, apenas.",
				"I e IV, apenas.",
				"II e III, apenas.",
				"II e IV, apenas.",
			},
			2,
		),
		longMC(
			"[ENADE 2011 · Q09] Sobre técnicas de levantamento de requisitos:\n\nI. Workshop de requisitos consiste em reuniões estruturadas e delimitadas entre analistas e representantes do cliente.\nII. Cenário consiste na observação das ações do funcionário na realização de uma tarefa, para verificar os passos necessários.\nIII. As entrevistas são realizadas com os stakeholders e podem ser abertas ou fechadas.\nIV. A prototipagem é uma versão inicial do sistema, baseada em requisitos levantados em outros sistemas da organização.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"I e III.",
				"II e IV.",
				"I, III e IV.",
			},
			1,
		),
		longMC(
			"[ENADE 2011 · Q12] Analise as afirmações sobre a UML:\n\nI. A UML é uma metodologia para o desenvolvimento OO, pois fornece representações gráficas e semântica para modelagem.\nII. O diagrama de casos de uso demonstra o comportamento externo do sistema, na perspectiva do usuário; é o mais abstrato, flexível e informal da UML.\nIII. Relacionamento de extensão de um caso de uso A para B significa que toda vez que A for executado ele incorporará o comportamento de B.\nIV. Os diagramas de comportamento da UML demonstram como ocorrem as trocas de mensagens entre os objetos para se atingir um objetivo.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"II e IV.",
				"III e IV.",
				"I, II e III.",
			},
			1,
		),
		longMC(
			"Um sistema de prontuário eletrônico será usado por médicos e enfermeiros. A diretoria exige relatórios para o conselho, mas não acessa a tela clínica. A ANS regula prazos. O servidor de laboratório envia laudos via API.\n\nI. Médico e enfermeiro são atores (e usuários) do sistema clínico.\nII. A diretoria é stakeholder, mesmo sem ser ator do diagrama de casos de uso clínicos.\nIII. A ANS, por ser órgão externo, nunca pode ser stakeholder.\nIV. O sistema laboratorial pode ser modelado como ator (sistema externo).\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"I, II e IV.",
				"II e III.",
				"I, III e IV.",
			},
			1,
		),
		{
			text: "A secretaria diz: “aluno com matrícula trancada não emite histórico escolar”.\n\nEssa frase, no vocabulário da análise, é principalmente:",
			options: []string{
				"Uma regra de negócio da instituição, que o software deverá respeitar (virando RF/restrição verificável).",
				"Um estilo arquitetural cliente-servidor.",
				"Um requisito não funcional de portabilidade.",
				"Um padrão GoF (Observer).",
			},
			correctIndex: 0,
		},
		longMC(
			"Portal de clínica:\n\nI. O paciente agenda consulta informando especialidade e convênio.\nII. A confirmação do agendamento aparece em no máximo 3 segundos.\nIII. Dados de saúde só são acessíveis a profissionais autenticados, em conformidade com a LGPD.\nIV. O sistema emite lembrete por SMS na véspera.\n\nSão requisitos funcionais apenas",
			[]string{
				"I e II.",
				"I e IV.",
				"II e III.",
				"I, III e IV.",
			},
			1,
		),
		{
			text: "“Estamos construindo o produto certo?” versus “estamos construindo o produto corretamente?” (Boehm).\n\nA primeira pergunta corresponde a:",
			options: []string{
				"Validação — o sistema (ou os requisitos) atende ao que o stakeholder realmente precisa.",
				"Verificação apenas de compilação.",
				"Estudo de viabilidade econômica só.",
				"Acoplamento entre microsserviços.",
			},
			correctIndex: 0,
		},
		{
			text: "A matriz liga o RF “matricular aluno” ao caso de uso Matricular, à classe MatriculaService e ao teste T-04.\n\nIsso serve principalmente para:",
			options: []string{
				"Seguir o requisito da origem até implementação e teste, apoiando mudança e auditoria.",
				"Substituir a elicitação (não precisamos mais entrevistar).",
				"Provar que o sistema é monolítico.",
				"Definir o estilo pipes-and-filters.",
			},
			correctIndex: 0,
		},
		{
			text: "No sistema de e-commerce, “Finalizar compra” sempre executa “Calcular frete”. “Aplicar cupom” só ocorre se o cliente informar um código.\n\nA modelagem UML adequada é:",
			options: []string{
				"Finalizar compra «include» Calcular frete; Aplicar cupom «extend» Finalizar compra.",
				"Tudo é «extend», porque cupom existe.",
				"Tudo é «include», inclusive o cupom.",
				"UML não modela loja virtual.",
			},
			correctIndex: 0,
		},
		longMC(
			"[Inspirada no ENADE 2011 · Q11] Atores são bonequinhos; cada interação de valor é uma elipse nomeada.\n\nI. Essa notação descreve o diagrama de casos de uso.\nPORQUE\nII. A UML é um processo de gerenciamento de projeto (substitui o Scrum) e por isso casos de uso tornaram-se obrigatórios em todo software.\n\nAcerca dessas asserções:",
			[]string{
				"I verdadeira, II verdadeira, e II justifica I.",
				"I verdadeira, II verdadeira, mas II não justifica I.",
				"I verdadeira e II falsa.",
				"I falsa e II verdadeira.",
			},
			2,
		),
		{
			text: "O analista precisa mostrar que, ao confirmar matrícula, a tela chama o controlador, que consulta a regra de vagas e então persiste.\n\nO diagrama UML mais adequado para a ordem das mensagens no tempo é:",
			options: []string{
				"Diagrama de sequência (ou de comunicação).",
				"Somente diagrama de classes, porque tem atributos.",
				"Diagrama de implantação, porque tem servidor.",
				"Diagrama de casos de uso, porque tem tempo implícito.",
			},
			correctIndex: 0,
		},
		longMC(
			"I. Alto acoplamento entre módulos dificulta manutenção isolada.\nII. Alta coesão significa que o módulo trata de responsabilidades próximas.\nIII. Microsserviços bem desenhados buscam baixo acoplamento (comunicação por contratos/API), não “tudo no mesmo banco sem fronteira”.\nIV. Coesão alta é indesejável; o ideal é um único módulo com todas as regras da empresa.\n\nÉ correto apenas",
			[]string{
				"I e II.",
				"I, II e III.",
				"III e IV.",
				"I e IV.",
			},
			1,
		),
		{
			text: "O time discute se o sistema de streaming usará pipes-and-filters na ingestão de vídeo e, numa tela, o padrão Observer para atualizar a barra de progresso.\n\nA distinção correta é:",
			options: []string{
				"Pipes-and-filters é estilo (organização ampla do fluxo); Observer é padrão de projeto (colaboração local entre classes).",
				"Os dois são padrões GoF da mesma escala.",
				"Observer é estilo de implantação em nuvem.",
				"Pipes-and-filters é metodologia de elicitação.",
			},
			correctIndex: 0,
		},
		{
			text: "Um laboratório terá dezenas de PCs magros só com navegador; o processamento e os dados ficam num servidor (ou na nuvem). Novos postos devem entrar sem reinstalar o núcleo do sistema.\n\nA descrição mais fiel é:",
			options: []string{
				"Arquitetura cliente-servidor (clientes thin + servidor de aplicação/dados).",
				"P2P, porque cada PC é igual ao servidor.",
				"Somente blackboard de IA.",
				"Monolito é impossível nesse cenário.",
			},
			correctIndex: 0,
		},
		{
			text: "No portal da biblioteca, o aluno clica em “renovar”. A tela não calcula multa; uma classe de domínio calcula e grava; depois a tela mostra o novo prazo.\n\nMapeamento MVC:",
			options: []string{
				"Clique chega ao Controle; regra e persistência no Modelo; novo prazo na Visão.",
				"A Visão calcula a multa, porque é ela que o aluno vê.",
				"O Modelo só desenha HTML.",
				"Não há Controle em aplicações web.",
			},
			correctIndex: 0,
		},
		{
			text: "Quatro pessoas, produto de agenda para clínicas, regras de convênio ainda mudando toda semana, poucas dezenas de usuários no piloto.\n\nA escolha inicial mais razoável é:",
			options: []string{
				"Monolito modular bem fatiado internamente, com possibilidade de extrair serviços depois.",
				"Trinta microsserviços no primeiro sprint, cada tabela um serviço.",
				"Proibir API e usar só planilha na rede.",
				"Uma função serverless para cada campo do formulário.",
			},
			correctIndex: 0,
		},
		{
			text: "Há serviços de agenda, prontuário e faturamento. O app móvel não deve conhecer o endereço interno de cada um nem repetir autenticação em todos.\n\nO API Gateway serve sobretudo para:",
			options: []string{
				"Entrada única: roteamento, autenticação e políticas transversais.",
				"Substituir o banco relacional.",
				"Gerar diagrama de casos de uso automaticamente.",
				"Eliminar requisitos não funcionais.",
			},
			correctIndex: 0,
		},
		longMC(
			"I. REST organiza a API em recursos e verbos HTTP.\nII. Container empacota a aplicação e dependências de forma reproduzível.\nIII. Kubernetes orquestra muitos containers (subir, escalar, reiniciar).\nIV. Docker e Kubernetes são a mesma ferramenta; um container substitui a análise de requisitos.\n\nÉ correto apenas",
			[]string{
				"I e II.",
				"I, II e III.",
				"III e IV.",
				"II e IV.",
			},
			1,
		),
		{
			text: "A faculdade usa e-mail no navegador (sem instalar servidor de correio). O laboratório de ADS aluga VMs e instala o próprio Linux. A disciplina de nuvem publica só o código e a plataforma executa.\n\nA ordem correta é:",
			options: []string{
				"SaaS (e-mail); IaaS (VMs); PaaS (só o código).",
				"IaaS para o e-mail, porque há servidor em algum lugar.",
				"Tudo é SaaS, inclusive a VM.",
				"PaaS é o e-mail; SaaS é a VM.",
			},
			correctIndex: 0,
		},
		{
			text: "O serviço de pagamento publica “pagamento aprovado”. O serviço de matrícula escuta e libera a disciplina alguns segundos depois. Por instantes, a tela de pagamento já mostrou sucesso e a de matrícula ainda não.\n\nIsso ilustra principalmente:",
			options: []string{
				"Arquitetura orientada a eventos com consistência eventual (dados convergem, não são instantaneamente iguais em todos os serviços).",
				"Falha definitiva de RF: matrícula nunca será liberada.",
				"Modelo von Neumann quebrado.",
				"Alto acoplamento desejável entre os dois serviços via banco compartilhado sem contrato.",
			},
			correctIndex: 0,
		},
		{
			text: "App de clínicas: agenda, prontuário e convênio. Qual desenho de arquitetura é o mais adequado?",
			options: []string{
				"App + API + serviços + dados, com contratos claros, LGPD no prontuário e plano de disponibilidade.",
				"Só escolher a linguagem da moda, sem falar de qualidade nem segurança.",
				"Copiar um tutorial de framework, sem ligar ao caso da clínica.",
				"Ignorar RNF, porque sistema de saúde só tem requisito funcional.",
			},
			correctIndex: 0,
		},
		longMC(
			"Uma clínica tem recepcionista, procedimento de triagem, software de agenda, o PC da recepção e o banco de pacientes.\n\nI. O software sozinho já é o sistema computacional completo.\nII. Pessoas, processos, software, hardware e dados fazem parte do sistema.\nIII. Arquitetura de software recorta os programas; o sistema inclui o contexto organizacional.\nIV. O switch da rede e o treinamento da secretaria nunca importam para requisitos de qualidade.\n\nÉ correto apenas",
			[]string{
				"I e II.",
				"II e III.",
				"I e IV.",
				"II, III e IV.",
			},
			1,
		),
		{
			text: "A documentação OpenAPI de POST /matriculas descreve o que o portal pode pedir ao serviço acadêmico. Isso é, principalmente:",
			options: []string{
				"Interface de componente (contrato entre programas) — não confundir com a tela (UI).",
				"A camada de Visão do MVC, porque o aluno vê a matrícula.",
				"Arquitetura Harvard, pois há duas memórias.",
				"Técnica de elicitação por observação/etnografia.",
			},
			correctIndex: 0,
		},
		{
			text: "RH, acadêmico e tesouraria pertencem a donos diferentes e precisam interoperar por contratos (WSDL ou API), sem virar um único sistema compartilhado.\n\nA organização que melhor descreve esse recorte é:",
			options: []string{
				"SOA (ou APIs + eventos): serviços reutilizáveis com contrato estável.",
				"P2P, porque cada setor é simétrico ao servidor.",
				"Harvard, porque há dois tipos de memória.",
				"Apenas pipes-and-filters de compilador.",
			},
			correctIndex: 0,
		},
		{
			text: "A apresentação não chama SQL diretamente; o negócio fica no meio; a persistência embaixo. Os três andares rodam no mesmo servidor.\n\nIsso é:",
			options: []string{
				"Arquitetura em camadas (recorte lógico; pode existir dentro de um monolito).",
				"Microsserviços obrigatórios, um serviço por tabela.",
				"P2P entre as camadas.",
				"Prova de que UML é metodologia de processo.",
			},
			correctIndex: 0,
		},
		longMC(
			"I. Na arquitetura de von Neumann, programa e dados compartilham a mesma memória e o mesmo canal com a CPU (gargalo).\nII. Na arquitetura Harvard, há vias separadas para instruções e para dados — comum em microcontroladores e DSPs.\nIII. Dezenas de PCs magros + um servidor de laboratório formam uma arquitetura Harvard.\nIV. Von Neumann explica a rede cliente-servidor de um portal acadêmico.\n\nÉ correto apenas",
			[]string{
				"I e II.",
				"I e III.",
				"II e IV.",
				"I, II e III.",
			},
			0,
		),
		{
			text: "A sessão do aluno fica só na RAM do servidor. A matrícula precisa sobreviver se a máquina cair ou faltar energia.\n\nA afirmação correta é:",
			options: []string{
				"RAM é volátil; matrícula durável vai à memória secundária (SSD/banco) e a backup.",
				"RAM é persistente e substitui o banco de dados.",
				"Cache de CPU é o mesmo que CDN.",
				"A arquitetura de von Neumann resolve sozinha o backup da matrícula.",
			},
			correctIndex: 0,
		},
	}
}
