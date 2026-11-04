package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	sorveteriaQuizPrefix = "seed-heranca-sorveteria-"
	sorveteriaTitle      = "Quest 10 — Herança em C#: sobremesa, sorvete e milk-shake"
	sorveteriaDesc       = "Com base no material da sorveteria 2D: 7 objetivas (1 min) e 3 dissertativas (5 min) sobre abstract, herança, base, override, sobrecarga e polimorfismo. Sem modificadores de acesso."
)

func sorveteriaPack() pack {
	return pack{
		idPrefix:  sorveteriaQuizPrefix,
		title:     sorveteriaTitle,
		desc:      sorveteriaDesc,
		questions: sorveteriaQuestions(),
	}
}

// EnsureSorveteriaQuiz cria o quiz da herança da sorveteria se o professor ainda não o tiver.
func EnsureSorveteriaQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return sorveteriaPack().ensure(ctx, st, teacherID)
}

// IsSorveteriaSeedQuiz reports whether a quiz was created by this seed pack.
func IsSorveteriaSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, sorveteriaQuizPrefix)
}

const sorveteriaSobremesa = `namespace Sorveteria
{
    abstract class Sobremesa
    {
        public string Nome;

        public Sobremesa(string nome)
        {
            Nome = nome;
        }

        public abstract void Servir();

        public override string ToString()
        {
            return Nome;
        }
    }
}`

