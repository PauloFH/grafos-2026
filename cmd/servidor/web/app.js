// app.js — orquestracao do front EcoUrbano.
//
// Estado + chamadas a API (/api/datasets, /api/graph, /api/rota) + laco de
// animacao do caminhao + ligacao dos controles, paineis e timeline.

const VBOX_W = 1000;
const VBOX_H = 680;
const PX_POR_SEG = 150; // velocidade base do caminhao (unidades da viewBox / s)

const state = {
  datasets: [],
  datasetId: null,
  graph: null,
  positions: null,
  directed: false,
  start: null,
  classification: null,
  route: null, // { sequencia, circuito, arestas }
  segmentos: [], // { de, para, a, b, c, len, chave }
  segIndex: 0,
  segDist: 0, // px ja percorridos no segmento atual
  playing: false,
  speed: 1.5,
  totalLen: 0,
  ultimoAngulo: 0,
};

const $ = (id) => document.getElementById(id);

// ---------- API ----------

async function api(url) {
  const r = await fetch(url);
  if (!r.ok) {
    let msg = `HTTP ${r.status}`;
    try {
      const j = await r.json();
      if (j && j.erro) msg = j.erro;
    } catch (_) {}
    throw new Error(msg);
  }
  return r.json();
}

async function carregarDatasets() {
  state.datasets = await api("/api/datasets");
  const cont = $("modos");
  cont.replaceChildren();
  for (const d of state.datasets) {
    const b = document.createElement("button");
    b.className = "modo-btn";
    b.dataset.id = d.id;
    b.innerHTML = `<strong>${d.direcionado ? "Mão única" : "Mão dupla"}</strong>
      <span>${d.direcionado ? "Digrafo" : "Grafo"} · ${d.arestas} ruas</span>`;
    b.addEventListener("click", () => selecionarDataset(d.id));
    cont.appendChild(b);
  }
  if (state.datasets.length) await selecionarDataset(state.datasets[0].id);
}

async function selecionarDataset(id) {
  state.datasetId = id;
  for (const b of document.querySelectorAll(".modo-btn")) {
    b.classList.toggle("ativo", b.dataset.id === id);
  }
  state.graph = await api(`/api/graph?dataset=${encodeURIComponent(id)}`);
  state.directed = state.graph.direcionado;
  state.positions = computeLayout(state.graph.vertices, state.graph.arestas, VBOX_W, VBOX_H);
  Mapa.render(state.graph, state.positions);
  preencherInicioSelect();
  // Calcula a rota a partir do inicio sugerido pelo backend.
  await computarRota("");
}

async function computarRota(inicio) {
  let resp;
  try {
    const q = `/api/rota?dataset=${encodeURIComponent(state.datasetId)}${
      inicio ? `&inicio=${encodeURIComponent(inicio)}` : ""
    }`;
    resp = await api(q);
  } catch (e) {
    toast("Erro ao calcular rota: " + e.message);
    return;
  }

  state.classification = resp.classificacao;
  atualizarCardEuler(resp.classificacao);

  if (!resp.ok || !resp.trilha) {
    state.route = null;
    state.segmentos = [];
    Mapa.limparColetadas();
    Mapa.esconderCaminhao();
    $("btn-play").disabled = true;
    $("scrubber").disabled = true;
    toast(resp.erro || "Sem trilha euleriana para este ponto.");
    atualizarProgresso();
    construirTimeline();
    return;
  }

  $("btn-play").disabled = false;
  $("scrubber").disabled = false;
  state.route = resp.trilha;
  state.start = resp.trilha.sequencia[0];
  $("inicio").value = state.start;
  Mapa.definirInicio(state.start);
  construirSegmentos();
  construirTimeline();
  resetar();
}

// ---------- Segmentos e geometria ----------

function construirSegmentos() {
  const seq = state.route.sequencia;
  const segs = [];
  let total = 0;
  for (let i = 0; i < seq.length - 1; i++) {
    const de = seq[i];
    const para = seq[i + 1];
    const a = state.positions[de];
    const b = state.positions[para];
    const c = controlePonto(a, b, state.directed);
    const len = quadComprimento(a, c, b);
    total += len;
    segs.push({ de, para, a, b, c, len, chave: chaveAresta(de, para, state.directed) });
  }
  state.segmentos = segs;
  state.totalLen = total;
  const sc = $("scrubber");
  sc.max = segs.length;
}

