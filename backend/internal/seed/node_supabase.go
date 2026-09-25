package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	nodeSupabaseQuizPrefix = "seed-node-supabase-"
	nodeSupabaseTitle      = "Quest 4 — Node.js + Supabase: do zero ao CRUD"
	nodeSupabaseDesc       = "30 questões de 1 min 30 s sobre o material Node.js + Supabase: conceitos, SQL, policies, async/await, API Express e a tela. Sair da tela perde a questão."
	timeLimitNodeSupabase  = 90
)

func nodeSupabasePack() pack {
	return pack{
		idPrefix:  nodeSupabaseQuizPrefix,
		title:     nodeSupabaseTitle,
		desc:      nodeSupabaseDesc,
		questions: nodeSupabaseQuestions(),
	}
}

// EnsureNodeSupabaseQuiz cria o quiz do material Node.js + Supabase se o professor ainda não o tiver.
func EnsureNodeSupabaseQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return nodeSupabasePack().ensure(ctx, st, teacherID)
}

// IsNodeSupabaseSeedQuiz reports whether a quiz was created by this seed pack.
func IsNodeSupabaseSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, nodeSupabaseQuizPrefix)
}

func nq(text string, options []string, correct int) draftQuestion {
	return draftQuestion{
		text:         text,
		options:      options,
		correctIndex: correct,
		timeLimitSec: timeLimitNodeSupabase,
	}
}

func nqCode(text, code, lang string, options []string, correct int) draftQuestion {
	q := nq(text, options, correct)
	q.code = code
	q.codeLanguage = lang
	return q
}

