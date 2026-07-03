// render.js — desenho SVG do front U3.
//
// Objeto global MapaPCV, sem fetch e sem conhecimento de endpoints: recebe
// pontos (coords MDS normalizadas em [0.05, 0.95]) e rotas em ids originais
// e desenha no SVG 1000x680. Camadas (de tras para frente): rota-antes,
// rota, cidades, badges, caminhao.

const MapaPCV = (() => {
  const VBOX_W = 1000;
  const VBOX_H = 680;
  const NS = "http://www.w3.org/2000/svg";
  const RAIO_CIDADE = 9; // raio do circulo da cidade, em unidades da viewBox

  let svg = null;
  let camadaRotaAntes = null;
  let camadaRota = null;
  let camadaPercorrido = null;
  let camadaCidades = null;
  let camadaBadges = null;
  let camadaCaminhao = null;
  let caminhao = null;
  let pontoPorId = new Map(); // id original -> {x, y} em pixels da viewBox

  // el cria um elemento SVG com atributos.
  function el(tag, attrs) {
    const e = document.createElementNS(NS, tag);
    for (const [k, v] of Object.entries(attrs || {})) e.setAttribute(k, v);
    return e;
  }

  // escala converte coordenada normalizada [0,1] para pixels da viewBox.
  function escala(p) {
    return { x: p.x * VBOX_W, y: p.y * VBOX_H };
  }

  // init prende o MapaPCV a um <svg> e cria as camadas na ordem certa.
  function init(svgEl) {
    svg = svgEl;
    svg.replaceChildren();
    const defs = el("defs");
    defs.innerHTML =
      '<filter id="sombra-pcv" x="-40%" y="-40%" width="180%" height="180%">' +
      '<feDropShadow dx="0" dy="2" stdDeviation="2.5" flood-color="#03101d" flood-opacity="0.5"/>' +
      "</filter>";
    svg.appendChild(defs);
    camadaRotaAntes = el("g", { class: "camada-rota-antes" });
    camadaRota = el("g", { class: "camada-rota" });
    camadaPercorrido = el("g", { class: "camada-percorrido" });
    camadaCidades = el("g", { class: "camada-cidades" });
    camadaBadges = el("g", { class: "camada-badges" });
    camadaCaminhao = el("g", { class: "camada-caminhao" });
    svg.append(camadaRotaAntes, camadaRota, camadaPercorrido, camadaCidades, camadaBadges, camadaCaminhao);
    criaCaminhao();
  }

  // criaCaminhao monta o caminhao visto de cima apontando para +x (mesmo
  // desenho do front da U2, reduzido para a escala dos nos da U3); comeca
  // invisivel — o app.js o posiciona ao animar a rota.
  function criaCaminhao() {
    camadaCaminhao.replaceChildren();
    const g = el("g", { class: "caminhao", filter: "url(#sombra-pcv)", visibility: "hidden" });
    g.innerHTML =
      '<g transform="scale(0.72)">' +
      '<circle class="caminhao-halo" r="22"/>' +
      '<rect x="-20" y="-11" width="40" height="22" rx="5" class="caminhao-corpo"/>' +
      '<rect x="-19" y="-9" width="13" height="18" rx="2" class="caminhao-bau"/>' +
      '<rect x="7" y="-9" width="12" height="18" rx="3" class="caminhao-cabine"/>' +
      '<rect x="11" y="-6" width="6" height="12" rx="1.5" class="caminhao-vidro"/>' +
      '<circle cx="19.5" cy="-6" r="1.7" class="caminhao-farol"/>' +
      '<circle cx="19.5" cy="6" r="1.7" class="caminhao-farol"/>' +
      '<text x="-12.5" y="0" text-anchor="middle" dy="0.35em" class="caminhao-simbolo">🥛</text>' +
      "</g>";
    camadaCaminhao.appendChild(g);
    caminhao = g;
  }

  // renderPontos desenha as cidades (circulo + id + nome) e destaca o
  // deposito. pontos = [{id, nome, x, y}] com x/y normalizados; deposito =
  // id original do deposito (1). Limpa rotas, badges e caminhao anteriores.
  function renderPontos(pontos, deposito) {
    limparRotas();
    camadaCidades.replaceChildren();
    pontoPorId = new Map();

    for (const p of pontos) {
      const { x, y } = escala(p);
      pontoPorId.set(p.id, { x, y });

      const g = el("g", {
        class: p.id === deposito ? "cidade deposito" : "cidade",
        transform: `translate(${x} ${y})`,
      });
      g.appendChild(el("circle", { r: RAIO_CIDADE }));

      const tId = el("text", { class: "cidade-id" });
      tId.textContent = p.id;
      g.appendChild(tId);

      const tNome = el("text", { class: "cidade-nome", y: 19 });
      tNome.textContent = p.nome;
      g.appendChild(tNome);

      const title = el("title");
      title.textContent = `${p.id} — ${p.nome}`;
      g.appendChild(title);

      camadaCidades.appendChild(g);
    }
  }

  // caminhoFechado monta o atributo d de um <path> ligando os ids na ordem
  // dada e fechando o ciclo de volta ao primeiro ponto (Z).
  function caminhoFechado(ordem) {
    const partes = [];
    for (let i = 0; i < ordem.length; i++) {
      const p = pontoPorId.get(ordem[i]);
      if (!p) continue;
      partes.push(`${i === 0 ? "M" : "L"} ${p.x} ${p.y}`);
    }
    if (partes.length) partes.push("Z"); // fecha o ciclo (volta ao deposito)
    return partes.join(" ");
  }

  // pontosDaOrdem devolve os pontos da rota em pixels da viewBox, na ordem
  // de visita (sem repetir o deposito).
  function pontosDaOrdem(ordem) {
    const pts = [];
    for (const id of ordem || []) {
      const p = pontoPorId.get(id);
      if (p) pts.push(p);
    }
    return pts;
  }

  // caminhosDaOrdem devolve uma polilinha por PERNA do ciclo (cidade k ->
  // cidade k+1, fechando no deposito). No diagrama cada perna e um trecho
  // reto de 2 pontos; o motor de animacao do app.js consome este formato.
  function caminhosDaOrdem(ordem) {
    const pts = pontosDaOrdem(ordem);
    const pernas = [];
    for (let i = 0; i < pts.length && pts.length >= 2; i++) {
      pernas.push([pts[i], pts[(i + 1) % pts.length]]);
    }
    return pernas;
  }

  // renderRota desenha a rota final como ciclo fechado + badges com a ordem
  // de visita (1 = deposito). ordem = ids originais SEM repetir o deposito
  // no fim (o fechamento e responsabilidade do desenho).
  function renderRota(ordem) {
    camadaRota.replaceChildren();
    camadaBadges.replaceChildren();
    if (!ordem || !ordem.length) return;

    camadaRota.appendChild(el("path", { class: "rota-path", d: caminhoFechado(ordem) }));

    // badge deslocado do centro com a posicao de visita
    for (let i = 0; i < ordem.length; i++) {
      const p = pontoPorId.get(ordem[i]);
      if (!p) continue;
      const g = el("g", { class: "badge", transform: `translate(${p.x + 11} ${p.y - 11})` });
      g.appendChild(el("circle", { r: 5.5 }));
      const t = el("text");
      t.textContent = i + 1;
      g.appendChild(t);
      camadaBadges.appendChild(g);
    }
  }

  // renderRotaAntes desenha a rota do construtivo (antes da busca local)
  // tracejada, atras da rota final. ordem no mesmo formato de renderRota.
  function renderRotaAntes(ordem) {
    camadaRotaAntes.replaceChildren();
    if (!ordem || !ordem.length) return;
    camadaRotaAntes.appendChild(
      el("path", { class: "rota-antes-path", d: caminhoFechado(ordem) })
    );
  }

  // renderPercorrido desenha o rastro do trecho ja percorrido pelo caminhao
  // como polilinha ABERTA sobre a rota. pts = pontos em pixels da viewBox
  // (vertices completados + posicao atual do caminhao);
  function renderPercorrido(pts) {
    camadaPercorrido.replaceChildren();
    if (!pts || pts.length < 2) return;
    const partes = pts.map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`);
    camadaPercorrido.appendChild(
      el("path", { class: "rota-percorrida", d: partes.join(" ") })
    );
  }

  // posicionarCaminhao mostra o caminhao em (x, y) da viewBox apontando na
  // direcao anguloGraus (0 = +x).
  function posicionarCaminhao(x, y, anguloGraus) {
    if (!caminhao) return;
    caminhao.setAttribute("visibility", "visible");
    caminhao.setAttribute("transform", `translate(${x} ${y}) rotate(${anguloGraus})`);
  }

  // esconderCaminhao oculta o caminhao (rota limpa ou animacao desligada).
  function esconderCaminhao() {
    if (caminhao) caminhao.setAttribute("visibility", "hidden");
  }

  // limparRotas remove as rotas, o rastro, a numeracao e o
  // caminhao, mantendo as cidades desenhadas.
  function limparRotas() {
    camadaRotaAntes.replaceChildren();
    camadaRota.replaceChildren();
    camadaPercorrido.replaceChildren();
    camadaBadges.replaceChildren();
    esconderCaminhao();
  }

  return {
    init,
    renderPontos,
    renderRota,
    renderRotaAntes,
    renderPercorrido,
    limparRotas,
    pontosDaOrdem,
    caminhosDaOrdem,
    posicionarCaminhao,
    esconderCaminhao,
    velocidadeBase: 150,
  };
})();
