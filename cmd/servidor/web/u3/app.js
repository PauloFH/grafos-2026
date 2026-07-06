// app.js — orquestracao do front U3


const METODOS = [
  { valor: "vmp", rotulo: "VMP + 2-opt", nome: "Vizinho mais Próximo + 2-opt", usaSemente: false },
  { valor: "imb", rotulo: "IMB + Or-opt", nome: "Inserção mais Barata + Or-opt", usaSemente: false },
  { valor: "ag", rotulo: "Alg. Genético", nome: "Algoritmo Genético", usaSemente: true },
  { valor: "memetico", rotulo: "Memético", nome: "Algoritmo Memético", usaSemente: true },
];

const state = {
  instancias: [], // [{id, medida, n}]
  instancia: null, // InstanciaResponse corrente
  nomePorId: new Map(), // id original -> nome da cidade
  ultimaResp: null, // ultimo SolveResponse ok
  resolvendo: false,
};

// mapa e o renderer ativo
let mapa = MapaPCV;


const anim = {
  ordem: [], // ids originais da rota corrente (sem repetir o deposito)
  segs: [],
  inicioPerna: [], // indice do primeiro micro-segmento de cada perna
  totalPernas: 0,
  segIndex: 0,
  segDist: 0, // unidades do renderer ja percorridas no segmento atual
  playing: false,
  speed: 1.5,
  ultimoAngulo: 0,
};

const $ = (id) => document.getElementById(id);

// ---------- API ----------

// api faz um GET e devolve o JSON; em status nao-2xx tenta extrair j.erro
// do corpo (mesma logica do front da U2).
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

// ---------- carga e troca de instancia ----------

async function carregarInstancias() {
  state.instancias = await api("/api/u3/instancias");
  const sel = $("instancia");
  sel.replaceChildren();
  for (const inst of state.instancias) {
    const opt = document.createElement("option");
    opt.value = inst.id;
    opt.textContent = `${inst.id} · ${inst.medida} · N=${inst.n}`;
    sel.appendChild(opt);
  }
  if (state.instancias.length) await selecionarInstancia(state.instancias[0].id);
}

async function selecionarInstancia(id) {
  let resp;
  try {
    resp = await api(`/api/u3/instancia?id=${encodeURIComponent(id)}`);
  } catch (e) {
    toast("Erro ao carregar instância: " + e.message);
    return;
  }
  state.instancia = resp;
  state.nomePorId = new Map(resp.pontos.map((p) => [p.id, p.nome]));
  state.ultimaResp = null;
  limparAnimacao();
  mapa.renderPontos(resp.pontos, resp.deposito);
  limparResultado();
}

// ---------- resolve ----------

async function resolver(metodo) {
  if (!state.instancia || state.resolvendo) return;

  const cfg = METODOS.find((m) => m.valor === metodo);
  let url = `/api/u3/resolve?id=${encodeURIComponent(state.instancia.id)}&metodo=${encodeURIComponent(metodo)}`;
  if (cfg && cfg.usaSemente) {
    const semente = $("semente").value.trim();
    if (semente !== "") url += `&semente=${encodeURIComponent(semente)}`;
  }

  setResolvendo(true);
  let resp;
  try {
    resp = await api(url);
  } catch (e) {
    toast("Erro ao resolver: " + e.message);
    return;
  } finally {
    setResolvendo(false);
  }

  if (!resp.ok) {
    // stub/erro de negocio: mensagem sem quebrar a pagina, mantendo os
    // pontos ja desenhados
    toast(resp.erro || "Falha ao resolver.");
    return;
  }

  state.ultimaResp = resp;
  mapa.renderRota(resp.rota.ordem);
  mapa.renderRotaAntes(resp.antes ? resp.antes.ordem : []);
  preencherResultado(cfg, resp);

  montarAnimacao(resp.rota.ordem);
  play(); // comeca a percorrer; o usuario pausa/ajusta nos controles
}

// ---------- alternancia diagrama MDS <-> mapa real ----------

