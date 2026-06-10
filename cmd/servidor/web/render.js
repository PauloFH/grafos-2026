// render.js — desenho do mapa em SVG e helpers de geometria.
//
// Expoe:
//   - helpers globais de geometria (chaveAresta, controlePonto, quad*, etc.)
//     usados tambem pelo app.js para animar o caminhao;
//   - o objeto `Mapa`, que cria/atualiza os elementos SVG (ruas, cruzamentos,
//     deposito e caminhao) e dispara callbacks de clique/arraste nos nos.

const SVG_NS = "http://www.w3.org/2000/svg";
const RAIO_NO = 17; // raio do cruzamento, em unidades da viewBox

function el(tag, attrs) {
  const e = document.createElementNS(SVG_NS, tag);
  for (const k in attrs) e.setAttribute(k, attrs[k]);
  return e;
}

// menorRotulo replica a ordenacao do backend: numerica quando possivel.
function menorRotulo(a, b) {
  const na = Number(a);
  const nb = Number(b);
  if (!Number.isNaN(na) && !Number.isNaN(nb)) return na < nb;
  return a < b;
}

// chaveAresta identifica uma rua de forma estavel. Casa com a deduplicacao do
// backend: par canonico (lo|hi) para grafo; sentido preservado (de>para) p/ digrafo.
function chaveAresta(de, para, direcionado) {
  if (direcionado) return de + ">" + para;
  return menorRotulo(de, para) ? de + "|" + para : para + "|" + de;
}

// controlePonto: ponto de controle da curva quadratica de uma rua.
// Grafo (reta) -> ponto medio. Digrafo -> deslocado para a esquerda do sentido,
// para separar arcos opostos (1->20 e 20->1).
function controlePonto(a, b, direcionado) {
  const mx = (a.x + b.x) / 2;
  const my = (a.y + b.y) / 2;
  if (!direcionado) return { x: mx, y: my };
  const dx = b.x - a.x;
  const dy = b.y - a.y;
  const d = Math.hypot(dx, dy) || 0.01;
  const nx = -dy / d;
  const ny = dx / d;
  const off = 26;
  return { x: mx + nx * off, y: my + ny * off };
}

function quadAt(a, c, b, t) {
  const mt = 1 - t;
  return {
    x: mt * mt * a.x + 2 * mt * t * c.x + t * t * b.x,
    y: mt * mt * a.y + 2 * mt * t * c.y + t * t * b.y,
  };
}

function quadDeriv(a, c, b, t) {
  const mt = 1 - t;
  return {
    x: 2 * mt * (c.x - a.x) + 2 * t * (b.x - c.x),
    y: 2 * mt * (c.y - a.y) + 2 * t * (b.y - c.y),
  };
}

function quadComprimento(a, c, b) {
  let len = 0;
  let prev = a;
  const N = 16;
  for (let i = 1; i <= N; i++) {
    const p = quadAt(a, c, b, i / N);
    len += Math.hypot(p.x - prev.x, p.y - prev.y);
    prev = p;
  }
  return len;
}

function aparar(p, alvo, r) {
  const dx = alvo.x - p.x;
  const dy = alvo.y - p.y;
  const d = Math.hypot(dx, dy) || 0.01;
  return { x: p.x + (dx / d) * r, y: p.y + (dy / d) * r };
}

