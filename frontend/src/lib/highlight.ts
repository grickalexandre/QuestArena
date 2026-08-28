export type TokenKind =
  | 'txt'
  | 'kw'
  | 'type'
  | 'str'
  | 'num'
  | 'com'
  | 'fn'
  | 'var'
  | 'ann'
  | 'op'
  | 'punc'

export type Token = { t: string; k: TokenKind }

export type CodeLanguage =
  | 'plain'
  | 'java'
  | 'python'
  | 'javascript'
  | 'typescript'
  | 'csharp'
  | 'c'
  | 'cpp'
  | 'go'
  | 'php'
  | 'sql'
  | 'kotlin'
  | 'html'
  | 'css'
  | 'json'
  | 'shell'

export const CODE_LANGUAGES: { id: CodeLanguage; label: string }[] = [
  { id: 'plain', label: 'Texto simples' },
  { id: 'java', label: 'Java' },
  { id: 'python', label: 'Python' },
  { id: 'javascript', label: 'JavaScript' },
  { id: 'typescript', label: 'TypeScript' },
  { id: 'csharp', label: 'C#' },
  { id: 'c', label: 'C' },
  { id: 'cpp', label: 'C++' },
  { id: 'go', label: 'Go' },
  { id: 'php', label: 'PHP' },
  { id: 'sql', label: 'SQL' },
  { id: 'kotlin', label: 'Kotlin' },
  { id: 'html', label: 'HTML' },
  { id: 'css', label: 'CSS' },
  { id: 'json', label: 'JSON' },
  { id: 'shell', label: 'Shell / Terminal' },
]

export function languageLabel(id: string): string {
  return CODE_LANGUAGES.find((l) => l.id === id)?.label ?? 'Texto simples'
}

type LangDef = {
  keywords: string[]
  builtins?: string[]
  line?: string[]
  block?: [string, string]
  quotes?: string[]
  /** Quote characters that may span multiple lines (template literals, raw strings). */
  multiline?: string[]
  triple?: boolean
  caseInsensitive?: boolean
  /** Treat Identifiers starting with an uppercase letter as types. */
  pascalTypes?: boolean
  markup?: boolean
}

const words = (s: string) => s.split(/\s+/).filter(Boolean)

const JS_KEYWORDS =
  'as async await break case catch class const continue debugger default delete do else export extends finally for from function get if import in instanceof let new of return set static super switch this throw try typeof var void while yield'
const TS_KEYWORDS = `${JS_KEYWORDS} abstract declare enum implements interface is keyof namespace private protected public readonly satisfies type`