// usarMapaReal troca o renderer ativo e re-renderiza a instancia e a rota
// correntes no novo espaco (a animacao recomeca do inicio).
function usarMapaReal(ativo) {
  if (ativo && !MapaReal.disponivel()) {
    toast("Mapa real indisponível — sem acesso ao Leaflet/tiles (internet).");
    $("toggle-mapa-real").checked = false;
    return;
  }
  pausar();
  mapa.limparRotas();
  mapa = ativo ? MapaReal : MapaPCV;
  document.body.classList.toggle("modo-mapa-real", ativo);
  if (ativo) {
    MapaReal.init("mapa-real"); // no-op se ja criado
    MapaReal.ativar(); // recalcula tamanho apos sair de display:none
  }
  if (state.instancia) {
    mapa.renderPontos(state.instancia.pontos, state.instancia.deposito);
    if (state.ultimaResp) {
      mapa.renderRota(state.ultimaResp.rota.ordem);
      if ($("toggle-antes").checked) {
        mapa.renderRotaAntes(state.ultimaResp.antes ? state.ultimaResp.antes.ordem : []);
      }
      montarAnimacao(state.ultimaResp.rota.ordem);
    }
  }
}

// setResolvendo desabilita os botoes de metodo durante o request (AG e
// memetico podem levar alguns segundos em N=48).
function setResolvendo(ativo) {
  state.resolvendo = ativo;
  for (const b of document.querySelectorAll(".metodo-btn")) b.disabled = ativo;
}

// montarAnimacao converte a rota nas pernas do ciclo fechado (via renderer
// ativo), explode cada perna em micro-segmentos, habilita os controles e
// reconstroi a folha de rota. Nao inicia o percurso.
function montarAnimacao(ordem) {
  anim.ordem = ordem || [];
  anim.segs = [];
  anim.inicioPerna = [];
  const pernas = mapa.caminhosDaOrdem(anim.ordem);
  anim.totalPernas = pernas.length;
  for (let k = 0; k < pernas.length; k++) {
    anim.inicioPerna.push(anim.segs.length);
    const pts = pernas[k];
    for (let i = 0; i + 1 < pts.length; i++) {
      const a = pts[i];
      const b = pts[i + 1];
      const len = Math.hypot(b.x - a.x, b.y - a.y);
      if (len < 1e-9) continue;
      anim.segs.push({
        a,
        b,
        len,
        perna: k,
        ang: (Math.atan2(b.y - a.y, b.x - a.x) * 180) / Math.PI,
      });
    }
  }
  $("scrubber").max = anim.totalPernas;
  habilitarControles(anim.segs.length > 0);
  construirTimeline();
  resetar();
}

// pernaAtual devolve o indice da perna em curso.
function pernaAtual() {
  if (!anim.segs.length) return 0;
  return anim.segIndex >= anim.segs.length ? anim.totalPernas : anim.segs[anim.segIndex].perna;
}

// limparAnimacao descarta o percurso corrente (troca de instancia).
function limparAnimacao() {
  anim.ordem = [];
  anim.segs = [];
  anim.inicioPerna = [];
  anim.totalPernas = 0;
  anim.segIndex = 0;
  anim.segDist = 0;
  anim.playing = false;
  atualizarBotaoPlay();
  habilitarControles(false);
  $("scrubber").max = 0;
  $("scrubber").value = 0;
  $("timeline").replaceChildren();
  mapa.esconderCaminhao();
}

// laco roda continuamente (requestAnimationFrame); so avanca quando playing.
let ultimoTs = 0;
function laco(ts) {
  const dt = ultimoTs ? Math.min((ts - ultimoTs) / 1000, 0.05) : 0;
  ultimoTs = ts;
  if (anim.playing && anim.segs.length) {
    avancar(mapa.velocidadeBase * anim.speed * dt);
    desenharCaminhao();
  }
  requestAnimationFrame(laco);
}

function avancar(px) {
  const pernaAntes = pernaAtual();
  let rest = px;
  while (rest > 0 && anim.segIndex < anim.segs.length) {
    const seg = anim.segs[anim.segIndex];
    const falta = seg.len - anim.segDist;
    if (rest < falta) {
      anim.segDist += rest;
      rest = 0;
    } else {
      rest -= falta;
      anim.segIndex++;
      anim.segDist = 0;
    }
  }
  if (pernaAtual() !== pernaAntes) {
    marcarTimeline(pernaAtual());
    atualizarProgresso();
  }
  if (anim.segIndex >= anim.segs.length) finalizar();
}