const Mapa = {
  svg: null,
  defs: null,
  gArestas: null,
  gNos: null,
  gCaminhao: null,
  positions: null,
  direcionado: false,
  inicio: null,
  arestaEls: {}, // chave -> { grupo, forma, rotulo, de, para }
  noEls: {}, // id -> { grupo, circulo }
  caminhao: null,
  handlers: {},

  init(svg, handlers) {
    this.svg = svg;
    this.handlers = handlers || {};
    this.defs = el("defs", {});
    this.defs.innerHTML = this._defsHTML();
    this.svg.appendChild(this.defs);
    this.gArestas = el("g", { class: "camada-arestas" });
    this.gNos = el("g", { class: "camada-nos" });
    this.gCaminhao = el("g", { class: "camada-caminhao" });
    this.svg.appendChild(this.gArestas);
    this.svg.appendChild(this.gNos);
    this.svg.appendChild(this.gCaminhao);
    this._ligarArraste();
  },

  _defsHTML() {
    const seta = (id, cls) =>
      `<marker id="${id}" class="${cls}" viewBox="0 0 10 10" refX="9" refY="5"
         markerWidth="7" markerHeight="7" orient="auto-start-reverse">
         <path d="M0,0 L10,5 L0,10 z"/></marker>`;
    return (
      seta("seta", "m-seta") +
      seta("seta-ok", "m-seta-ok") +
      seta("seta-atual", "m-seta-atual") +
      `<filter id="sombra" x="-40%" y="-40%" width="180%" height="180%">
         <feDropShadow dx="0" dy="2" stdDeviation="2.5" flood-color="#0b3d1f" flood-opacity="0.45"/>
       </filter>`
    );
  },

  // render reconstroi o mapa inteiro a partir do grafo e das posicoes.
  render(graph, positions) {
    this.positions = positions;
    this.direcionado = graph.direcionado;
    this.inicio = null;
    this.arestaEls = {};
    this.noEls = {};
    this.gArestas.replaceChildren();
    this.gNos.replaceChildren();
    this.gCaminhao.replaceChildren();

    for (const e of graph.arestas) this._criaAresta(e.de, e.para);
    for (const v of graph.vertices) this._criaNo(v, graph.graus[v]);
    this._criaCaminhao();
  },

  _criaAresta(de, para) {
    const chave = chaveAresta(de, para, this.direcionado);
    const grupo = el("g", { class: "aresta", "data-chave": chave });
    const forma = el(this.direcionado ? "path" : "line", { class: "rua" });
    if (this.direcionado) forma.setAttribute("marker-end", "url(#seta)");
    const rotulo = el("text", { class: "rua-ordem", "text-anchor": "middle", dy: "0.32em" });
    grupo.appendChild(forma);
    grupo.appendChild(rotulo);
    this.gArestas.appendChild(grupo);
    this.arestaEls[chave] = { grupo, forma, rotulo, de, para };
    this._geoAresta(chave);
  },

  _geoAresta(chave) {
    const a = this.positions[this.arestaEls[chave].de];
    const b = this.positions[this.arestaEls[chave].para];
    const { forma, rotulo } = this.arestaEls[chave];
    const c = controlePonto(a, b, this.direcionado);
    if (this.direcionado) {
      const aT = aparar(a, c, RAIO_NO);
      const bT = aparar(b, c, RAIO_NO + 4);
      forma.setAttribute("d", `M ${aT.x} ${aT.y} Q ${c.x} ${c.y} ${bT.x} ${bT.y}`);
    } else {
      forma.setAttribute("x1", a.x);
      forma.setAttribute("y1", a.y);
      forma.setAttribute("x2", b.x);
      forma.setAttribute("y2", b.y);
    }
    rotulo.setAttribute("x", c.x);
    rotulo.setAttribute("y", c.y);
  },

  _criaNo(id, grau) {
    const p = this.positions[id];
    const grupo = el("g", {
      class: "no",
      "data-id": id,
      transform: `translate(${p.x} ${p.y})`,
    });
    const circ = el("circle", { class: "no-circulo", r: RAIO_NO });
    const txt = el("text", { class: "no-rotulo", "text-anchor": "middle", dy: "0.34em" });
    txt.textContent = id;
    const titulo = el("title", {});
    titulo.textContent = `Cruzamento ${id} — grau ${grau ?? "?"}`;
    grupo.appendChild(circ);
    grupo.appendChild(txt);
    grupo.appendChild(titulo);
    this.gNos.appendChild(grupo);
    this.noEls[id] = { grupo, circ };
  },

  _criaCaminhao() {
    const g = el("g", { class: "caminhao", filter: "url(#sombra)", visibility: "hidden" });
    // Caminhao visto de cima, apontando para +x (gira sem ficar de cabeca p/ baixo).
    g.innerHTML = `
      <circle class="caminhao-halo" r="22"/>
      <rect x="-20" y="-11" width="40" height="22" rx="5" class="caminhao-corpo"/>
      <rect x="-19" y="-9" width="13" height="18" rx="2" class="caminhao-cacamba"/>
      <rect x="7" y="-9" width="12" height="18" rx="3" class="caminhao-cabine"/>
      <rect x="11" y="-6" width="6" height="12" rx="1.5" class="caminhao-vidro"/>
      <circle cx="19.5" cy="-6" r="1.7" class="caminhao-farol"/>
      <circle cx="19.5" cy="6" r="1.7" class="caminhao-farol"/>
      <text x="-12.5" y="0" text-anchor="middle" dy="0.35em" class="caminhao-simbolo">♻</text>`;
    this.gCaminhao.appendChild(g);
    this.caminhao = g;
  },

  // atualizarGeometria reposiciona arestas e nos apos arrastar um cruzamento.
  atualizarGeometria() {
    for (const chave in this.arestaEls) this._geoAresta(chave);
    for (const id in this.noEls) {
      const p = this.positions[id];
      this.noEls[id].grupo.setAttribute("transform", `translate(${p.x} ${p.y})`);
    }
  },

  definirInicio(id) {
    if (this.inicio && this.noEls[this.inicio]) {
      this.noEls[this.inicio].grupo.classList.remove("deposito");
    }
    this.inicio = id;
    if (id && this.noEls[id]) this.noEls[id].grupo.classList.add("deposito");
  },

  limparColetadas() {
    for (const chave in this.arestaEls) {
      const a = this.arestaEls[chave];
      a.grupo.classList.remove("coletada", "atual");
      a.rotulo.textContent = "";
      if (this.direcionado) a.forma.setAttribute("marker-end", "url(#seta)");
    }
  },

  marcarColetada(chave, ordem) {
    const a = this.arestaEls[chave];
    if (!a) return;
    a.grupo.classList.add("coletada");
    a.grupo.classList.remove("atual");
    a.rotulo.textContent = ordem;
    if (this.direcionado) a.forma.setAttribute("marker-end", "url(#seta-ok)");
  },

  destacarAtual(chave) {
    for (const k in this.arestaEls) {
      const a = this.arestaEls[k];
      const ativo = k === chave;
      a.grupo.classList.toggle("atual", ativo);
      if (this.direcionado && !a.grupo.classList.contains("coletada")) {
        a.forma.setAttribute("marker-end", ativo ? "url(#seta-atual)" : "url(#seta)");
      }
    }
  },

  posicionarCaminhao(x, y, anguloGraus) {
    if (!this.caminhao) return;
    this.caminhao.setAttribute("visibility", "visible");
    this.caminhao.setAttribute("transform", `translate(${x} ${y}) rotate(${anguloGraus})`);
  },

  esconderCaminhao() {
    if (this.caminhao) this.caminhao.setAttribute("visibility", "hidden");
  },

  pulsarNo(id) {
    const n = this.noEls[id];
    if (!n) return;
    n.grupo.classList.remove("pulso");
    // reinicia a animacao
    void n.grupo.getBBox();
    n.grupo.classList.add("pulso");
  },

  // --- arraste / clique nos cruzamentos ---
  _ligarArraste() {
    let alvo = null;
    let moveu = false;
    let p0 = null;

    const paraSVG = (evt) => {
      const pt = this.svg.createSVGPoint();
      pt.x = evt.clientX;
      pt.y = evt.clientY;
      return pt.matrixTransform(this.svg.getScreenCTM().inverse());
    };

    this.gNos.addEventListener("pointerdown", (evt) => {
      const g = evt.target.closest(".no");
      if (!g) return;
      alvo = g.getAttribute("data-id");
      moveu = false;
      p0 = paraSVG(evt);
      g.setPointerCapture?.(evt.pointerId);
      evt.preventDefault();
    });

    this.gNos.addEventListener("pointermove", (evt) => {
      if (!alvo) return;
      const p = paraSVG(evt);
      if (!moveu && Math.hypot(p.x - p0.x, p.y - p0.y) < 4) return;
      moveu = true;
      this.positions[alvo] = { x: p.x, y: p.y };
      this.atualizarGeometria();
      this.handlers.onArrastar && this.handlers.onArrastar(alvo);
    });

    const fim = () => {
      if (alvo && !moveu) this.handlers.onClicar && this.handlers.onClicar(alvo);
      alvo = null;
    };
    this.gNos.addEventListener("pointerup", fim);
    this.gNos.addEventListener("pointercancel", fim);
  },
};