func nodeSupabaseQuestions() []draftQuestion {
	return []draftQuestion{
		nq(
			"No material, o que é o Node.js?",
			[]string{
				"Um framework de front-end para desenhar telas no navegador",
				"O ambiente que executa JavaScript fora do navegador, no computador ou em um servidor",
				"O banco PostgreSQL instalado na sua máquina",
				"O editor que substitui o VS Code",
			},
			1,
		),
		nq(
			"Para que serve o npm, que vem com o Node.js?",
			[]string{
				"Baixar e organizar bibliotecas: npm install nome-da-biblioteca",
				"Criar a tabela alunos no painel do Supabase",
				"Ligar o Row Level Security da tabela",
				"Gerar a senha do banco na hora de criar o projeto",
			},
			0,
		),
		nqCode(
			"Depois de instalar o Node no Windows, o que estes comandos confirmam?",
			"node -v\nnpm -v",
			"shell",
			[]string{
				"Que a anon key foi copiada para o .env",
				"Que a tabela alunos já existe no Supabase",
				"Que o Node e o npm estão no PATH; se o comando não existir, falta reinstalar com Add to PATH e reabrir o terminal",
				"Que o servidor Express já está escutando na porta 3000",
			},
			2,
		),
		nq(
			"O que é o Supabase, do jeito que a aula usa?",
			[]string{
				"Um backend na nuvem sobre PostgreSQL; o foco da aula é banco de dados + cliente JavaScript",
				"Um banco que roda só dentro do seu computador",
				"Um substituto do npm para instalar pacotes",
				"O servidor HTTP que dispensa o Express",
			},
			0,
		),
		nq(
			"No CRUD do material, a letra U corresponde a qual comando SQL?",
			[]string{
				"INSERT — cadastrar um registro novo",
				"SELECT — consultar registros",
				"UPDATE — alterar um registro existente",
				"DELETE — apagar um registro",
			},
			2,
		),
		nq(
			"Na arquitetura do exercício, por onde o dado passa até a tabela alunos?",
			[]string{
				"O PostgreSQL roda na sua máquina e o Node só abre o arquivo .sql",
				"Node.js (index.js ou server.js) → @supabase/supabase-js → HTTPS → Supabase Cloud",
				"O navegador grava direto no Table Editor, sem biblioteca",
				"O Express substitui o Supabase e guarda os alunos na memória",
			},
			1,
		),
		nq(
			"Ao criar o projeto no Supabase, por que a aula pede a região mais próxima?",
			[]string{
				"Porque só South America aceita a chave anon",
				"Porque a região escolhe o nome das colunas da tabela",
				"Porque a senha do banco só vale nessa região",
				"Porque uma região mais próxima reduz a latência",
			},
			3,
		),
		nq(
			"O que o material diz sobre a chave service_role?",
			[]string{
				"É chave de admin e nunca deve ir para o front-end",
				"É a chave que colocamos em SUPABASE_ANON_KEY",
				"É a senha do Database Password, anotada no .env",
				"É o papel anon dentro do Postgres",
			},
			0,
		),
		nq(
			"A anon key pode aparecer no código e ser vista no navegador. O que realmente decide o que ela pode fazer?",
			[]string{
				"O sigilo da chave, que funciona como senha",
				"O arquivo .env, que criptografa o banco",
				"As policies de RLS: a chave é o crachá de visitante; as regras da portaria é que abrem a porta",
				"O HTTPS, que impede qualquer leitura da tabela",
			},
			2,
		),
		nq(
			"A palavra anon aparece em dois lugares. Qual leitura o material faz?",
			[]string{
				"São a mesma coisa: a chave do .env já é a policy",
				"A chave no .env identifica o visitante; o Postgres executa a query como o papel (role) anon, e a policy to anon decide o que libera",
				"anon no .env é o service_role; anon na policy é o usuário logado",
				"O papel anon só existe depois que alguém faz login",
			},
			1,
		),
		nqCode(
			"No script da tabela, o que a linha do RA garante?",
			`ra text not null unique`,
			"sql",
			[]string{
				"Que o RA é preenchido automaticamente pelo banco",
				"Que o RA é obrigatório e não pode se repetir",
				"Que o RA pode ficar em branco, desde que não se repita",
				"Que o RA passa a ser a chave primária no lugar do id",
			},
			1,
		),
		nqCode(
			"O que generated always as identity faz na coluna id?",
			`id bigint generated always as identity primary key`,
			"sql",
			[]string{
				"Obriga o cliente a enviar o id em todo INSERT",
				"Impede que a tabela tenha chave primária",
				"O próprio banco gera 1, 2, 3… (o equivalente ao AUTO_INCREMENT)",
				"Preenche a data e a hora do cadastro",
			},
			2,
		),
		nq(
			"Com o RLS ligado e nenhuma policy criada, o que o papel anon consegue fazer?",
			[]string{
				"SELECT, INSERT, UPDATE e DELETE, porque a chave do .env está correta",
				"Só SELECT; gravar continua bloqueado",
				"Só o que o Table Editor permitir, ignorando o Node",
				"Nada: não lê, não grava, não altera e não apaga",
			},
			3,
		),
		nqCode(
			"Nesta policy de INSERT, qual cláusula o material usa, e o que ela olha?",
			`create policy "alunos_insert_anon" on alunos
  for insert to anon
  with check (true);`,
			"sql",
			[]string{
				"WITH CHECK: olha o dado novo que está entrando",
				"USING: olha as linhas que já estão gravadas",
				"USING e WITH CHECK juntos, como no UPDATE",
				"Nenhuma: INSERT não precisa de condição",
			},
			0,
		),
		nq(
			"No policies.sql da aula, por que o UPDATE usa USING e WITH CHECK?",
			[]string{
				"USING escolhe a tabela; WITH CHECK escolhe o papel anon",
				"USING pergunta se posso alterar esta linha; WITH CHECK pergunta se o valor final é permitido",
				"Os dois são sinônimos; o script repete por segurança",
				"USING vale para INSERT; WITH CHECK vale só para DELETE",
			},
			1,
		),
		nq(
			"Qual é a ordem correta que o checklist do material pede?",
			[]string{
				"npm start, depois policies.sql, depois tabela_alunos.sql",
				"policies.sql primeiro, porque o RLS cria a tabela",
				"tabela_alunos.sql, depois policies.sql, e só então npm run demo ou npm start",
				"Só o Table Editor; os arquivos .sql são opcionais",
			},
			2,
		),
		nq(
			"O que o material avisa sobre using (true) e with check (true)?",
			[]string{
				"É o modo certo para um app real com dados sensíveis",
				"Bloqueia a anon key, então ninguém de fora lê a tabela",
				"Só libera SELECT; INSERT continua negado",
				"Qualquer pessoa com a chave anon pode ler e apagar todos os alunos; serve para estudar, não para produção",
			},
			3,
		),
		nq(
			"O Node responde violates row-level security policy. O que a aula manda fazer?",
			[]string{
				"Rodar de novo o policies.sql: a policy falta ou o RLS está sem regra",
				"Trocar a anon key pela service_role no front-end",
				"Apagar a tabela e criar outra com outro nome",
				"Tirar o await da função, porque a Promise ainda está pendente",
			},
			0,
		),
		nqCode(
			"O que cada pacote instalado abaixo faz no projeto?",
			`npm install @supabase/supabase-js dotenv express`,
			"shell",
			[]string{
				"supabase-js cria o servidor, dotenv instala o Node e express lê o .env",
				"Os três são opcionais: o Node já vem com tudo isso pronto",
				"supabase-js conversa com o Supabase, dotenv lê o .env e express cria a API HTTP",
				"supabase-js instala o PostgreSQL na máquina e os outros dois são atalhos",
			},
			2,
		),
		nq(
			"O que require('dotenv').config() faz, e o que não pode ir para o GitHub?",
			[]string{
				"Sobe o Express na porta 3000; o package.json é que não pode ser commitado",
				"Lê o .env e coloca cada linha em process.env; o .env entra no .gitignore e não vai para o GitHub",
				"Cria as policies de RLS a partir do .env",
				"Rotaciona a anon key automaticamente a cada npm start",
			},
			1,
		),
		nqCode(
			"O que este código imprime, e por quê?",
			`function buscar() {
  const resultado = supabase.from('alunos').select('*');
  console.log(resultado);
}`,
			"javascript",
			[]string{
				"Uma Promise pendente: faltou await, então o código não esperou os dados",
				"O array de alunos, porque select('*') já é síncrono",
				"null, porque a função não é async",
				"O erro violates row-level security, sempre",
			},
			0,
		),
		nqCode(
			"Para que serve o .single() no final do insert?",
			`const { data, error } = await supabase
  .from('alunos')
  .insert([{ nome, ra, curso }])
  .select()
  .single();`,
			"javascript",
			[]string{
				"Garante que só uma linha seja inserida por vez",
				"Devolve um único objeto, em vez de um array",
				"Ignora o erro se o RA já existir",
				"Ordena o resultado pelo id",
			},
			1,
		),
		nqCode(
			"Na busca por RA, por que a aula usa maybeSingle() e não single()?",
			`.eq('ra', ra)
.maybeSingle();`,
			"javascript",
			[]string{
				"maybeSingle ordena pelo nome; single não ordena",
				"maybeSingle ignora o RLS; single exige policy",
				"Os dois lançam erro quando o RA não existe",
				"maybeSingle aceita 0 ou 1 linha e devolve null se não achar; single lançaria erro se não houvesse linha",
			},
			3,
		),
		nqCode(
			"O que aconteceria se o .eq('ra', ra) fosse esquecido nesse update?",
			`const { data, error } = await supabase
  .from('alunos')
  .update(campos)
  .eq('ra', ra)
  .select();`,
			"javascript",
			[]string{
				"O Supabase recusaria a operação automaticamente",
				"Apenas a primeira linha da tabela seria alterada",
				"O update alteraria todas as linhas da tabela",
				"Nada mudaria: o .eq() só documenta a intenção",
			},
			2,
		),
		nqCode(
			"O que o filtro abaixo procura?",
			`.ilike('nome', '%' + parteNome + '%')`,
			"javascript",
			[]string{
				"Nomes que contenham o texto em qualquer posição, sem diferenciar maiúsculas",
				"Somente nomes que comecem exatamente com o texto",
				"Somente nomes idênticos, respeitando maiúsculas",
				"Nomes que terminem com o texto, e só se tiverem acento",
			},
			0,
		),
		nqCode(
			"Sem a linha app.use(express.json()), o que chega em req.body no POST?",
			`app.use(express.json());

app.post('/alunos', async (req, res) => {
  const { nome, ra, curso } = req.body;
});`,
			"javascript",
			[]string{
				"O JSON já convertido, porque o Express lê o body sozinho",
				"Uma string com o texto cru do PowerShell",
				"O parâmetro :ra da URL",
				"undefined: o middleware é quem transforma o JSON em objeto",
			},
			3,
		),
		nqCode(
			"Depois de um POST bem-sucedido, a API responde assim. O que o 201 significa no material?",
			`res.status(201).json(aluno);`,
			"javascript",
			[]string{
				"Bad Request: faltou nome ou ra",
				"Created: o aluno foi criado",
				"Not Found: o RA não existe",
				"Internal Server Error: o banco quebrou",
			},
			1,
		),
		nq(
			"No navegador, http://localhost:3000/alunos mostra []. O que a aula diz?",
			[]string{
				"Não é erro: a tabela está vazia (não rodou npm run demo, ou os dados foram apagados)",
				"A porta 3000 está ocupada (EADDRINUSE)",
				"Falta o Content-Type, então o Express não montou o array",
				"O service_role foi parar no front-end",
			},
			0,
		),
		nq(
			"No frontend da aula, o que é o DOM?",
			[]string{
				"O banco PostgreSQL visto pelo Table Editor",
				"A árvore de elementos da página que o JavaScript lê e muda, por exemplo com getElementById",
				"O arquivo .env com a URL e a anon key",
				"O middleware que lê o JSON do POST",
			},
			1,
		),
		nqCode(
			"Para que serve evento.preventDefault() ao salvar o aluno?",
			`async function salvarAluno(evento) {
  evento.preventDefault();
}`,
			"javascript",
			[]string{
				"Impede o RLS de bloquear o INSERT",
				"Converte o formulário em JSON antes do fetch",
				"Cancela o recarregar da página, que é o comportamento padrão do formulário",
				"Fecha o modal e chama listarAlunos",
			},
			2,
		),
	}
}