// desenharCaminhao posiciona o caminhao e o rastro percorrido no ponto
// atual do percurso.
function desenharCaminhao() {
  if (!anim.segs.length) {
    mapa.esconderCaminhao();
    mapa.renderPercorrido([]);
    return;
  }
  let p;
  if (anim.segIndex >= anim.segs.length) {
    const ult = anim.segs[anim.segs.length - 1];
    p = ult.b;
  } else {
    const seg = anim.segs[anim.segIndex];
    const t = seg.len > 0 ? Math.min(anim.segDist / seg.len, 1) : 0;
    p = { x: seg.a.x + (seg.b.x - seg.a.x) * t, y: seg.a.y + (seg.b.y - seg.a.y) * t };
    anim.ultimoAngulo = seg.ang;
  }
  mapa.posicionarCaminhao(p.x, p.y, anim.ultimoAngulo);

  // rastro: vertices dos segmentos ja completados + posicao atual
  const rastro = [];
  for (let i = 0; i < Math.min(anim.segIndex, anim.segs.length); i++) rastro.push(anim.segs[i].a);
  if (anim.segIndex < anim.segs.length) rastro.push(anim.segs[anim.segIndex].a);
  rastro.push(p);
  mapa.renderPercorrido(rastro);
}

// finalizar pausa no deposito ao completar o ciclo.
function finalizar() {
  anim.segIndex = anim.segs.length;
  if (anim.playing) {
    anim.playing = false;
    atualizarBotaoPlay();
    toast("Entregas concluídas — caminhão de volta à usina de Angicos 🥛🎉", "sucesso");
  }
  atualizarProgresso();
}

// ---------- controles ----------

function play() {
  if (!anim.segs.length) return;
  if (anim.segIndex >= anim.segs.length) resetar(); // replay do inicio
  anim.playing = true;
  atualizarBotaoPlay();
}

function pausar() {
  anim.playing = false;
  atualizarBotaoPlay();
}

function alternarPlay() {
  anim.playing ? pausar() : play();
}

function resetar() {
  anim.playing = false;
  anim.segIndex = 0;
  anim.segDist = 0;
  atualizarBotaoPlay();
  desenharCaminhao();
  atualizarProgresso();
  marcarTimeline(0);
}

function passo(delta) {
  pausar();
  irPara(pernaAtual() + delta);
}

// irPara posiciona o caminhao no INICIO da perna k (k cidades visitadas).
function irPara(k) {
  if (!anim.segs.length) return;
  k = Math.max(0, Math.min(k, anim.totalPernas));
  anim.segIndex = k >= anim.totalPernas ? anim.segs.length : anim.inicioPerna[k];
  anim.segDist = 0;
  desenharCaminhao();
  atualizarProgresso();
  marcarTimeline(k);
}

function habilitarControles(ativo) {
  for (const id of ["btn-play", "btn-reset", "btn-passo-tras", "btn-passo-frente", "scrubber"]) {
    $(id).disabled = !ativo;
  }
}

function atualizarBotaoPlay() {
  $("btn-play").innerHTML = anim.playing ? "⏸&nbsp; Pausar" : "▶&nbsp; Percorrer";
}

function atualizarProgresso() {
  const k = pernaAtual();
  $("scrubber").value = k;
  $("prog-cidades").textContent = `${k} / ${anim.totalPernas}`;
}

// construirTimeline monta um chip por parada na ordem de visita + o chip
// final da volta ao deposito; clicar num chip leva o caminhao ate ele.
function construirTimeline() {
  const tl = $("timeline");
  tl.replaceChildren();
  for (let i = 0; i < anim.ordem.length; i++) {
    const id = anim.ordem[i];
    tl.appendChild(chipTimeline(i, `${i + 1} · ${state.nomePorId.get(id) || id}`));
  }
  if (anim.ordem.length) {
    const dep = anim.ordem[0];
    tl.appendChild(chipTimeline(anim.ordem.length, `↺ ${state.nomePorId.get(dep) || dep}`));
  }
}

function chipTimeline(k, texto) {
  const b = document.createElement("button");
  b.className = "chip";
  b.textContent = texto;
  b.title = "Levar o caminhão até esta parada";
  b.addEventListener("click", () => {
    pausar();
    irPara(k);
  });
  return b;
}

// marcarTimeline pinta os chips ja visitados e destaca a proxima parada.
function marcarTimeline(k) {
  const chips = $("timeline").children;
  for (let i = 0; i < chips.length; i++) {
    chips[i].classList.toggle("feito", i < k);
    chips[i].classList.toggle("atual", i === k);
  }
  const atual = chips[k];
  if (atual && atual.scrollIntoView) {
    atual.scrollIntoView({ block: "nearest", inline: "nearest" });
  }
}

// ---------- paineis ----------

function limparResultado() {
  $("resultado").innerHTML = '<p class="vazio">Escolha um método para resolver.</p>';
  $("rota-texto").textContent = "—";
}

