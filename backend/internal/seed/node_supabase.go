package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	nodeSupabaseQuizPrefix = "seed-node-supabase-"
	nodeSupabaseTitle      = "Quest 4 — Node.js + Supabase: do zero ao CRUD"
	nodeSupabaseDesc       = "15 questões de 1 minuto sobre o material Node.js + Supabase: conceitos, SQL, policies, async/await e a API Express."
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

func nodeSupabaseQuestions() []draftQuestion {
	return []draftQuestion{
		{
			text: "O que é o Node.js?",
			options: []string{
				"Um framework de front-end para desenhar telas no navegador",
				"Um ambiente que executa JavaScript fora do navegador",
				"Um banco de dados relacional instalado na sua máquina",
				"Um editor de código que substitui o VS Code",
			},
			correctIndex: 1,
		},
		{
			text: "O que é o Supabase, do jeito que usamos nesta aula?",
			options: []string{
				"Um backend na nuvem construído sobre o PostgreSQL",
				"Uma biblioteca que roda o banco dentro do seu computador",
				"Um substituto do npm para instalar pacotes",
				"Um servidor HTTP que dispensa o Express",
			},
			correctIndex: 0,
		},
		{
			text: "O que cada pacote instalado abaixo faz no projeto?",
			code: `npm init -y
npm install @supabase/supabase-js dotenv express`,
			codeLanguage: "shell",
			options: []string{
				"supabase-js cria o servidor, dotenv instala o Node e express lê o .env",
				"Os três são opcionais: o Node já vem com tudo isso pronto",
				"supabase-js conversa com o Supabase, dotenv lê o .env e express cria a API HTTP",
				"supabase-js instala o PostgreSQL na sua máquina e os outros dois são atalhos",
			},
			correctIndex: 2,
		},
		{
			text: "No script da tabela, o que a linha do RA garante?",
			code: `create table if not exists alunos (
  id bigint generated always as identity primary key,
  nome text not null,
  ra text not null unique,
  curso text,
  criado_em timestamptz default now()
);`,
			codeLanguage: "sql",
			options: []string{
				"Que o RA é preenchido automaticamente pelo banco",
				"Que o RA é obrigatório e não pode se repetir na tabela",
				"Que o RA pode ficar em branco, desde que não se repita",
				"Que o RA passa a ser a chave primária no lugar do id",
			},
			correctIndex: 1,
		},
		{
			text: "O que essa policy libera na tabela alunos?",
			code: `alter table alunos enable row level security;

create policy "alunos_select_anon" on alunos
  for select to anon
  using (true);`,
			codeLanguage: "sql",
			options: []string{
				"Leitura de qualquer linha para quem usa a chave anon",
				"Todas as operações do CRUD para qualquer usuário",
				"Apenas a inserção de novas linhas pelo painel do Supabase",
				"Nada: policies só valem para usuários que fizeram login",
			},
			correctIndex: 0,
		},
		{
			text: "A anon key aparece no código e pode ser vista no navegador. O que realmente protege os dados?",
			options: []string{
				"O sigilo da chave, que funciona como uma senha",
				"O HTTPS, que impede qualquer leitura indevida da tabela",
				"O arquivo .env, que criptografa o banco de dados",
				"As policies de RLS, que decidem o que o papel anon pode fazer",
			},
			correctIndex: 3,
		},
		{
			text: "O que o console.log abaixo imprime?",
			code: `function buscar() {
  const resultado = supabase.from('alunos').select('*');
  console.log(resultado);
}`,
			codeLanguage: "javascript",
			options: []string{
				"A lista de alunos cadastrados na tabela",
				"Uma Promise pendente, e não os dados",
				"null, porque a tabela ainda está vazia",
				"Um erro de sintaxe, porque falta o await",
			},
			correctIndex: 1,
		},
		{
			text: "Para que serve o .single() no final da consulta?",
			code: `async function inserirAluno(nome, ra, curso) {
  const { data, error } = await supabase
    .from('alunos')
    .insert([{ nome, ra, curso }])
    .select()
    .single();

  if (error) throw error;
  return data;
}`,
			codeLanguage: "javascript",
			options: []string{
				"Faz o retorno vir como um único objeto, em vez de um array",
				"Garante que apenas uma linha seja inserida por vez",
				"Ignora o erro caso o RA já exista na tabela",
				"Ordena o resultado pelo id antes de devolver",
			},
			correctIndex: 0,
		},
		{
			text: "Por que essa busca usa .maybeSingle() em vez de .single()?",
			code: `async function buscarPorRa(ra) {
  const { data, error } = await supabase
    .from('alunos')
    .select('*')
    .eq('ra', ra)
    .maybeSingle();

  if (error) throw error;
  return data;
}`,
			codeLanguage: "javascript",
			options: []string{
				"Porque é mais rápido quando a tabela tem muitas linhas",
				"Porque .single() só pode ser usado depois de um insert",
				"Porque aceita 0 ou 1 resultado sem lançar erro quando o RA não existe",
				"Porque .maybeSingle() devolve sempre um array, mesmo vazio",
			},
			correctIndex: 2,
		},
		{
			text: "O que o filtro abaixo procura no banco?",
			code: `async function buscarPorNome(parteNome) {
  const { data, error } = await supabase
    .from('alunos')
    .select('*')
    .ilike('nome', '%' + parteNome + '%')
    .order('nome', { ascending: true });

  if (error) throw error;
  return data;
}`,
			codeLanguage: "javascript",
			options: []string{
				"Nomes que contenham o texto em qualquer posição, sem diferenciar maiúsculas",
				"Somente nomes que comecem exatamente com o texto informado",
				"Somente nomes idênticos ao texto, respeitando maiúsculas e minúsculas",
				"Nomes que terminem com o texto, ignorando acentos",
			},
			correctIndex: 0,
		},
		{
			text: "O que aconteceria se o .eq('ra', ra) fosse esquecido nesse update?",
			code: `async function atualizarCurso(ra, curso) {
  const { data, error } = await supabase
    .from('alunos')
    .update({ curso })
    .eq('ra', ra)
    .select();

  if (error) throw error;
  return data;
}`,
			codeLanguage: "javascript",
			options: []string{
				"O Supabase recusaria a operação automaticamente",
				"Apenas a primeira linha da tabela seria alterada",
				"Nada mudaria: o .eq() serve só para documentar a intenção",
				"O update alteraria o curso de todas as linhas da tabela",
			},
			correctIndex: 3,
		},
		{
			text: "O que acontece no POST se a linha app.use(express.json()) for removida?",
			code: `const app = express();

app.use(express.json());

app.post('/alunos', async (req, res) => {
  const { nome, ra, curso } = req.body;
  // ...
});`,
			codeLanguage: "javascript",
			options: []string{
				"req.body chega undefined, porque ninguém interpretou o JSON recebido",
				"O servidor não inicia e o terminal mostra erro na porta",
				"As rotas GET também param de responder",
				"O Express converte o body em texto e o código segue funcionando",
			},
			correctIndex: 0,
		},
		{
			text: "Na rota abaixo, o que é o :ra e de onde vem o 404?",
			code: `app.get('/alunos/:ra', async (req, res) => {
  try {
    const aluno = await buscarPorRa(req.params.ra);
    if (!aluno) return res.status(404).json({ erro: 'Aluno não encontrado' });
    res.json(aluno);
  } catch (e) {
    res.status(500).json({ erro: e.message });
  }
});`,
			codeLanguage: "javascript",
			options: []string{
				"É um comentário do Express; o 404 vem direto do Supabase",
				"É um parâmetro de rota lido em req.params.ra, e o 404 sai quando o RA não existe",
				"É o corpo enviado pelo cliente, e o 404 aparece quando o body está vazio",
				"É um caminho literal: só a URL /alunos/:ra funciona nessa rota",
			},
			correctIndex: 1,
		},
		{
			text: "Sobre o arquivo .env do projeto, qual atitude está correta?",
			code: `SUPABASE_URL=https://SEU_CODIGO.supabase.co
SUPABASE_ANON_KEY=cole_aqui_a_chave_anon_inteira
PORT=3000`,
			codeLanguage: "shell",
			options: []string{
				"Fazer commit do .env para o time inteiro ter as chaves",
				"Usar a chave service_role no front-end para evitar bloqueios de RLS",
				"Nunca versionar o .env no Git e nunca usar a service_role no front-end",
				"Colar as chaves direto no supabase.js e apagar o .env",
			},
			correctIndex: 2,
		},
		{
			text: "Para que serve o -ContentType 'application/json' neste teste?",
			code: `Invoke-RestMethod -Method Post -Uri http://localhost:3000/alunos -ContentType 'application/json' -Body '{"nome":"João Souza","ra":"2026002","curso":"ADS"}'`,
			codeLanguage: "shell",
			options: []string{
				"Avisa o servidor que o corpo é JSON, para o express.json() conseguir interpretá-lo",
				"Define a porta em que o servidor Express vai escutar",
				"Converte a resposta recebida em uma tabela do PowerShell",
				"Autentica a requisição usando a anon key do Supabase",
			},
			correctIndex: 0,
		},
	}
}
