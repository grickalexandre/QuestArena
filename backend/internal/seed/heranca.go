package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	herancaQuizPrefix = "seed-heranca-conta-"
	herancaTitle      = "Quest 3 — Herança até a oficina 6.1"
	herancaDesc       = "15 certo ou errado com trechos de C# da aula: herança, classe abstrata, Laboratório A (Animal), Laboratório B (Conta) e construtor com : base. Sem public/private/protected. Tempo folgado — a nota vale pelo acerto, não pela corrida de XP."
	timeLimitHeranca  = 180
)

func herancaPack() pack {
	return pack{
		idPrefix:  herancaQuizPrefix,
		title:     herancaTitle,
		desc:      herancaDesc,
		questions: herancaQuestions(),
	}
}

// EnsureHerancaQuiz cria o quiz da aula (até a oficina 6.1) se o professor ainda não o tiver.
func EnsureHerancaQuiz(ctx context.Context, st store.Store, teacherID string) error {
	return herancaPack().ensure(ctx, st, teacherID)
}

func vf(text, code string, certa bool) draftQuestion {
	opts := []string{"Certo", "Errado"}
	idx := 0
	if !certa {
		idx = 1
	}
	// Alterna a ordem dos botões para o colega do lado não copiar “o da esquerda”.
	if len(text)%2 == 0 {
		opts[0], opts[1] = opts[1], opts[0]
		idx = 1 - idx
	}
	return draftQuestion{
		text:         text,
		options:      opts,
		correctIndex: idx,
		code:         strings.TrimRight(strings.Trim(code, "\r\n"), " \t"),
		codeLanguage: "csharp",
		timeLimitSec: timeLimitHeranca,
	}
}

func herancaQuestions() []draftQuestion {
	return []draftQuestion{
		vf(
			"Rex consegue usar Dormir() mesmo a classe Cachorro não tendo esse método escrito nela: ele vem da herança do Animal.",
			`abstract class Animal
{
    string Nome;

    void Dormir()
    {
        Console.WriteLine(Nome + " esta dormindo... zzz");
    }

    abstract void EmitirSom();
}

class Cachorro : Animal
{
    override void EmitirSom()
    {
        Console.WriteLine(Nome + " late: Au au!");
    }
}

Cachorro rex = new Cachorro();
rex.Nome = "Rex";
rex.Dormir();`,
			true,
		),
		vf(
			"Como Pocao só herda atributos, o método Vender() precisa ser copiado de novo dentro da classe Pocao.",
			`class Item
{
    string Nome;
    double Preco;

    void Vender()
    {
        Console.WriteLine("Vendeu " + Nome);
    }
}

class Pocao : Item
{
    int Cura;
}`,
			false,
		),
		vf(
			"A linha class ContaCorrente : Conta diz que a corrente herda de Conta — é um tipo de conta.",
			`abstract class Conta
{
    string Titular;
    decimal Saldo;
    abstract void Sacar(decimal valor);
}

class ContaCorrente : Conta
{
    decimal Limite;
    override void Sacar(decimal valor) { /* usa saldo + limite */ }
}`,
			true,
		),
		vf(
			"Polimorfismo e herança são sinônimos: basta herdar para o Usar() já ter efeito diferente em cada filha.",
			`class Item
{
    string Nome;
    void Usar()
    {
        Console.WriteLine("Voce usou " + Nome + "!");
    }
}

class Pocao : Item { }
class Espada : Item { }

Pocao p = new Pocao();
p.Nome = "Pocao de Vida";
Espada e = new Espada();
e.Nome = "Excalibur";
p.Usar();
e.Usar();`,
			false,
		),
		vf(
			"Depois destas linhas, item.Usar() executa o Usar da Poção (cura), porque o objeto real é Poção.",
			`abstract class Item
{
    abstract void Usar();
}

class Pocao : Item
{
    override void Usar()
    {
        Console.WriteLine("Curou 20 de vida");
    }
}

Item item = new Pocao();
item.Usar();`,
			true,
		),
		vf(
			"Como a variável foi declarada como Item, o C# ignora o new Pocao e sempre roda só o Usar do molde Item.",
			`Item item = new Pocao("Cura", 10);
item.Usar();`,
			false,
		),
		vf(
			"Se Item for abstract, o comando new Item() cria um item genérico vazio e o programa segue.",
			`abstract class Item
{
    string Nome;
    abstract void Usar();
}

Item x = new Item();`,
			false,
		),
		vf(
			"Pode existir uma variável do tipo Item apontando para new Pocao(): o rótulo é o molde, o objeto é concreto.",
			`abstract class Item { }
class Pocao : Item { }

Item pocao = new Pocao();
// Item vazio = new Item();  // nao compila`,
			true,
		),
		vf(
			"Um método abstract no pai (Usar, EmitirSom ou Sacar) obriga cada filha concreta a escrever o próprio com override.",
			`abstract class Animal
{
    abstract void EmitirSom();
}

class Cachorro : Animal
{
    override void EmitirSom()
    {
        Console.WriteLine("Au au!");
    }
}

class Gato : Animal
{
    override void EmitirSom()
    {
        Console.WriteLine("Miau!");
    }
}`,
			true,
		),
		vf(
			"No Laboratório A, zoologico.Add(new Animal()) funciona, porque a lista é do tipo Animal.",
			`abstract class Animal
{
    abstract void EmitirSom();
}

List<Animal> zoologico = new List<Animal>();
zoologico.Add(new Animal());`,
			false,
		),
		vf(
			"No foreach do zoológico, a.EmitirSom() late ou mia conforme o objeto real, sem converter o tipo na hora de emitir o som.",
			`List<Animal> zoologico = new List<Animal>();
zoologico.Add(new Cachorro());
zoologico.Add(new Gato());

foreach (Animal a in zoologico)
{
    a.EmitirSom();
}`,
			true,
		),
		vf(
			"No trecho da conta corrente, o saque de 200 é aceito e o saldo fica -150, porque o disponível é saldo + limite.",
			`ContaCorrente cc = new ContaCorrente();
cc.Titular = "Ana";
cc.Saldo = 50;
cc.Limite = 200;
cc.Sacar(200);`,
			true,
		),
		vf(
			"No Laboratório B, a poupança usa a mesma regra de saque da corrente, porque Sacar está declarado no pai Conta.",
			`abstract class Conta
{
    decimal Saldo;
    abstract void Sacar(decimal valor);
}

Conta poupanca = new ContaPoupanca();
poupanca.Saldo = 50;
poupanca.Sacar(200);`,
			false,
		),
		vf(
			"Ao executar new Pocao(...), o construtor do Item roda primeiro (por causa do : base) e só depois o da Poção.",
			`class Item
{
    Item(string nome, double preco)
    {
        Console.WriteLine("Item nascendo");
    }
}

class Pocao : Item
{
    Pocao(string nome, double preco, int cura)
        : base(nome, preco)
    {
        Console.WriteLine("Pocao nascendo");
    }
}

Pocao p = new Pocao("Cura", 10, 20);`,
			true,
		),
		vf(
			"O : base(nome, preco) pode ser escrito dentro das chaves do construtor da filha, na primeira linha, como se fosse um método.",
			`class Pocao : Item
{
    Pocao(string nome, double preco, int cura)
    {
        base(nome, preco);
        Cura = cura;
    }
}`,
			false,
		),
	}
}

// IsHerancaSeedQuiz reports whether a quiz was created by this seed pack.
func IsHerancaSeedQuiz(quizID string) bool {
	return strings.HasPrefix(quizID, herancaQuizPrefix)
}