function preencherResultado(cfg, resp) {
  const otimo = state.instancia.otimo;
  const linhas = [
    ["Método", cfg ? cfg.nome : resp.metodo],
    ["Custo", resp.rota.custo.toFixed(2)],
    ["Ótimo", otimo ? otimo.toFixed(2) : "—"],
    ["Gap", `${resp.gapPercent.toFixed(2)}%`],
    ["Tempo", `${resp.tempoMs.toFixed(2)} ms`],
  ];
  if (resp.antes) {
    linhas.push(["Antes → depois", `${resp.antes.custo.toFixed(2)} → ${resp.rota.custo.toFixed(2)}`]);
  }
  // quilometragem rodoviaria real (OSRM) — informativa; 
  const kmReal = MapaReal.kmDaOrdem(resp.rota.ordem);
  if (kmReal) {
    linhas.push(["Estrada (OSRM)", `${kmReal.toFixed(1)} km`]);
  }

  const dl = document.createElement("dl");
  for (const [rotulo, valor] of linhas) {
    const dt = document.createElement("dt");
    dt.textContent = rotulo;
    const dd = document.createElement("dd");
    dd.textContent = valor;
    dl.append(dt, dd);
  }
  $("resultado").replaceChildren(dl);

  // exibicao textual repete o deposito no fim para explicitar o ciclo
  const ordem = resp.rota.ordem;
  $("rota-texto").textContent = ordem.length ? [...ordem, ordem[0]].join(" → ") : "—";
}

// ---------- toast ----------

let toastTimer = null;
function toast(msg, tipo) {
  const t = $("toast");
  t.textContent = msg;
  // sucesso pinta o aviso de verde; o padrao (erro) segue vermelho
  t.classList.toggle("sucesso", tipo === "sucesso");
  t.classList.add("visivel");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => t.classList.remove("visivel"), 4000);
}

// ---------- inicializacao ----------

document.addEventListener("DOMContentLoaded", async () => {
  MapaPCV.init($("mapa"));

  // botoes de metodo gerados a partir do array de configuracao
  const cont = $("metodos");
  for (const m of METODOS) {
    const b = document.createElement("button");
    b.className = "metodo-btn";
    b.textContent = m.rotulo;
    b.title = m.nome;
    b.addEventListener("click", () => resolver(m.valor));
    cont.appendChild(b);
  }

  $("instancia").addEventListener("change", (e) => selecionarInstancia(e.target.value));
  $("toggle-antes").addEventListener("change", (e) => {
    document.body.classList.toggle("mostrar-antes", e.target.checked);
    // o CSS so alcanca o SVG; no mapa real a camada e re-renderizada
    if (mapa === MapaReal && state.ultimaResp) {
      mapa.renderRotaAntes(
        e.target.checked && state.ultimaResp.antes ? state.ultimaResp.antes.ordem : []
      );
    }
  });
  $("toggle-mapa-real").addEventListener("change", (e) => usarMapaReal(e.target.checked));
  if (!MapaReal.disponivel()) {
    $("toggle-mapa-real").disabled = true;
    $("toggle-mapa-real").parentElement.title = "Indisponível sem internet (Leaflet/tiles via CDN)";
  }
  $("toggle-numeros").addEventListener("change", (e) => {
    document.body.classList.toggle("mostrar-numeros", e.target.checked);
  });

  // controles do caminhao
  $("btn-play").addEventListener("click", alternarPlay);
  $("btn-reset").addEventListener("click", resetar);
  $("btn-passo-tras").addEventListener("click", () => passo(-1));
  $("btn-passo-frente").addEventListener("click", () => passo(1));
  $("velocidade").addEventListener("input", (e) => {
    anim.speed = parseFloat(e.target.value);
    $("velocidade-val").textContent = `${anim.speed.toFixed(2)}×`;
  });
  $("scrubber").addEventListener("input", (e) => {
    pausar();
    irPara(parseInt(e.target.value, 10) || 0);
  });
  habilitarControles(false);
  requestAnimationFrame(laco);


  await MapaReal.carregarRotas();

  try {
    await carregarInstancias();
  } catch (e) {
    toast("Erro ao listar instâncias: " + e.message);
  }

  const params = new URLSearchParams(location.search);
  if (MapaReal.disponivel() && params.get("mapa") !== "mds") {
    $("toggle-mapa-real").checked = true;
    usarMapaReal(true);
  }
  const metodoAuto = params.get("metodo");
  if (metodoAuto && METODOS.some((m) => m.valor === metodoAuto)) {
    resolver(metodoAuto);
  }
});
