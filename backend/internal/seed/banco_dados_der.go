package seed

// Reconstruções didáticas (notação Chen). Não são o scan do caderno do INEP.
const derEntidades2011 = `<svg viewBox="0 0 760 340" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="DER Ent1 Ent2 Ent3 Ent4 Ent5">
  <rect x="36" y="18" width="176" height="92" rx="6" fill="#122033" stroke="#2dd4bf" stroke-width="2"/>
  <text x="124" y="42" text-anchor="middle" fill="#fbbf24" font-size="12" font-weight="700">PK: cent11, cent12</text>
  <text x="124" y="68" text-anchor="middle" fill="#f4f7fb" font-size="18" font-weight="700">Ent1</text>
  <text x="124" y="90" text-anchor="middle" fill="#9db0c7" font-size="11">forte</text>

  <polygon points="258,64 286,86 258,108 230,86" fill="#1a2740" stroke="#ff8a3d" stroke-width="2"/>
  <text x="258" y="90" text-anchor="middle" fill="#ff8a3d" font-size="10">N:N</text>
  <line x1="212" y1="64" x2="230" y2="86" stroke="#9db0c7" stroke-width="1.5"/>
  <line x1="286" y1="86" x2="304" y2="64" stroke="#9db0c7" stroke-width="1.5"/>

  <rect x="304" y="18" width="160" height="92" rx="6" fill="#122033" stroke="#2dd4bf" stroke-width="2"/>
  <text x="384" y="42" text-anchor="middle" fill="#fbbf24" font-size="12" font-weight="700">PK: cent21</text>
  <text x="384" y="68" text-anchor="middle" fill="#f4f7fb" font-size="18" font-weight="700">Ent2</text>
  <text x="384" y="90" text-anchor="middle" fill="#9db0c7" font-size="11">forte</text>

  <line x1="124" y1="110" x2="124" y2="208" stroke="#9db0c7" stroke-width="1.5"/>
  <text x="148" y="150" fill="#f4f7fb" font-size="12">0..1</text>
  <text x="148" y="196" fill="#f4f7fb" font-size="12">0..N</text>

  <line x1="384" y1="110" x2="384" y2="208" stroke="#9db0c7" stroke-width="1.5"/>
  <text x="408" y="150" fill="#f4f7fb" font-size="12">1</text>
  <text x="408" y="196" fill="#f4f7fb" font-size="12">0..N</text>

  <rect x="36" y="208" width="176" height="92" rx="22" fill="#122033" stroke="#86efac" stroke-width="2"/>
  <text x="124" y="232" text-anchor="middle" fill="#fbbf24" font-size="11" font-weight="700">cent51 + PK de Ent1</text>
  <text x="124" y="258" text-anchor="middle" fill="#86efac" font-size="17" font-weight="700">Ent5</text>
  <text x="124" y="280" text-anchor="middle" fill="#9db0c7" font-size="11">fraca (cantos redondos)</text>

  <rect x="304" y="208" width="160" height="92" rx="22" fill="#122033" stroke="#86efac" stroke-width="2"/>
  <text x="384" y="232" text-anchor="middle" fill="#fbbf24" font-size="11" font-weight="700">cent31 + PK de Ent2</text>
  <text x="384" y="258" text-anchor="middle" fill="#86efac" font-size="17" font-weight="700">Ent3</text>
  <text x="384" y="280" text-anchor="middle" fill="#9db0c7" font-size="11">fraca (cantos redondos)</text>

  <polygon points="530,254 558,276 530,298 502,276" fill="#1a2740" stroke="#ff8a3d" stroke-width="2"/>
  <text x="530" y="280" text-anchor="middle" fill="#ff8a3d" font-size="10">N:N</text>
  <line x1="464" y1="254" x2="502" y2="276" stroke="#9db0c7" stroke-width="1.5"/>
  <line x1="558" y1="276" x2="588" y2="254" stroke="#9db0c7" stroke-width="1.5"/>

  <rect x="588" y="208" width="150" height="92" rx="6" fill="#122033" stroke="#2dd4bf" stroke-width="2"/>
  <text x="663" y="232" text-anchor="middle" fill="#fbbf24" font-size="12" font-weight="700">PK: cent41</text>
  <text x="663" y="258" text-anchor="middle" fill="#f4f7fb" font-size="18" font-weight="700">Ent4</text>
  <text x="663" y="280" text-anchor="middle" fill="#9db0c7" font-size="11">forte</text>

  <text x="380" y="328" text-anchor="middle" fill="#9db0c7" font-size="11">teal = forte · verde arredondado = fraca · losango = relacionamento</text>
</svg>`

const derPessoaCavalo2017 = `<svg viewBox="0 0 680 220" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="DER Pessoa possui Cavalo">
  <rect x="28" y="40" width="200" height="110" rx="6" fill="#122033" stroke="#2dd4bf" stroke-width="2"/>
  <text x="128" y="68" text-anchor="middle" fill="#fbbf24" font-size="13" font-weight="700">id (identificador)</text>
  <text x="128" y="96" text-anchor="middle" fill="#f4f7fb" font-size="20" font-weight="700">Pessoa</text>
  <text x="128" y="122" text-anchor="middle" fill="#9db0c7" font-size="13">nome, rg (descritivo)</text>

  <line x1="228" y1="95" x2="292" y2="95" stroke="#9db0c7" stroke-width="1.6"/>
  <polygon points="340,72 372,95 340,118 308,95" fill="#1a2740" stroke="#ff8a3d" stroke-width="2"/>
  <text x="340" y="99" text-anchor="middle" fill="#ff8a3d" font-size="12">possui</text>
  <line x1="372" y1="95" x2="452" y2="95" stroke="#9db0c7" stroke-width="1.6"/>

  <text x="250" y="36" fill="#f4f7fb" font-size="14" font-weight="700">0..1</text>
  <text x="400" y="36" fill="#f4f7fb" font-size="14" font-weight="700">1</text>

  <rect x="452" y="40" width="200" height="110" rx="6" fill="#122033" stroke="#2dd4bf" stroke-width="2"/>
  <text x="552" y="88" text-anchor="middle" fill="#f4f7fb" font-size="20" font-weight="700">Cavalo</text>
  <text x="552" y="116" text-anchor="middle" fill="#9db0c7" font-size="13">participação total</text>

  <text x="340" y="188" text-anchor="middle" fill="#9db0c7" font-size="12">pessoa pode não ter cavalo · cavalo deve ter dono · no máximo um cavalo por pessoa</text>
  <text x="340" y="208" text-anchor="middle" fill="#9db0c7" font-size="11">rg não é o identificador no desenho</text>
</svg>`
