package seed

import (
	"context"
	"strings"

	"github.com/questarena/questarena/internal/store"
)

const (
	herancaQuizPrefix = "seed-heranca-conta-"
	herancaTitle      = "Quest 3 — Herança até a oficina 6.1"
	herancaDesc       = "Nova série: 20 certo ou errado com trechos de C# da aula (herança, classe abstrata, Laboratório A, Laboratório B e : base). Sem public/private/protected. 1 minuto por questão."
	timeLimitHeranca  = 60
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
			"Mimi usa o atributo Nome mesmo a classe Gato não declarando Nome: isso é herança do Animal.",
			`abstract class Animal
{
    string Nome;
    abstract void EmitirSom();
}

class Gato : Animal
{
    override void EmitirSom()
    {
        Console.WriteLine(Nome + " mia: Miau!");
    }
}

Gato mimi = new Gato();
mimi.Nome = "Mimi";
mimi.EmitirSom();`,
			true,
		),
		vf(
			"A Espada precisa copiar o campo Nome dentro da própria classe, senão o objeto fica sem nome.",
			`class Item
{
    string Nome;
    double Preco;
}

class Espada : Item
{
    int Dano;
}

Espada e = new Espada();
e.Nome = "Excalibur";
e.Dano = 12;`,
			false,
		),
		vf(
			"A declaração class Papagaio : Animal indica que Papagaio herda de Animal.",
			`abstract class Animal
{
    string Nome;
    abstract void EmitirSom();
}

class Papagaio : Animal
{
    override void EmitirSom()
    {
        Console.WriteLine(Nome + " fala: Ola!");
    }
}`,
			true,
		),
		vf(
			"Herdar de Item já faz a Espada atacar e a Poção curar, sem ninguém escrever override.",
			`class Item
{
    string Nome;
    void Usar()
    {
        Console.WriteLine("Voce usou " + Nome);
    }
}

class Pocao : Item { }
class Espada : Item { }`,
			false,
		),
		vf(
			"No trecho, i.Usar() imprime o ataque da Espada. Isso é polimorfismo: o objeto real escolhe o efeito.",
			`abstract class Item
{
    abstract void Usar();
}

class Espada : Item
{
    override void Usar()
    {
        Console.WriteLine("Causou 12 de dano");
    }
}

Item i = new Espada();
i.Usar();`,
			true,
		),
		vf(
			"Polimorfismo é só o nome chique da herança: se a classe tem os dois pontos, o comportamento já muda sozinho.",
			`class Item
{
    void Usar() { Console.WriteLine("usou"); }
}
class Pocao : Item { }
class Espada : Item { }`,
			false,
		),
		vf(
			"Como Conta é abstract, new Conta() abre uma conta genérica no banco e o programa segue.",
			`abstract class Conta
{
    string Titular;
    decimal Saldo;
    abstract void Sacar(decimal valor);
}

Conta x = new Conta();`,
			false,
		),
		vf(
			"A lista pode ser List<Animal> e guardar new Gato() e new Cachorro() ao mesmo tempo.",
			`List<Animal> zoo = new List<Animal>();
zoo.Add(new Gato());
zoo.Add(new Cachorro());`,
			true,
		),
		vf(
			"Se a classe concreta Cachorro esquecer o override de EmitirSom, o programa não compila.",
			`abstract class Animal
{
    abstract void EmitirSom();
}

class Cachorro : Animal
{
    // faltou override void EmitirSom()
}`,
			true,
		),
		vf(
			"new ContaCorrente() é inválido pelo mesmo motivo de new Conta(): as duas classes são moldes abstract.",
			`abstract class Conta
{
    abstract void Sacar(decimal valor);
}

class ContaCorrente : Conta
{
    decimal Limite;
    override void Sacar(decimal valor) { }
}

Conta c = new ContaCorrente();`,
			false,
		),
		vf(
			"No Laboratório A, Loro.EmitirSom() só funciona se a variável for Papagaio; se for Animal, o som some.",
			`Animal loro = new Papagaio();
loro.Nome = "Loro";
loro.EmitirSom();`,
			false,
		),
		vf(
			"No Laboratório A, o mesmo comando a.EmitirSom() no papagaio fala e no gato mia: polimorfismo no zoológico.",
			`List<Animal> zoo = new List<Animal>();
zoo.Add(new Papagaio());
zoo.Add(new Gato());

foreach (Animal a in zoo)
{
    a.EmitirSom();
}`,
			true,
		),
		vf(
			"No Laboratório B, com saldo 100 e limite 200, sacar 150 na conta corrente deixa o saldo em -50.",
			`ContaCorrente cc = new ContaCorrente();
cc.Titular = "Ana";
cc.Saldo = 100;
cc.Limite = 200;
cc.Sacar(150);`,
			true,
		),
		vf(
			"No Laboratório B, Limite mora na classe Conta; por isso a poupança também aceita cheque especial.",
			`abstract class Conta
{
    decimal Saldo;
    abstract void Sacar(decimal valor);
}

class ContaCorrente : Conta
{
    decimal Limite;
    override void Sacar(decimal valor) { /* saldo + limite */ }
}

class ContaPoupanca : Conta
{
    override void Sacar(decimal valor) { /* so o saldo */ }
}`,
			false,
		),
		vf(
			"Com saldo 100, a poupança recusa sacar 150, porque não há limite extra.",
			`ContaPoupanca p = new ContaPoupanca();
p.Titular = "Beto";
p.Saldo = 100;
p.Sacar(150);`,
			true,
		),
		vf(
			"O construtor tem o mesmo nome da classe, não declara tipo de retorno e é acionado pelo new.",
			`class Carro
{
    string Modelo;
    int Portas;

    Carro(string modelo, int portas)
    {
        Modelo = modelo;
        Portas = portas;
    }
}

Carro c = new Carro("Fusca", 2);`,
			true,
		),
		vf(
			"No new Carro, a tela mostra primeiro as portas e só depois o veículo, porque o construtor da filha roda antes do pai.",
			`class Veiculo
{
    Veiculo(string modelo)
    {
        Console.WriteLine("Veiculo: montando " + modelo);
    }
}

class Carro : Veiculo
{
    Carro(string modelo, int portas) : base(modelo)
    {
        Console.WriteLine("Carro: colocando " + portas + " portas");
    }
}

Carro c = new Carro("Fusca", 2);`,
			false,
		),
		vf(
			"Na linha do construtor, : base(modelo) manda só o que o pai precisa; as portas ficam no corpo do Carro.",
			`class Carro : Veiculo
{
    int Portas;

    Carro(string modelo, int portas) : base(modelo)
    {
        Portas = portas;
    }
}`,
			true,
		),
		vf(
			"Se Veiculo só tem construtor com modelo, o Carro não compila se faltar : base(modelo) na linha do construtor.",
			`class Veiculo
{
    Veiculo(string modelo)
    {
        Modelo = modelo;
    }
}

class Carro : Veiculo
{
    int Portas;

    Carro(string modelo, int portas)
    {
        Portas = portas;
    }
}`,
			true,
		),
		vf(
			": base(preco, nome) neste Item(string nome, int preco) só troca a ordem, mas o programa aceita e grava certo.",
			`class Item
{
    string Nome;
    int Preco;

    Item(string nome, int preco)
    {
        Nome = nome;
        Preco = preco;
    }
}

class Pocao : Item
{
    Pocao(string nome, int preco, int cura) : base(preco, nome)
    {
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