func sorveteriaQuestions() []draftQuestion {
	return []draftQuestion{
		{
			text: "Na classe Sobremesa, o uso de abstract class significa que:",
			code: sorveteriaSobremesa,
			codeLanguage: "csharp",
			options: []string{
				"A classe pode ser instanciada com new Sobremesa(\"Copo\")",
				"A classe não pode ser herdada",
				"Não é permitido criar new Sobremesa(...); é preciso uma classe concreta (ex.: Sorvete ou MilkShake)",
				"Todos os métodos da classe devem ser estáticos",
			},
			correctIndex: 2,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em class Sorvete : Sobremesa, o símbolo : indica que:",
			code: `class Sorvete : Sobremesa
{
    // ...
}`,
			codeLanguage: "csharp",
			options: []string{
				"Sorvete e Sobremesa são a mesma classe",
				"Sorvete herda de Sobremesa",
				"Sobremesa herda de Sorvete",
				"Sorvete implementa apenas uma interface",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "No construtor public Sorvete(string nome, string sabor) : base(nome), a expressão : base(nome):",
			code: `class Sorvete : Sobremesa
{
    public string Sabor;

    public Sorvete(string nome, string sabor) : base(nome)
    {
        Sabor = sabor;
    }
}`,
			codeLanguage: "csharp",
			options: []string{
				"Cria um segundo objeto Sobremesa separado",
				"Chama o construtor da classe base passando o nome",
				"Substitui o valor de Sabor pelo de Nome",
				"Impede a herança de Sobremesa",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em Sobremesa, o método public abstract void Servir(); obriga que:",
			code: `abstract class Sobremesa
{
    public abstract void Servir();
}`,
			codeLanguage: "csharp",
			options: []string{
				"A classe base implemente o corpo de Servir",
				"Nenhuma classe filha possa ter Servir",
				"Toda classe concreta que herdar de Sobremesa implemente Servir",
				"Servir seja sempre estático",
			},
			correctIndex: 2,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em Sorvete, existem Servir() com override e Servir(int bolas). Essa diferença representa:",
			code: `class Sorvete : Sobremesa
{
    public override void Servir()
    {
        Console.WriteLine("Na taça: " + Nome + " de " + Sabor);
    }

    public void Servir(int bolas)
    {
        Console.WriteLine(bolas + " bolas de " + Sabor);
    }
}`,
			codeLanguage: "csharp",
			options: []string{
				"Dois construtores da classe",
				"Sobrescrita (override) e sobrecarga (outro método com parâmetros diferentes)",
				"Dois namespaces distintos",
				"Dois tipos de herança múltipla",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em MilkShake, o método Servir() chama base.Servir() e depois acrescenta outra mensagem. Isso significa que:",
			code: `class MilkShake : Sorvete
{
    public override void Servir()
    {
        base.Servir();
        Console.WriteLine("Com canudo e chantilly!");
    }
}`,
			codeLanguage: "csharp",
			options: []string{
				"O Servir da classe base é ignorado",
				"Só o método da neta é executado, sem usar a base",
				"Primeiro executa o Servir herdado/base e, em seguida, o comportamento específico do MilkShake",
				"base.Servir() apaga o método da classe filha",
			},
			correctIndex: 2,
			timeLimitSec: timeLimitOneMin,
		},
		{
			text: "Em ToString() do Sorvete, a expressão base.ToString() + \" de \" + Sabor serve para:",
			code: `public override string ToString()
{
    return base.ToString() + " de " + Sabor;
}`,
			codeLanguage: "csharp",
			options: []string{
				"Ignorar o nome definido em Sobremesa",
				"Reaproveitar o ToString da classe base e complementar com o sabor",
				"Converter Sabor em número inteiro",
				"Impedir a sobrescrita de ToString nas classes filhas",
			},
			correctIndex: 1,
			timeLimitSec: timeLimitOneMin,
		},

		{
			text: "Explique a cadeia de herança MilkShake : Sorvete : Sobremesa. Em sua resposta, diga o que um objeto MilkShake “é” em relação a Sorvete e a Sobremesa, e cite um método que existe só em MilkShake.",
			code: `class MilkShake : Sorvete
{
    public void Bater()
    {
        Console.WriteLine("Bateu o " + Nome + " no liquidificador.");
    }
}`,
			codeLanguage: "csharp",
			expectedAnswer: "MilkShake herda de Sorvete, que herda de Sobremesa. Assim, um milk-shake é um sorvete e, por isso, também é uma sobremesa. O método Bater() existe apenas em MilkShake.",
			expectedAnswers: []string{
				"MilkShake é um Sorvete e também uma Sobremesa pela cadeia de herança. Bater() só existe no MilkShake.",
			},
			keyTerms: []string{"MilkShake", "Sorvete", "Sobremesa", "herda", "Bater"},
			criteria: `40% — descreve a cadeia MilkShake herda de Sorvete, que herda de Sobremesa (ou equivalente).
30% — afirma que o milk-shake é um sorvete e também uma sobremesa.
30% — cita Bater() (ou equivalente) como método exclusivo do MilkShake.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Por que é possível guardar um Sorvete e um MilkShake em uma List<Sobremesa>? O que acontece, em termos de polimorfismo, quando se chama item.Servir() para cada elemento dessa lista?",
			code: `List<Sobremesa> cardapio = new List<Sobremesa>();
cardapio.Add(new Sorvete("Copo", "napolitano"));
cardapio.Add(new MilkShake("Shake", "morango"));

foreach (Sobremesa item in cardapio)
{
    item.Servir();
}`,
			codeLanguage: "csharp",
			expectedAnswer: "Como Sorvete e MilkShake herdam (direta ou indiretamente) de Sobremesa, ambos podem ser tratados como Sobremesa. Em tempo de execução, item.Servir() chama a implementação da classe real do objeto (Sorvete ou MilkShake).",
			expectedAnswers: []string{
				"A lista aceita porque ambos são Sobremesa. O Servir chamado é o do objeto real (polimorfismo).",
			},
			keyTerms: []string{"List", "Sobremesa", "polimorfismo", "objeto real", "Servir"},
			criteria: `50% — explica que entram na lista porque herdam de Sobremesa (são sobremesas).
50% — explica polimorfismo: Servir executado é o da classe real do objeto, não um único da base.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
		{
			text: "Explique o papel da verificação if (item is MilkShake batido) ao percorrer a lista de Sobremesa. Por que Bater() não pode ser chamado diretamente em todo item da lista?",
			code: `foreach (Sobremesa item in cardapio)
{
    item.Servir();
    if (item is MilkShake batido)
    {
        batido.Bater();
    }
}`,
			codeLanguage: "csharp",
			expectedAnswer: "A lista enxerga cada elemento como Sobremesa. Bater() existe só em MilkShake, então não está disponível para qualquer Sobremesa. O is MilkShake confirma se o objeto é um milk-shake e permite chamar Bater() com segurança apenas nesse caso.",
			expectedAnswers: []string{
				"Bater só existe no MilkShake. O is testa o tipo real antes de chamar, porque a variável da lista é Sobremesa.",
			},
			keyTerms: []string{"is", "MilkShake", "Bater", "Sobremesa", "tipo"},
			criteria: `50% — diz que Bater() não está em Sobremesa / só existe em MilkShake, por isso não dá para chamar em todo item.
50% — explica que is (ou verificação de tipo) confirma o objeto real antes de chamar Bater com segurança.`,
			threshold:    0.5,
			timeLimitSec: timeLimitFiveMin,
		},
	}
}
