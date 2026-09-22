package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	bancoDadosQuizPrefix = "seed-banco-dados-"
	bancoDadosTitle      = "Quest 7 — Banco de Dados (ENADE)"
	bancoDadosDesc       = "30 objetivas no formato Enade (2 min cada): os itens oficiais reconstruídos do material (2011, 2017, 2005) e extras no mesmo molde (JOIN, GROUP BY, chaves, normalização, ACID, visão, índice). Enunciados parafraseados — o caderno oficial está no INEP."
)

func bancoDadosPack() pack {
	return pack{
		idPrefix:  bancoDadosQuizPrefix,
		title:     bancoDadosTitle,
		desc:      bancoDadosDesc,
		questions: bancoDadosQuestions(),
	}
}

// EnsureBancoDadosQuiz cria o quiz de banco de dados (ENADE) se o professor ainda não o tiver.
func EnsureBancoDadosQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return bancoDadosPack().ensure(ctx, st, teacherID)
}

// IsBancoDadosSeedQuiz reports whether a quiz was created by this seed pack.
func IsBancoDadosSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, bancoDadosQuizPrefix)
}

func enadeMC(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: 120,
	}
}

func enadeSQL(text, code string, options []string, correct int) draftQuestion {
	q := enadeMC(text, options, correct)
	q.code = strings.TrimRight(strings.Trim(code, "\r\n"), " \t")
	q.codeLanguage = "sql"
	return q
}

func enadeDER(text, svg string, options []string, correct int) draftQuestion {
	q := enadeMC(text, options, correct)
	q.diagram = strings.TrimSpace(svg)
	return q
}