// Recalcula geometria dos segmentos apos arrastar um cruzamento (sequencia igual).
function recomputarGeometria() {
  if (!state.segmentos.length) return;
  for (const s of state.segmentos) {
    s.a = state.positions[s.de];
    s.b = state.positions[s.para];
    s.c = controlePonto(s.a, s.b, state.directed);
    s.len = quadComprimento(s.a, s.c, s.b);
  }
  desenharCaminhao();
}

// ---------- Animacao ----------

let ultimoTs = 0;
function laco(ts) {
  const dt = ultimoTs ? (ts - ultimoTs) / 1000 : 0;
  ultimoTs = ts;
  if (state.playing && state.segmentos.length) {
    avancar(PX_POR_SEG * state.speed * dt);
  }
  desenharCaminhao();
  requestAnimationFrame(laco);
}

function avancar(px) {
  let rest = px;
  while (rest > 0 && state.segIndex < state.segmentos.length) {
    const seg = state.segmentos[state.segIndex];
    const falta = seg.len - state.segDist;
    if (rest < falta) {
      state.segDist += rest;
      rest = 0;
    } else {
      rest -= falta;
      completarSegmento(state.segIndex);
      state.segIndex++;
      state.segDist = 0;
    }
  }
  if (state.segIndex >= state.segmentos.length) finalizar();
  destacarAtual();
}

function completarSegmento(i) {
  Mapa.marcarColetada(state.segmentos[i].chave, i + 1);
  atualizarProgresso();
  marcarTimeline(i + 1);
}

function desenharCaminhao() {
  if (!state.segmentos.length) {
    Mapa.esconderCaminhao();
    return;
  }
  if (state.segIndex >= state.segmentos.length) {
    const ult = state.segmentos[state.segmentos.length - 1];
    Mapa.posicionarCaminhao(ult.b.x, ult.b.y, state.ultimoAngulo);
    return;
  }
  const seg = state.segmentos[state.segIndex];
  const t = seg.len > 0 ? Math.min(state.segDist / seg.len, 1) : 0;
  const p = quadAt(seg.a, seg.c, seg.b, t);
  const d = quadDeriv(seg.a, seg.c, seg.b, t);
  state.ultimoAngulo = (Math.atan2(d.y, d.x) * 180) / Math.PI;
  Mapa.posicionarCaminhao(p.x, p.y, state.ultimoAngulo);
}

function destacarAtual() {
  if (state.segIndex < state.segmentos.length) {
    Mapa.destacarAtual(state.segmentos[state.segIndex].chave);
  } else {
    Mapa.destacarAtual(null);
  }
}

function finalizar() {
  state.segIndex = state.segmentos.length;
  if (state.playing) {
    state.playing = false;
    atualizarBotaoPlay();
    if (state.route && state.route.circuito) {
      toast("Missão cumprida! Circuito fechado de volta ao depósito 🎉♻️");
      confete();
    } else {
      toast("Rota concluída — todas as ruas coletadas ✅");
    }
  }
  atualizarProgresso();
}

// ---------- Controles ----------

function play() {
  if (!state.segmentos.length) return;
  if (state.segIndex >= state.segmentos.length) resetar();
  state.playing = true;
  atualizarBotaoPlay();
}
function pausar() {
  state.playing = false;
  atualizarBotaoPlay();
}
function alternarPlay() {
  state.playing ? pausar() : play();
}

function resetar() {
  state.playing = false;
  state.segIndex = 0;
  state.segDist = 0;
  Mapa.limparColetadas();
  atualizarBotaoPlay();
  desenharCaminhao();
  destacarAtual();
  atualizarProgresso();
  marcarTimeline(0);
}

function passo(delta) {
  pausar();
  irPara(state.segIndex + delta);
}