const LANGS: Record<string, LangDef> = {
  plain: { keywords: [] },

  javascript: {
    keywords: words(JS_KEYWORDS),
    builtins: words(
      'true false null undefined NaN Infinity console document window Math JSON Object Array String Number Boolean Promise Map Set Symbol Error',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'", '`'],
    multiline: ['`'],
    pascalTypes: true,
  },

  typescript: {
    keywords: words(TS_KEYWORDS),
    builtins: words(
      'true false null undefined any unknown never string number boolean object void console Promise Record Array Partial Readonly',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'", '`'],
    multiline: ['`'],
    pascalTypes: true,
  },

  java: {
    keywords: words(
      'abstract assert break case catch class const continue default do else enum extends final finally for goto if implements import instanceof interface native new package private protected public return static strictfp super switch synchronized this throw throws transient try volatile while var record sealed permits yield',
    ),
    builtins: words(
      'boolean byte char double float int long short void String Integer Double Boolean Character Long Object System List ArrayList Map HashMap Set HashSet Scanner Math true false null',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    pascalTypes: true,
  },

  kotlin: {
    keywords: words(
      'as break by catch class companion const continue crossinline data do else enum external final finally for fun get if import in infix init inline interface internal is lateinit noinline object open operator out override package private protected public reified return sealed set super suspend this throw try typealias val var vararg when while',
    ),
    builtins: words(
      'Int Long Double Float Boolean String Char Any Unit Nothing List MutableList Map Set Array println print true false null it',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    multiline: ['"'],
    pascalTypes: true,
  },

  csharp: {
    keywords: words(
      'abstract as base break case catch checked class const continue default delegate do else enum event explicit extern finally fixed for foreach goto if implicit in interface internal is lock namespace new operator out override params private protected public readonly ref return sealed sizeof stackalloc static struct switch this throw try typeof unchecked unsafe using virtual volatile while async await get set var record yield',
    ),
    builtins: words(
      'bool byte char decimal double float int long object sbyte short string uint ulong ushort void List Dictionary Console Math true false null value',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    pascalTypes: true,
  },

  c: {
    keywords: words(
      'auto break case const continue default do else enum extern for goto if inline register restrict return sizeof static struct switch typedef union volatile while',
    ),
    builtins: words(
      'char double float int long short signed unsigned void size_t FILE NULL printf scanf malloc free include define true false',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
  },

  cpp: {
    keywords: words(
      'auto break case catch class const constexpr continue default delete do else enum explicit export extern for friend goto if inline namespace new operator private protected public return sizeof static struct switch template this throw try typedef typename union using virtual volatile while',
    ),
    builtins: words(
      'bool char double float int long short signed unsigned void string vector map set cout cin endl std nullptr NULL true false include define',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    pascalTypes: true,
  },

  go: {
    keywords: words(
      'break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var',
    ),
    builtins: words(
      'bool byte complex64 complex128 error float32 float64 int int8 int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr append cap close copy delete len make new panic print println recover nil true false iota fmt',
    ),
    line: ['//'],
    block: ['/*', '*/'],
    quotes: ['"', "'", '`'],
    multiline: ['`'],
    pascalTypes: true,
  },

  python: {
    keywords: words(
      'and as assert async await break class continue def del elif else except finally for from global if import in is lambda nonlocal not or pass raise return try while with yield match case',
    ),
    builtins: words(
      'True False None self cls print len range int str float bool list dict set tuple input open enumerate zip map filter sum min max abs sorted type isinstance super Exception ValueError TypeError',
    ),
    line: ['#'],
    quotes: ['"', "'"],
    triple: true,
    pascalTypes: true,
  },

  php: {
    keywords: words(
      'abstract and array as break callable case catch class clone const continue declare default do echo else elseif empty enddeclare endfor endforeach endif endswitch endwhile enum extends final finally fn for foreach function global goto if implements include include_once instanceof insteadof interface isset list match namespace new or print private protected public readonly require require_once return static switch throw trait try unset use var while xor yield',
    ),
    builtins: words(
      'true false null int float string bool object void self parent this count strlen array_map var_dump printf sprintf',
    ),
    line: ['//', '#'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    pascalTypes: true,
  },

  sql: {
    keywords: words(
      'select from where insert into values update set delete create table alter drop truncate join inner left right full outer cross on group by order having limit offset distinct as and or not null is in between like exists union all case when then else end primary key foreign references constraint index view database schema default check unique cascade begin commit rollback transaction with recursive returning',
    ),
    builtins: words(
      'int integer bigint smallint decimal numeric float double real varchar char text date datetime timestamp time boolean serial uuid json count sum avg min max coalesce nullif cast now current_date current_timestamp true false',
    ),
    line: ['--'],
    block: ['/*', '*/'],
    quotes: ['"', "'"],
    caseInsensitive: true,
  },

  shell: {
    keywords: words(
      'if then elif else fi for in do done while until case esac function select break continue return exit local export readonly declare source alias',
    ),
    builtins: words(
      'echo cd ls cp mv rm mkdir rmdir cat grep sed awk curl wget git npm node python pip docker sudo chmod chown touch pwd env set unset printf read test',
    ),
    line: ['#'],
    quotes: ['"', "'"],
  },

  css: {
    keywords: words(
      'important media supports keyframes import charset font-face from to and not only screen print',
    ),
    builtins: words(
      'inherit initial unset auto none flex grid block inline inline-block absolute relative fixed sticky static hidden visible bold italic center left right solid dashed dotted transparent',
    ),
    block: ['/*', '*/'],
    quotes: ['"', "'"],
  },

  json: {
    keywords: [],
    builtins: words('true false null'),
    quotes: ['"'],
  },

  html: { keywords: [], markup: true },
}

const IDENT_START = /[A-Za-z_$]/
const IDENT_PART = /[A-Za-z0-9_$-]/

function isIdentStart(ch: string) {
  return IDENT_START.test(ch)
}

function tokenizeCode(src: string, def: LangDef): Token[] {
  const out: Token[] = []
  const quotes = def.quotes ?? []
  const multiline = def.multiline ?? []
  const keywords = new Set(def.caseInsensitive ? def.keywords.map((k) => k.toLowerCase()) : def.keywords)
  const builtins = new Set(
    (def.builtins ?? []).map((k) => (def.caseInsensitive ? k.toLowerCase() : k)),
  )
  const n = src.length
  let i = 0
  let buf = ''

  const flush = () => {
    if (buf) {
      out.push({ t: buf, k: 'txt' })
      buf = ''
    }
  }
  const push = (text: string, kind: TokenKind) => {
    if (!text) return
    flush()
    out.push({ t: text, k: kind })
  }

  while (i < n) {
    const ch = src[i]

    if (def.block && src.startsWith(def.block[0], i)) {
      const found = src.indexOf(def.block[1], i + def.block[0].length)
      const stop = found === -1 ? n : found + def.block[1].length
      push(src.slice(i, stop), 'com')
      i = stop
      continue
    }

    const lineComment = (def.line ?? []).find((c) => src.startsWith(c, i))
    if (lineComment) {
      const found = src.indexOf('\n', i)
      const stop = found === -1 ? n : found
      push(src.slice(i, stop), 'com')
      i = stop
      continue
    }

    if (quotes.includes(ch)) {
      const isTriple = def.triple && src.startsWith(ch.repeat(3), i)
      const stop = isTriple
        ? closingTriple(src, i, ch)
        : closingQuote(src, i, ch, multiline.includes(ch))
      push(src.slice(i, stop), 'str')
      i = stop
      continue
    }

    if (ch >= '0' && ch <= '9') {
      let j = i
      while (j < n && /[0-9a-fA-FxXoObB_.]/.test(src[j])) j++
      if (/[eE]/.test(src[j - 1] ?? '') && /[+-]/.test(src[j] ?? '')) {
        j++
        while (j < n && /[0-9]/.test(src[j])) j++
      }
      push(src.slice(i, j), 'num')
      i = j
      continue
    }

    if (ch === '@' && isIdentStart(src[i + 1] ?? '')) {
      let j = i + 1
      while (j < n && IDENT_PART.test(src[j])) j++
      push(src.slice(i, j), 'ann')
      i = j
      continue
    }

    if (isIdentStart(ch)) {
      let j = i
      while (j < n && IDENT_PART.test(src[j])) j++
      const word = src.slice(i, j)
      const lookup = def.caseInsensitive ? word.toLowerCase() : word
      let k: TokenKind = 'txt'
      if (keywords.has(lookup)) k = 'kw'
      else if (builtins.has(lookup)) k = 'type'
      else if (word.startsWith('$')) k = 'var'
      else {
        let p = j
        while (p < n && (src[p] === ' ' || src[p] === '\t')) p++
        if (src[p] === '(') k = 'fn'
        else if (def.pascalTypes && /^[A-Z]/.test(word)) k = 'type'
      }
      push(word, k)
      i = j
      continue
    }

    if ('{}()[];,'.includes(ch)) {
      push(ch, 'punc')
      i++
      continue
    }

    if ('+-*/%=<>!&|^~?:.'.includes(ch)) {
      push(ch, 'op')
      i++
      continue
    }

    buf += ch
    i++
  }

  flush()
  return out
}

function closingQuote(src: string, start: number, quote: string, allowNewline: boolean): number {
  let i = start + 1
  while (i < src.length) {
    const ch = src[i]
    if (ch === '\\') {
      i += 2
      continue
    }
    if (ch === quote) return i + 1
    if (ch === '\n' && !allowNewline) return i
    i++
  }
  return src.length
}

function closingTriple(src: string, start: number, quote: string): number {
  const mark = quote.repeat(3)
  const found = src.indexOf(mark, start + 3)
  return found === -1 ? src.length : found + 3
}

function tokenizeMarkup(src: string): Token[] {
  const out: Token[] = []
  const n = src.length
  let i = 0
  let text = ''
  const flushText = () => {
    if (text) {
      out.push({ t: text, k: 'txt' })
      text = ''
    }
  }

  while (i < n) {
    if (src.startsWith('<!--', i)) {
      const found = src.indexOf('-->', i)
      const stop = found === -1 ? n : found + 3
      flushText()
      out.push({ t: src.slice(i, stop), k: 'com' })
      i = stop
      continue
    }

    if (src[i] === '<') {
      const found = src.indexOf('>', i)
      const stop = found === -1 ? n : found + 1
      flushText()
      out.push(...tokenizeTag(src.slice(i, stop)))
      i = stop
      continue
    }

    text += src[i]
    i++
  }

  flushText()
  return out
}

function tokenizeTag(tag: string): Token[] {
  const out: Token[] = []
  const parts = tag.split(/("[^"]*"|'[^']*')/)
  parts.forEach((part, idx) => {
    if (!part) return
    if (idx % 2 === 1) {
      out.push({ t: part, k: 'str' })
      return
    }
    const pieces = part.split(/([A-Za-z_:][\w:.-]*)/)
    let tagNameSeen = false
    pieces.forEach((piece, pIdx) => {
      if (!piece) return
      if (pIdx % 2 === 1) {
        out.push({ t: piece, k: tagNameSeen ? 'ann' : 'kw' })
        tagNameSeen = true
      } else {
        out.push({ t: piece, k: 'punc' })
      }
    })
  })
  return out
}

export function tokenize(src: string, language: string): Token[] {
  const def = LANGS[language] ?? LANGS.plain
  if (def.markup) return tokenizeMarkup(src)
  if (!def.keywords.length && !def.quotes?.length && !def.line && !def.block) {
    return src ? [{ t: src, k: 'txt' }] : []
  }
  return tokenizeCode(src, def)
}

/** Tokenizes and splits the result per source line, so gutters stay aligned. */
export function tokenizeLines(src: string, language: string): Token[][] {
  const lines: Token[][] = [[]]
  for (const token of tokenize(src, language)) {
    const parts = token.t.split('\n')
    parts.forEach((part, idx) => {
      if (idx > 0) lines.push([])
      if (part) lines[lines.length - 1].push({ t: part, k: token.k })
    })
  }
  return lines
}