func bancoDadosQuestions() []draftQuestion {
	return []draftQuestion{
		enadeSQL(
			"[ENADE 2011 · Q22] Sistema de malotes com três relações: FUNCIONARIOS (MATRICULA PK, NOME), MALOTES (CODIGO_MALOTE PK, MATRICULA FK, CODIGO_CONTEUDO FK, SITUACAO_MALOTE) e CONTEUDOS (CODIGO_CONTEUDO PK, DESCRICAO).\n\nO relatório pede NOME de quem enviou, CODIGO_MALOTE, DESCRICAO e SITUACAO_MALOTE. Qual consulta realiza as duas junções internas com atributos qualificados?",
			`SELECT NOME, CODIGO_MALOTE, DESCRICAO, SITUACAO_MALOTE
FROM MALOTES
INNER JOIN CONTEUDOS
  ON MALOTES.CODIGO_CONTEUDO = CONTEUDOS.CODIGO_CONTEUDO
INNER JOIN FUNCIONARIOS
  ON MALOTES.MATRICULA = FUNCIONARIOS.MATRICULA;`,
			[]string{
				"ON (CODIGO_CONTEUDO = CODIGO_CONTEUDO) e ON (MATRICULA = MATRICULA), sem prefixo de tabela",
				"FROM MALOTES, CONTEUDOS, FUNCIONARIOS WHERE com os mesmos predicados sem prefixo",
				"Dois INNER JOIN seguidos e os dois ON empilhados depois, na ordem invertida da gramática SQL",
				"Cada INNER JOIN com o ON imediatamente após o par, usando MALOTES.coluna = outra.coluna",
			},
			3,
		),
		enadeDER(
			"[ENADE 2011 · Q23] Use o DER abaixo (forte = retângulo; fraca = cantos arredondados; PK no topo). Mapeie para o relacional em 3FN.\n\nQuantas colunas implementam chaves estrangeiras (não o número de relacionamentos)?",
			derEntidades2011,
			[]string{
				"6 colunas (uma FK simples por relacionamento, mais duas extras)",
				"7 colunas (esqueceu que a PK de Ent1 tem dois atributos)",
				"8 colunas (esqueceu a PK composta da fraca na associativa com Ent4)",
				"9 colunas (3 da associativa Ent1–Ent2 + 1 em Ent3 + 2 em Ent5 + 3 da associativa Ent3–Ent4)",
			},
			3,
		),
		enadeDER(
			"[ENADE 2017 · item 34] Diálogo + DER Pessoa–Cavalo. rg é atributo descritivo, não identificador.\n\nI. As pessoas do diálogo podem ser cadastradas na entidade Pessoa.\nII. Cada uma delas pode ter mais de um cavalo cadastrado.\nIII. O atributo rg de Pessoa pode exercer o papel de chave primária.\nIV. Todo cavalo deve ter no mínimo uma pessoa; uma pessoa pode ser cadastrada sem cavalo.\n\nÉ correto apenas o que se afirma em",
			derPessoaCavalo2017,
			[]string{
				"I e II.",
				"I e IV.",
				"II e III.",
				"III e IV.",
			},
			1,
		),
		enadeMC(
			"[ENADE 2005 · Q29] T1 transfere R$ 100 da conta X para Y: lê X, X←X−100, escreve X, lê Y, Y←Y+100, escreve Y. T2, ao mesmo tempo, faz Escrita(Y) depois da leitura de Y por T1 e antes da escrita de Y por T1.\n\nQual propriedade ACID o intercalamento viola de imediato?",
			[]string{
				"Atomicidade — T1 não concluiu todas as escritas",
				"Isolamento — T2 interferiu no Y que T1 já tinha lido e ainda não tinha escrito",
				"Durabilidade — as escritas não sobreviveriam a uma queda",
				"Distributividade — as contas X e Y não podem ser atualizadas em paralelo",
			},
			1,
		),
		enadeMC(
			"[ENADE 2005 · Q72] Três transações sobre A e B. T1 bloqueia A, atualiza, desbloqueia A e só então bloqueia B. T2 e T3 pedem B e depois A. Analise:\n\nI. (T1, T2) não é serializável e há perigo de deadlock.\nII. (T1, T3) não é serializável, mas sem perigo de deadlock.\nIII. (T2, T3) é serializável e sem deadlock.\n\nÉ correto apenas o que se afirma em",
			[]string{
				"I e II.",
				"I e III.",
				"II e III.",
				"I, II e III.",
			},
			2,
		),
		enadeSQL(
			"[ENADE 2017 · reconstrução] Eleição: Candidato (id, nome, partido) e Votacao (candidato_id, uf, votos). O total de votos por partido não está numa coluna isolada.\n\nQual consulta devolve o total por partido?",
			`SELECT c.partido, SUM(v.votos) AS total
FROM candidato c
INNER JOIN votacao v ON v.candidato_id = c.id
GROUP BY c.partido;`,
			[]string{
				"SELECT partido, votos FROM candidato — o total já está no cadastro do candidato",
				"SELECT c.partido, SUM(v.votos) FROM candidato c INNER JOIN votacao v ON v.candidato_id = c.id GROUP BY c.partido",
				"SELECT c.partido, AVG(v.votos) FROM candidato c, votacao v GROUP BY c.nome",
				"SELECT SUM(votos) FROM votacao — sem JOIN nem GROUP BY por partido",
			},
			1,
		),
		enadeSQL(
			"[ENADE 2021 · reconstrução] Esquema Cliente (id, nome) e Pedido (id, cliente_id, valor). Pedido: média do valor por cliente, só quem tem pedido, da maior média para a menor.\n\nQual consulta atende?",
			`SELECT c.nome, AVG(p.valor) AS media
FROM cliente c
INNER JOIN pedido p ON p.cliente_id = c.id
GROUP BY c.nome
ORDER BY media DESC;`,
			[]string{
				"SELECT c.nome, AVG(p.valor) FROM cliente c INNER JOIN pedido p ON p.cliente_id = c.id — sem GROUP BY",
				"SELECT c.nome, AVG(p.valor) FROM cliente c LEFT JOIN pedido p ON p.id = c.id GROUP BY c.nome",
				"SELECT c.nome, AVG(p.valor) AS media FROM cliente c INNER JOIN pedido p ON p.cliente_id = c.id GROUP BY c.nome ORDER BY media DESC",
				"SELECT AVG(valor) FROM pedido ORDER BY cliente_id — a média geral já responde por cliente",
			},
			2,
		),
		enadeMC(
			"[ENADE 2017 · reconstrução] Reengenharia: sistema em Cobol com arquivos ADABAS/ISAM. A equipe migra para um SGBD relacional.\n\nO ganho principal que o item espera reconhecer é",
			[]string{
				"O programa passa a gravar planilhas na nuvem, o que dispensa catálogo e transação",
				"Independência de dados, catálogo, controle de concorrência e restrições no SGBD — o que o arquivo legado não garantia",
				"ADABAS já era relacional; a migração só troca o nome do produto",
				"ISAM implementa JOIN e GROUP BY; o SGBD só serve para tela",
			},
			1,
		),
		enadeMC(
			"No mapeamento MER → relacional, um relacionamento N:N entre Candidato e UF com atributo votos vira",
			[]string{
				"Uma coluna votos na tabela Candidato, repetida para cada estado",
				"Uma tabela associativa (Votacao) com as duas FKs e a coluna votos",
				"Duas FKs cruzadas: Candidato.uf e UF.candidato_id, sem tabela no meio",
				"Um atributo multivalorado votos na entidade UF, em 1FN",
			},
			1,
		),
		enadeMC(
			"INNER JOIN versus LEFT JOIN, no padrão cobrado em 2021: o enunciado pede a média só de clientes que têm pedido.\n\nA escolha correta é",
			[]string{
				"LEFT JOIN, para incluir clientes sem pedido com média nula",
				"INNER JOIN, porque só entram tuplas com correspondência nas duas relações",
				"Produto cartesiano (vírgula no FROM) e depois DISTINCT no nome",
				"FULL OUTER JOIN, obrigatório sempre que houver AVG",
			},
			1,
		),
		enadeMC(
			"Uma relação AlunoTemTelefones (ra, nome, telefone1, telefone2) viola a primeira forma normal porque",
			[]string{
				"RA determina nome, então falta 2FN",
				"Telefone é atributo multivalorado / repetido na mesma tupla; 1FN exige valor atômico e tabela filha",
				"Falta chave estrangeira para a operadora",
				"A ordem das colunas telefone1 e telefone2 importa no modelo relacional",
			},
			1,
		),
		enadeMC(
			"Matrícula (ra, disciplina, nome_aluno, nome_disciplina, nota). RA → nome_aluno e disciplina → nome_disciplina. A relação não está em 2FN porque",
			[]string{
				"Há atributo multivalorado, o que já quebra 1FN",
				"Nome_aluno e nome_disciplina dependem só de parte da chave composta (ra, disciplina)",
				"Nota depende de toda a chave; isso sozinho impede 2FN",
				"Toda tabela com três colunas descritivas está automaticamente em 3FN",
			},
			1,
		),
		enadeMC(
			"A terceira forma normal (3FN) impede dependência transitiva. Qual situação viola 3FN?",
			[]string{
				"Pedido (id, cliente_id, valor) com cliente_id FK — o valor depende da chave",
				"Aluno (ra, nome, cod_curso, nome_curso) em que nome_curso depende de cod_curso, não de ra",
				"Votacao (candidato_id, uf, votos) com PK composta — votos depende do par",
				"Entidade fraca ItemPedido com PK (pedido_id, nro) — o discriminador faz parte da chave",
			},
			1,
		),
		enadeMC(
			"Álgebra relacional e SQL: a expressão π NOME, CODIGO (FUNCIONARIOS ⋈ MALOTES) corresponde a",
			[]string{
				"SELECT * FROM funcionarios, malotes — produto cartesiano, sem predicado",
				"SELECT nome, codigo FROM funcionarios INNER JOIN malotes ON funcionarios.matricula = malotes.matricula",
				"SELECT nome FROM funcionarios WHERE codigo IN (SELECT codigo FROM malotes) — só σ, sem ⋈",
				"UNION das duas tabelas, porque π é união",
			},
			1,
		),
		enadeMC(
			"HAVING versus WHERE no molde Enade (média > 100 por cliente):",
			[]string{
				"WHERE AVG(valor) > 100 filtra depois do GROUP BY; HAVING filtra linhas cruas",
				"WHERE filtra tuplas antes de agregar; HAVING filtra grupos depois de AVG/SUM",
				"Os dois são sinônimos; o SGBD escolhe um",
				"HAVING só existe com COUNT(*); média usa apenas WHERE",
			},
			1,
		),
		enadeMC(
			"T1 escreve X, falha antes do COMMIT e o SGBD restaura X. A propriedade ACID em ação é",
			[]string{
				"Isolamento — outra transação leu X sujo",
				"Atomicidade — tudo ou nada; a transação incompleta é desfeita",
				"Durabilidade — o valor novo de X já estava confirmado",
				"Consistência — o invariante de negócio mudou de propósito",
			},
			1,
		),
		enadeMC(
			"Deadlock no controle de concorrência ocorre quando",
			[]string{
				"Uma transação libera A antes de pedir B (quebra 2PL), mesmo sem espera circular",
				"Há espera circular: cada transação segura um item e espera o outro ao mesmo tempo",
				"Duas transações pedem os itens na mesma ordem (B depois A)",
				"O isolamento é READ UNCOMMITTED, independentemente de lock",
			},
			1,
		),
		enadeMC(
			"Integridade referencial: apagar Partido em ano de eleição, ainda com candidatos ligados.\n\nA política alinhada ao material é",
			[]string{
				"ON DELETE CASCADE no candidato — some o histórico de votos sem aviso",
				"ON DELETE RESTRICT (ou NO ACTION) — o SGBD recusa enquanto houver candidatos",
				"Apagar o partido no arquivo e deixar candidato.partido_id órfão",
				"CHECK (partido IS NULL) na PK do partido",
			},
			1,
		),
		enadeMC(
			"Autorrelacionamento de menu (item com id_pai). Para listar cada item e o nome do menu superior, inclusive a raiz sem pai, o SQL usa",
			[]string{
				"INNER JOIN da tabela com ela mesma — a raiz some, porque não tem pai",
				"LEFT JOIN menu AS filho com menu AS pai ON filho.id_pai = pai.id, com apelidos",
				"GROUP BY id_pai sem JOIN — o nome do pai não está na mesma tupla",
				"UNION ALL de todas as linhas com SELECT pai FROM dual",
			},
			1,
		),
		enadeMC(
			"Portaria 2026 / objetos da área: quando justificar NoSQL em vez do relacional num caso de prova?",
			[]string{
				"Sempre — o Enade de ADS abandonou o modelo relacional",
				"Quando o domínio pede documentos flexíveis ou escala horizontal, justificando a perda de JOIN/ACID clássicos; OLTP com restrições rígidas continua relacional",
				"Quando o enunciado cita GROUP BY — agregação só existe em documento",
				"Quando há entidade fraca — NoSQL é o único mapeamento de fraca",
			},
			1,
		),
		enadeMC(
			"Na arquitetura ANSI/SPARC, CREATE VIEW pertence a qual nível, e CREATE INDEX a qual?",
			[]string{
				"VIEW no interno (arquivo) e INDEX no externo (tela do usuário)",
				"VIEW no nível externo (recorte por perfil) e INDEX no nível interno (acesso físico)",
				"Os dois no nível conceitual, porque ambos são SQL",
				"VIEW no conceitual e INDEX no externo — o professor não enxerga índice",
			},
			1,
		),
		enadeMC(
			"O DBA cria um índice em CPF para acelerar o login. O SELECT da aplicação não muda. Isso ilustra",
			[]string{
				"Independência lógica — o esquema conceitual foi decomposto",
				"Independência física — mudou o esquema interno sem alterar programas",
				"Quebra de isolamento — o índice trava a transação",
				"Violação de 3FN — CPF passou a ser chave estrangeira",
			},
			1,
		),
		enadeMC(
			"GRANT SELECT ON vw_notas_turma TO secretaria. Essa sentença é",
			[]string{
				"DML — altera as notas dos alunos",
				"DCL — concede privilégio; o SGBD grava a autorização no catálogo",
				"TCL — confirma a transação de matrícula",
				"DDL — cria a tabela de notas no nível interno",
			},
			1,
		),
		enadeMC(
			"UNION versus UNION ALL, na equivalência álgebra ↔ SQL do material:",
			[]string{
				"UNION ALL é a união de conjuntos (∪) e elimina repetido; UNION mantém duplicata",
				"UNION implementa ∪ (elimina repetido); UNION ALL mantém multiconjunto e não é ∪ clássica",
				"Os dois exigem GROUP BY para serem equivalentes a ⋈",
				"UNION só existe entre tabelas com PK diferente",
			},
			1,
		),
		enadeMC(
			"Divisão relacional (÷): matrícula (aluno, disciplina) ÷ lista obrigatória de disciplinas. O resultado é",
			[]string{
				"Todo aluno que cursou pelo menos uma disciplina da lista (isso é junção)",
				"Só quem cursou todas as disciplinas da lista — no MySQL, HAVING COUNT ou NOT EXISTS",
				"A tabela partida ao meio, metade dos atributos de cada lado",
				"O produto cartesiano aluno × disciplina, depois DISTINCT",
			},
			1,
		),
		enadeMC(
			"Entidade fraca ItemPedido identificada por Pedido. A chave primária correta no relacional é",
			[]string{
				"Só o número do item — a fraca tem chave própria como a forte",
				"O par (pedido_id, nro_item): PK da identificadora + discriminador; pedido_id também é FK",
				"Um GUID solto, porque fraca não pode usar FK na PK",
				"Apenas pedido_id, com nro_item como atributo descritivo sem unicidade",
			},
			1,
		),
		enadeMC(
			"CREATE VIEW vw_totais_partido AS SELECT partido, SUM(votos) … GROUP BY partido. Sobre essa visão,",
			[]string{
				"Ela materializa as linhas e substitui as tabelas Candidato e Votacao",
				"Em geral não guarda as linhas: cada consulta recalcula; com GROUP BY costuma ser só leitura",
				"Equivale a um índice no nível interno e acelera INSERT",
				"É DCL: só o DBA pode SELECT nela, nunca GRANT",
			},
			1,
		),
		enadeMC(
			"T1 já fez COMMIT da transferência. Há queda de energia. O recovery refaz as escritas confirmadas. A propriedade é",
			[]string{
				"Atomicidade — a transação ainda não tinha terminado",
				"Durabilidade — após COMMIT o efeito persiste mesmo com falha",
				"Isolamento — T2 não pode ter lido o saldo",
				"Distributividade — as contas estavam em servidores distintos",
			},
			1,
		),
		enadeMC(
			"LGPD no objeto da área (portaria 2026): minimização de dados no cadastro acadêmico significa",
			[]string{
				"Guardar o máximo de atributos “por precaução”, inclusive dados sensíveis sem finalidade",
				"Coletar só o necessário à finalidade (matrícula, nota); não exigir dado extra sem base legal",
				"Publicar RA e nota em VIEW aberta para toda a internet",
				"Trocar o SGBD relacional por planilha, que não se submete à lei",
			},
			1,
		),
		enadeMC(
			"COUNT(*) versus COUNT(coluna) no SQL de prova:",
			[]string{
				"São sempre iguais, inclusive com nulos na coluna",
				"COUNT(*) conta linhas do grupo; COUNT(coluna) ignora nulos daquela coluna",
				"COUNT(coluna) conta tabelas; COUNT(*) conta bancos",
				"COUNT(*) só funciona sem GROUP BY",
			},
			1,
		),
	}
}