// irPara posiciona o caminhao no INICIO do segmento k (k ruas ja coletadas).
function irPara(k) {
  k = Math.max(0, Math.min(k, state.segmentos.length));
  state.segIndex = k;
  state.segDist = 0;
  Mapa.limparColetadas();
  for (let i = 0; i < k; i++) Mapa.marcarColetada(state.segmentos[i].chave, i + 1);
  desenharCaminhao();
  destacarAtual();
  atualizarProgresso();
  marcarTimeline(k);
}

// ---------- Paineis ----------

function atualizarCardEuler(c) {
  const badge = $("euler-badge");
  const mapaClasse = {
    CircuitoEuleriano: ["Euleriano", "ok", "Tem circuito: dá pra coletar tudo e voltar ao depósito."],
    CaminhoEuleriano: ["Semi-euleriano", "warn", "Tem caminho, mas começa e termina em cruzamentos diferentes."],
    NaoEuleriano: ["Não-euleriano", "bad", "Não existe rota que percorra cada rua exatamente uma vez."],
  };
  const [rotulo, cls, frase] = mapaClasse[c.classe] || ["—", "", ""];
  badge.textContent = rotulo;
  badge.className = "badge " + cls;
  $("euler-texto").textContent = c.texto;
  $("euler-conexo").textContent = c.conexo ? "Sim" : "Não";
  $("euler-inicio").textContent = c.verticeInicial || "—";
  $("euler-explica").textContent = frase;

  const det = $("euler-detalhe");
  if (state.directed) {
    const ent = Object.entries(c.desbalanceados || {});
    det.innerHTML =
      "<span class='rotulo'>Vértices desbalanceados (saída−entrada):</span> " +
      (ent.length ? ent.map(([v, n]) => `${v}: ${n > 0 ? "+" : ""}${n}`).join(", ") : "nenhum (todos balanceados)");
  } else {
    const imp = c.grausImpares || [];
    det.innerHTML =
      "<span class='rotulo'>Vértices de grau ímpar:</span> " +
      (imp.length ? imp.join(", ") : "nenhum (todos pares)");
  }
}

function atualizarProgresso() {
  const total = state.segmentos.length;
  const feitas = Math.min(state.segIndex, total);
  const pct = total ? (feitas / total) * 100 : 0;
  $("prog-ruas").textContent = `${feitas} / ${total}`;
  $("prog-barra").style.width = pct.toFixed(1) + "%";
  $("prog-percent").textContent = pct.toFixed(0) + "%";

  if (total && state.segIndex < total) {
    const s = state.segmentos[state.segIndex];
    $("prog-atual").textContent = `${s.de} ${state.directed ? "→" : "↔"} ${s.para}`;
  } else if (total) {
    $("prog-atual").textContent = "Concluído ✓";
  } else {
    $("prog-atual").textContent = "—";
  }
  $("prog-tipo").textContent = state.route
    ? state.route.circuito
      ? "Circuito (volta ao início)"
      : "Caminho (início ≠ fim)"
    : "—";
}

// ---------- Timeline ----------

function construirTimeline() {
  const tl = $("timeline");
  tl.replaceChildren();
  if (!state.route) {
    tl.innerHTML = "<span class='tl-vazia'>Sem rota.</span>";
    return;
  }
  const seq = state.route.sequencia;
  seq.forEach((v, i) => {
    if (i > 0) {
      const seta = document.createElement("span");
      seta.className = "tl-seta";
      seta.textContent = state.directed ? "→" : "·";
      tl.appendChild(seta);
    }
    const sp = document.createElement("button");
    sp.className = "tl-no";
    sp.textContent = v;
    sp.title = `Ir até a parada ${i}`;
    sp.addEventListener("click", () => {
      pausar();
      irPara(i);
    });
    tl.appendChild(sp);
  });
  marcarTimeline(0);
}

function marcarTimeline(k) {
  const nos = document.querySelectorAll("#timeline .tl-no");
  nos.forEach((n, i) => {
    n.classList.toggle("feito", i < k);
    n.classList.toggle("atual", i === k);
  });
  const atual = nos[k];
  if (atual) atual.scrollIntoView({ inline: "center", block: "nearest", behavior: "smooth" });
}

// ---------- UI helpers ----------

