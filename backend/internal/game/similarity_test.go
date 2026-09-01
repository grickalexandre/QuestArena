package game

import (
	"testing"

	"github.com/questarena/questarena/internal/models"
)

func TestSimilarityExact(t *testing.T) {
	if got := Similarity("Herança", "heranca"); got < 0.99 {
		t.Fatalf("expected near 1, got %v", got)
	}
}

func TestSimilarityPartial(t *testing.T) {
	got := Similarity("reaproveitamento de codigo com classes pai e filha", "reaproveitamento de código usando classe pai")
	if got < 0.4 {
		t.Fatalf("expected decent similarity, got %v", got)
	}
}

func TestSimilarityDifferent(t *testing.T) {
	got := Similarity("banana", "heranca em csharp")
	if got > 0.25 {
		t.Fatalf("expected low similarity, got %v", got)
	}
}

func TestSimilarityLongerParaphrase(t *testing.T) {
	got := Similarity(
		"Herança é o reaproveitamento de código de uma classe pai para uma classe filha",
		"reaproveitamento de código usando classe pai",
	)
	if got < 0.55 {
		t.Fatalf("longer correct explanation should score high, got %v", got)
	}
}

func TestBestSimilarityPicksClosest(t *testing.T) {
	got := BestSimilarity("polimorfismo", []string{
		"reaproveitamento de código",
		"capacidade de um objeto se comportar de várias formas",
		"polimorfismo",
	})
	if got < 0.99 {
		t.Fatalf("expected exact alternative to win, got %v", got)
	}
}

func TestGradeEssayUsesAlternatives(t *testing.T) {
	q := models.Question{
		ExpectedAnswer:  "reaproveitamento de código",
		ExpectedAnswers: []string{"herança", "classe pai e classe filha"},
	}
	got := GradeEssay("Herança", q)
	if got < 0.99 {
		t.Fatalf("alternative should match, got %v", got)
	}
}

func TestGradeEssayShortParaphraseHitsKeyTerms(t *testing.T) {
	q := models.Question{
		ExpectedAnswer: "Generalização ou especialização pergunta se isto é um tipo daquilo: Aluno é uma Pessoa. Autorrelacionamento é a mesma entidade em papéis diferentes: Funcionário supervisiona Funcionário (chefe e subordinado).",
		KeyTerms:       []string{"especialização", "pessoa", "autorrelacionamento", "funcionário"},
	}
	got := GradeEssay(
		"Especialização é aluno ser um tipo de pessoa. Autorrelacionamento é um funcionário chefe de outro.",
		q,
	)
	if got < 0.7 {
		t.Fatalf("short correct answer should pass via key terms, got %v", got)
	}
}

func TestGradeEssayWrongStaysLow(t *testing.T) {
	q := models.Question{
		ExpectedAnswer: "Especialização: Aluno é um tipo de Pessoa. Autorrelacionamento: Funcionário supervisiona Funcionário.",
		KeyTerms:       []string{"especialização", "pessoa", "autorrelacionamento", "funcionário"},
	}
	got := GradeEssay("São duas tabelas quaisquer no banco", q)
	if got >= 0.45 {
		t.Fatalf("unrelated answer should stay below threshold, got %v", got)
	}
}

func TestGradeEssayMatchesCompactIds(t *testing.T) {
	q := models.Question{
		ExpectedAnswer: "pasta com PastaId e PastaPaiId, autorrelacionamento 1:N",
		KeyTerms:       []string{"autorrelacionamento", "1:N", "PastaId", "PastaPaiId"},
	}
	got := GradeEssay("e autorrelacionamento um para muitos. tabela pasta pk pasta id fk pasta pai id", q)
	if got < 0.7 {
		t.Fatalf("compact/spaced ids and 1:N synonym should match, got %v", got)
	}
}
