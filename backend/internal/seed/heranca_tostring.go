package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	tostringQuizPrefix = "seed-heranca-tostring-"
	tostringTitle      = "Quest 9 — Herança: ToString() na lanchonete (fase 7)"
	tostringDesc       = "Fase 7, situações 1 a 4: 5 objetivas (1 min) e 5 dissertativas (5 min) sobre ToString de fábrica, override no pai, base.ToString() na filha e polimorfismo na lista."
)

func tostringPack() pack {
	return pack{
		idPrefix:  tostringQuizPrefix,
		title:     tostringTitle,
		desc:      tostringDesc,
		questions: tostringQuestions(),
	}
}

// EnsureToStringQuiz cria o quiz da fase 7 (ToString) se o professor ainda não o tiver.
func EnsureToStringQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return tostringPack().ensure(ctx, st, teacherID)
}

// IsToStringSeedQuiz reports whether a quiz was created by this seed pack.
func IsToStringSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, tostringQuizPrefix)
}

const tostringProdutoBase = `abstract class Produto
{
    public string Nome { get; protected set; }
    public double Preco { get; protected set; }

    protected Produto(string nome, double preco)
    {
        Nome = nome;
        Preco = preco;
    }

    public abstract void Servir();
}`

func tostringQuestions() []draftQuestion {
	return []draftQuestion{
		// Objetivas (1 minuto)
		{
			text: "Situação 1. Nenhuma classe escreveu ToString(). O que aparece na tela?",
			code: `namespace Lanchonete
{
    // Produto (abstrata, com Nome e Preco) e
    // Salgado : Produto (com Recheio), sem ToString

    class Program
    {
        static void Main()
        {
            Produto coxinha = new Salgado("Coxinha", 7.0, "frango");
            Console.WriteLine(coxinha);
        }
    }
}`,
			codeLanguage: "csharp",
			options: []string{
				"Coxinha (R$ 7)",
				"Lanchonete.Salgado",
				"Lanchonete.Produto, porque a variável é do tipo Produto",
				"Erro de compilação: Salgado não tem ToString()",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Situação 2. Por que o recheio NÃO pode ser escrito dentro do ToString() de Produto?",
			code: tostringProdutoBase[:len(tostringProdutoBase)-1] + `
    public override string ToString()
    {
        return Nome + " (R$ " + Preco + ") | recheio " + Recheio;
    }
}`,
			codeLanguage: "csharp",
			options: []string{
				"Porque ToString() só pode devolver atributos do tipo string",
				"Porque falta escrever override também na classe Salgado",
				"Porque Produto não tem Recheio: o pai não enxerga o que é só da filha, e o C# acusa que o nome não existe",
				"Porque Recheio tem private set e por isso não pode ser lido",
			},
			correctIndex: 2,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em Produto, o aluno escreveu o método SEM a palavra override. O que aparece em Console.WriteLine(coxinha)?",
			code: `// dentro de Produto
public string ToString()
{
    return Nome + " (R$ " + Preco + ")";
}

// no Main
Produto coxinha = new Salgado("Coxinha", 7.0, "frango");
Console.WriteLine(coxinha);`,
			codeLanguage: "csharp",
			options: []string{
				"Lanchonete.Salgado: sem override o método só esconde o de Object, e o WriteLine continua usando o de fábrica",
				"Coxinha (R$ 7), porque o nome do método é o mesmo",
				"Erro de compilação: todo ToString() exige override",
				"Uma linha vazia, porque o método novo não é virtual",
			},
			correctIndex: 0,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Situação 3. Produto devolve Nome + \" (R$ \" + Preco + \")\". O que aparece na tela?",
			code: `// dentro de Salgado
public override string ToString()
{
    return base.ToString() + " | recheio " + Recheio;
}

// no Main
Produto coxinha = new Salgado("Coxinha", 7.0, "frango");
Console.WriteLine(coxinha);`,
			codeLanguage: "csharp",
			options: []string{
				"Coxinha (R$ 7), porque a variável é Produto e não enxerga o recheio",
				"Lanchonete.Salgado | recheio frango",
				"Erro: uma variável do tipo Produto não pode acessar Recheio",
				"Coxinha (R$ 7) | recheio frango",
			},
			correctIndex: 3,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Situação 4. Na 3ª volta do foreach, como p está declarado e qual ToString() é executado?",
			code: `List<Produto> cardapio = new List<Produto>
{
    new Salgado("Coxinha", 7.0, "frango"),
    new Salgado("Pastel", 9.0, "carne"),
    new Doce("Bolo de pote", 10.0, "chocolate")
};

foreach (Produto p in cardapio)
{
    Console.WriteLine(p);
}`,
			codeLanguage: "csharp",
			options: []string{
				"Declarado como Doce; executa o ToString de Produto",
				"Declarado como Produto; executa o ToString de Doce, porque quem decide é o objeto real",
				"Declarado como Produto; executa o ToString de Produto, porque a lista é de Produto",
				"Declarado como Object; executa o ToString de fábrica",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},

		// Dissertativas (5 minutos)
		{
			text: "Situação 1. Sem nenhum ToString() escrito, Console.WriteLine(coxinha) mostrou \"Lanchonete.Salgado\". Explique por que saiu esse texto: de onde vem o ToString() que foi executado e como o C# chegou até ele.",
			expectedAnswer: "Console.WriteLine(coxinha) chama coxinha.ToString() sozinho. Toda classe herda de System.Object, que já tem um ToString() de fábrica. A busca começa no Salgado, sobe para Produto e, como nenhum dos dois escreveu ToString, usa o de Object, que só conhece o nome completo da classe: namespace + classe (Lanchonete.Salgado).",
			expectedAnswers: []string{
				"O WriteLine chama o ToString escondido. Ele vem de Object, o avô de todas as classes, e o de fábrica só mostra o nome da classe com o namespace.",
			},
			keyTerms: []string{"Object", "ToString", "herda", "nome da classe", "namespace"},
			criteria: `40% — diz que Console.WriteLine chama o ToString() automaticamente (chamada implícita).
30% — diz que o método vem da classe Object, herdada por toda classe (Salgado e Produto não escreveram o seu).
30% — explica que o ToString de fábrica só conhece o tipo, por isso mostra namespace + classe (Lanchonete.Salgado).
Não exigir a palavra "implícita"; vale explicar com outras palavras.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Situação 2. Escreva o ToString() da classe Produto para que a coxinha apareça como \"Coxinha (R$ 7)\" e explique para que serve a palavra override nessa linha.",
			code: tostringProdutoBase[:len(tostringProdutoBase)-1] + `
    // escreva aqui o ToString()
}`,
			codeLanguage: "csharp",
			expectedAnswer: "public override string ToString() { return Nome + \" (R$ \" + Preco + \")\"; } O override troca o ToString de fábrica, que vem de Object e é virtual, pelo da classe Produto. Sem override o método só esconde o de Object e continua aparecendo o nome da classe.",
			expectedAnswers: []string{
				"public override string ToString() { return $\"{Nome} (R$ {Preco})\"; } override = substituir o método herdado de Object.",
			},
			keyTerms: []string{"public override string ToString", "return", "Nome", "Preco", "Object"},
			criteria: `50% — código correto: public override string ToString() com return de um texto montado com Nome e Preco (concatenação ou interpolação valem).
Usar Console.WriteLine dentro do ToString em vez de return é erro grave (no máximo 20% nessa parte).
50% — explica que override substitui (sobrescreve) o ToString herdado de Object / de fábrica; mencionar virtual ou que sem override continua o nome da classe é um bônus, não obrigatório.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Na situação 2, o recheio era colado no Main: Console.WriteLine(coxinha + \" | recheio \" + coxinha.Recheio). Qual é o problema dessa solução e como a situação 3 resolve isso?",
			expectedAnswer: "Quem mostra o salgado precisa lembrar de colar o recheio em todo lugar; se esquecer, o recheio some. Além disso a variável tem que ser do tipo Salgado para enxergar Recheio. Na situação 3 o próprio Salgado faz override do ToString e cola o recheio no texto do pai (base.ToString() + recheio). O Main volta a ser só Console.WriteLine(coxinha), e funciona mesmo com variável do tipo Produto.",
			expectedAnswers: []string{
				"O Main fica responsável por um dado que é do salgado e repete código. A solução é o Salgado se apresentar sozinho com override ToString usando base.ToString().",
			},
			keyTerms: []string{"lembrar", "Salgado", "override", "base.ToString", "Main"},
			criteria: `50% — problema: depende de lembrar de colar o recheio toda vez que mostrar (repetição / pode esquecer), ou a variável precisa ser Salgado. Qualquer um desses dois argumentos bem explicado vale os 50%.
50% — solução: o Salgado escreve o próprio ToString (override) aproveitando o texto do pai com base.ToString() e acrescentando o recheio; o Main só faz Console.WriteLine(coxinha).
Dizer só "usar polimorfismo" sem explicar onde fica o ToString vale no máximo 20%.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Situação 3. Escreva o ToString() do Salgado aproveitando o texto do pai. Depois explique o que aconteceria se você escrevesse ToString() em vez de base.ToString() dentro dele, e qual a vantagem de usar base.",
			code: `class Salgado : Produto
{
    public string Recheio { get; private set; }

    // escreva aqui o ToString()
}`,
			codeLanguage: "csharp",
			expectedAnswer: "public override string ToString() { return base.ToString() + \" | recheio \" + Recheio; } Sem o base. o Salgado chamaria o próprio ToString de novo, sem parar: loop infinito (estouro de pilha). A vantagem é reaproveitar o texto do pai (DRY): se o Produto mudar o jeito de mostrar o preço, o Salgado muda junto.",
			expectedAnswers: []string{
				"override ToString com base.ToString() + recheio. Sem base vira recursão infinita (StackOverflow). base evita repetir código do pai.",
			},
			keyTerms: []string{"base.ToString", "override", "Recheio", "loop infinito", "reaproveitar"},
			criteria: `40% — código: public override string ToString() que retorna base.ToString() concatenado com o recheio.
30% — diz que sem base. o método chama a si mesmo sem parar (loop infinito, recursão infinita ou estouro de pilha / StackOverflow).
30% — vantagem: reaproveita o texto do pai sem repetir código (DRY); mudança no pai reflete na filha.
Copiar Nome e Preco de novo em vez de usar base.ToString() perde a parte do código pela metade.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Situação 4. No foreach abaixo a variável p é sempre Produto, mas a coxinha mostra o recheio e o bolo de pote mostra a cobertura. Explique por que isso acontece usando os termos tipo estático e tipo dinâmico, e diga o que muda no foreach se surgir uma classe Bebida.",
			code: `List<Produto> cardapio = new List<Produto>
{
    new Salgado("Coxinha", 7.0, "frango"),
    new Doce("Bolo de pote", 10.0, "chocolate")
};

foreach (Produto p in cardapio)
{
    Console.WriteLine(p);
}`,
			codeLanguage: "csharp",
			expectedAnswer: "É polimorfismo. O tipo estático de p é Produto (como foi declarada), mas o tipo dinâmico é o objeto real criado com new: Salgado ou Doce. Como ToString é virtual/override, o C# escolhe a versão na hora da execução olhando o objeto real (ligação dinâmica). Se surgir Bebida com o próprio ToString, ela entra na lista e o foreach não muda nada.",
			expectedAnswers: []string{
				"A lista é de Produto mas cada objeto continua sendo o que foi criado. Quem decide o ToString é o objeto real, não o tipo da variável. Uma Bebida nova entra sem mexer no foreach.",
			},
			keyTerms: []string{"polimorfismo", "tipo estático", "tipo dinâmico", "objeto real", "Bebida"},
			criteria: `40% — tipo estático = Produto (declaração da variável/lista); tipo dinâmico = classe real do objeto criado com new (Salgado, Doce).
30% — o ToString executado é o do objeto real, escolhido durante a execução (polimorfismo / ligação dinâmica / override).
30% — Bebida: o foreach não precisa mudar; basta a nova classe herdar de Produto e ter o próprio ToString.
Dizer que "a lista transforma tudo em Produto" é erro conceitual grave.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
	}
}
