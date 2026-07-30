# QuestArena

Plataforma estilo Kahoot: professores criam quizzes (peso + tempo por questão) e alunos competem ao vivo com PIN e ranking de XP.

## Stack

- **Backend:** Go (REST + WebSocket autoritativo)
- **Persistência/Auth:** Firebase Auth + Firestore (ou modo DEV em memória)
- **Frontend:** React + Vite + TypeScript

## Rodar em modo desenvolvimento (sem Firebase)

### 1. Backend

```bash
cd backend
copy .env.example .env
go run ./cmd/server
```

API em `http://localhost:8080`. Com `DEV_MODE=true` (padrão sem credenciais), o login do professor é local: qualquer e-mail/senha cria a conta na memória.

### 2. Frontend

```bash
cd frontend
copy .env.example .env
npm install
npm run dev
```

Abra `http://localhost:5173`.

### Fluxo rápido

1. **Área do professor** → registrar/entrar
2. Criar quiz → adicionar questões (marque a correta, defina **peso** e **tempo**)
3. **Ao vivo** → anote o PIN
4. Em outro navegador/aba: **Entrar na partida** → PIN + nickname
5. No host: **Iniciar partida** → avance as questões → pódio final

### Rodar back + front juntos

Na raiz do projeto:

```bash
npm install
npm run dev
```

Isso sobe API (`:8080`) e frontend (`:5173`) juntos.

## Publicar online (Railway)

O projeto já tem `Dockerfile` + `railway.toml` (API + frontend no mesmo serviço).

### 1. Conta e CLI (opcional)

1. Crie conta em [https://railway.app](https://railway.app)
2. Instale a CLI (opcional): [https://docs.railway.com/guides/cli](https://docs.railway.com/guides/cli)

### 2. Subir o código para o GitHub

Railway conecta pelo GitHub. Na pasta do projeto:

```bash
git init
git add .
git commit -m "QuestArena ready for Railway"
```

Crie um repositório no GitHub e faça push (`git remote add origin ...` + `git push -u origin main`).

### 3. Deploy no Railway

1. No Railway: **New Project** → **Deploy from GitHub repo**
2. Escolha o repositório Kahoot/QuestArena
3. Railway detecta o `Dockerfile` e faz o build
4. Em **Settings → Networking → Generate Domain** (gera a URL pública, ex.: `https://questarena-xxxx.up.railway.app`)
5. Variáveis (opcional no início):

| Variável | Valor |
|----------|--------|
| `DEV_MODE` | `true` (login local; dados somem ao reiniciar) |
| `PORT` | definido automaticamente pelo Railway |

6. Abra a URL gerada — professor e aluno usam o mesmo link.

### Observações

- Em `DEV_MODE=true` os quizzes não persistem entre deploys/reinícios.
- Para persistir de verdade, configure Firebase depois (`DEV_MODE=false` + credenciais).
- WebSocket usa o mesmo domínio (`wss://sua-url/ws`).

## Modo Firebase (produção)

1. Crie um projeto no Firebase e ative **Authentication** (e-mail/senha) e **Firestore**.
2. Gere uma service account e salve o JSON (ex.: `backend/firebase-service-account.json`).
3. Publique as regras em [`firebase/firestore.rules`](firebase/firestore.rules).
4. Configure o backend:

```env
DEV_MODE=false
FIREBASE_CREDENTIALS=./firebase-service-account.json
FIREBASE_PROJECT_ID=seu-project-id
PORT=8080
```

5. Configure o frontend (`.env`):

```env
VITE_FIREBASE_API_KEY=...
VITE_FIREBASE_AUTH_DOMAIN=...
VITE_FIREBASE_PROJECT_ID=...
VITE_FIREBASE_APP_ID=...
```

O frontend usa o ID token do Firebase; o Go valida o token e persiste quizzes no Firestore.

## API principal

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/health` | Health + auth mode |
| POST | `/api/auth/dev-login` | Login local (DEV) |
| POST | `/api/auth/session` | Valida token |
| CRUD | `/api/quizzes` | Quizzes do professor |
| CRUD | `/api/quizzes/:id/questions` | Questões |
| POST | `/api/sessions` | Cria sala ao vivo `{ pin }` |
| WS | `/ws` | Partida ao vivo |

### Eventos WebSocket

- Host: `host_join`, `start`, `next`, `end`
- Aluno: `join`, `answer`
- Servidor: `hosted`, `joined`, `player_joined`, `question`, `answer_ack`, `question_result`, `finished`, `error`

### Pontuação

```
base = weight * 1000
score = round(base * (0.2 + 0.8 * (1 - elapsed/timeLimit)))
```

Só respostas corretas dentro do prazo recebem XP. Cálculo apenas no servidor; a alternativa correta não é enviada aos alunos durante a pergunta.

## Estrutura

```
backend/          # Go API + game hub
frontend/         # React app (professor + aluno)
firebase/         # Firestore security rules
```

## Notas

- Estado da partida ao vivo fica **em memória** no processo Go (reiniciar o servidor encerra salas abertas).
- Alunos não precisam de conta: só PIN + nickname.
- Em DEV, dados somem ao reiniciar o backend.