function preencherInicioSelect() {
  const sel = $("inicio");
  sel.replaceChildren();
  for (const v of state.graph.vertices) {
    const o = document.createElement("option");
    o.value = v;
    o.textContent = "Cruzamento " + v;
    sel.appendChild(o);
  }
}

function atualizarBotaoPlay() {
  const b = $("btn-play");
  b.textContent = state.playing ? "⏸  Pausar" : "▶  Coletar";
  b.classList.toggle("tocando", state.playing);
}

let toastTimer = null;
function toast(msg) {
  const t = $("toast");
  t.textContent = msg;
  t.classList.add("visivel");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => t.classList.remove("visivel"), 4200);
}

// Confete simples em canvas overlay.
function confete() {
  const cv = $("confetti");
  const ctx = cv.getContext("2d");
  cv.width = cv.clientWidth;
  cv.height = cv.clientHeight;
  const cores = ["#34d399", "#10b981", "#a7f3d0", "#fbbf24", "#60a5fa", "#f472b6"];
  const ps = [];
  for (let i = 0; i < 140; i++) {
    ps.push({
      x: cv.width / 2 + (Math.random() - 0.5) * 120,
      y: cv.height / 3,
      vx: (Math.random() - 0.5) * 9,
      vy: Math.random() * -10 - 4,
      g: 0.3 + Math.random() * 0.2,
      s: 4 + Math.random() * 6,
      cor: cores[(Math.random() * cores.length) | 0],
      rot: Math.random() * Math.PI,
      vr: (Math.random() - 0.5) * 0.3,
    });
  }
  let frames = 0;
  function anim() {
    ctx.clearRect(0, 0, cv.width, cv.height);
    for (const p of ps) {
      p.vy += p.g;
      p.x += p.vx;
      p.y += p.vy;
      p.rot += p.vr;
      ctx.save();
      ctx.translate(p.x, p.y);
      ctx.rotate(p.rot);
      ctx.fillStyle = p.cor;
      ctx.fillRect(-p.s / 2, -p.s / 2, p.s, p.s * 0.6);
      ctx.restore();
    }
    frames++;
    if (frames < 150) requestAnimationFrame(anim);
    else ctx.clearRect(0, 0, cv.width, cv.height);
  }
  anim();
}

// ---------- Bootstrap ----------

window.addEventListener("DOMContentLoaded", async () => {
  Mapa.init($("mapa"), {
    onClicar: (id) => escolherInicio(id),
    onArrastar: () => recomputarGeometria(),
  });

  $("btn-play").addEventListener("click", alternarPlay);
  $("btn-reset").addEventListener("click", resetar);
  $("btn-passo-frente").addEventListener("click", () => passo(1));
  $("btn-passo-tras").addEventListener("click", () => passo(-1));
  $("inicio").addEventListener("change", (e) => escolherInicio(e.target.value));

  const vel = $("velocidade");
  vel.addEventListener("input", () => {
    state.speed = Number(vel.value);
    $("velocidade-val").textContent = state.speed.toFixed(2) + "×";
  });
  state.speed = Number(vel.value);
  $("velocidade-val").textContent = state.speed.toFixed(2) + "×";

  $("scrubber").addEventListener("input", (e) => {
    pausar();
    irPara(Number(e.target.value));
  });

  $("toggle-numeros").addEventListener("change", (e) => {
    document.body.classList.toggle("mostrar-numeros", e.target.checked);
  });
  document.body.classList.toggle("mostrar-numeros", $("toggle-numeros").checked);

  try {
    await carregarDatasets();
  } catch (e) {
    toast("Falha ao carregar dados: " + e.message);
  }

  requestAnimationFrame(laco);
});

async function escolherInicio(v) {
  if (!state.graph || !state.graph.vertices.includes(v)) return;
  state.start = v;
  $("inicio").value = v;
  Mapa.definirInicio(v);
  Mapa.pulsarNo(v);
  await computarRota(v);
}

// Mantem o scrubber sincronizado a cada frame de progresso.
const _atualizarProgressoOrig = atualizarProgresso;
atualizarProgresso = function () {
  _atualizarProgressoOrig();
  const sc = $("scrubber");
  if (sc && document.activeElement !== sc) sc.value = Math.min(state.segIndex, state.segmentos.length);
};
