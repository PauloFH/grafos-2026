// mapa-real.js — renderer de mapa real (Leaflet + tiles CARTO/OSM) do front
// U3, com a MESMA interface do MapaPCV (render.js): o app.js alterna entre
// os dois sem mudar a logica de animacao.
//
// Espaco de coordenadas exposto ao app: {x: longitude, y: -latitude} — o y
// invertido preserva a convencao de tela (y cresce para baixo) usada pelo
// motor de animacao, entao angulos e interpolacao funcionam sem ajuste.
//
// Depende de: Leaflet global `L` (CDN; sem internet o modo fica indisponivel
// e o front segue no diagrama MDS) e COORDS_RN (coords-rn.js).

const MapaReal = (() => {
  const VEL_BASE = 0.5; // graus/s — percorre o RN inteiro em ~20 s a 1.5x

  let map = null;
  let grupoCidades = null;
  let rotaLinha = null;
  let rotaAntesLinha = null;
  let percorridoLinha = null;
  let caminhaoMarker = null;
  let pontoPorId = new Map(); // id original -> {x: lon, y: -lat}
  let markerPorId = new Map();
  let nomePorId = new Map();
  let depositoId = null;
  let paresRodovia = null; // rotas-rn.json: "min|max" -> {km, p: [[lat,lon],...]}

  // carregarRotas busca o cache de geometrias rodoviarias (gerado por
  // scripts/gera_rotas_reais.py via OSRM). Sem o cache, as pernas caem em
  // linha reta — o front continua funcionando.
  async function carregarRotas() {
    if (paresRodovia) return;
    try {
      const r = await fetch("rotas-rn.json");
      if (r.ok) paresRodovia = (await r.json()).pares || null;
    } catch (_) {
      paresRodovia = null; // sem cache: fallback reto
    }
  }

  // geometriaEntre devolve a polilinha rodoviaria idA -> idB no espaco do
  // renderer ({x: lon, y: -lat}), com os pontos exatos das cidades nas
  // extremidades; sem cache do par, devolve o trecho reto.
  function geometriaEntre(idA, idB) {
    const a = pontoPorId.get(idA);
    const b = pontoPorId.get(idB);
    if (!a || !b) return [];
    const chave = idA < idB ? `${idA}|${idB}` : `${idB}|${idA}`;
    const par = paresRodovia && paresRodovia[chave];
    if (!par || !par.p || par.p.length < 2) return [a, b];
    let pts = par.p.map(([lat, lon]) => ({ x: lon, y: -lat }));
    if (idA > idB) pts = pts.slice().reverse();
    return [a, ...pts, b]; // ancora nas cidades (o OSRM gruda na estrada)
  }

  // kmDaOrdem soma a quilometragem rodoviaria real (OSRM) do ciclo completo;
  // null se o cache nao cobre todos os trechos.
  function kmDaOrdem(ordem) {
    if (!paresRodovia || !ordem || ordem.length < 2) return null;
    let total = 0;
    for (let i = 0; i < ordem.length; i++) {
      const a = ordem[i];
      const b = ordem[(i + 1) % ordem.length];
      const par = paresRodovia[a < b ? `${a}|${b}` : `${b}|${a}`];
      if (!par || typeof par.km !== "number") return null;
      total += par.km;
    }
    return total;
  }

  // caminhosDaOrdem devolve uma polilinha por PERNA do ciclo, pela estrada
  // real quando o cache cobre o par (mesmo formato do MapaPCV).
  function caminhosDaOrdem(ordem) {
    const pernas = [];
    for (let i = 0; i < (ordem || []).length && ordem.length >= 2; i++) {
      const geo = geometriaEntre(ordem[i], ordem[(i + 1) % ordem.length]);
      if (geo.length >= 2) pernas.push(geo);
    }
    return pernas;
  }

  // disponivel diz se o Leaflet (CDN) e as coordenadas carregaram.
  function disponivel() {
    return typeof L !== "undefined" && typeof COORDS_RN !== "undefined";
  }

  // init cria o mapa uma unica vez, com tiles escuros (combinam com o tema
  // PLP) e attribution obrigatoria do OSM/CARTO.
  function init(containerId) {
    if (map || !disponivel()) return;
    map = L.map(containerId, { zoomSnap: 0.5 });
    L.tileLayer("https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png", {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/attributions">CARTO</a>',
      subdomains: "abcd",
      minZoom: 6,
      maxZoom: 13,
    }).addTo(map);
    grupoCidades = L.layerGroup().addTo(map);
    map.setView([-5.8, -36.4], 8); // centro aproximado do RN

    // contorno oficial do RN (malha estadual do IBGE, embedada) — desenhado
    // antes das rotas/cidades para ficar por baixo delas
    fetch("rn-contorno.json")
      .then((r) => (r.ok ? r.json() : null))
      .then((gj) => {
        if (!gj || !map) return;
        L.geoJSON(gj, {
          interactive: false,
          style: {
            color: "#f2edd5",
            weight: 2,
            opacity: 0.55,
            fillColor: "#ee4c9b",
            fillOpacity: 0.04,
          },
        }).addTo(map);
      })
      .catch(() => {}); // sem contorno o mapa segue normal
  }

  // ativar deve ser chamado quando o container sai de display:none — o
  // Leaflet precisa recalcular o tamanho e reenquadrar os pontos.
  function ativar() {
    if (!map) return;
    map.invalidateSize();
    enquadrar();
  }

  function enquadrar() {
    const pts = [...pontoPorId.values()].map((p) => [-p.y, p.x]);
    if (pts.length) map.fitBounds(pts, { padding: [36, 36] });
  }

  // renderPontos desenha as cidades como circleMarkers nas coordenadas
  // reais (COORDS_RN por nome canonico) e destaca a usina (deposito).
  function renderPontos(pontos, deposito) {
    if (!map) return;
    limparRotas();
    grupoCidades.clearLayers();
    pontoPorId = new Map();
    markerPorId = new Map();
    nomePorId = new Map();
    depositoId = deposito;

    for (const p of pontos) {
      const c = COORDS_RN[p.nome];
      if (!c) continue; // nome fora do dicionario: ponto fica so no diagrama
      const [lat, lon] = c;
      pontoPorId.set(p.id, { x: lon, y: -lat });
      nomePorId.set(p.id, p.nome);

      const ehUsina = p.id === deposito;
      const m = L.circleMarker([lat, lon], {
        radius: ehUsina ? 9 : 6,
        color: ehUsina ? "#f59e0b" : "#ee4c9b",
        weight: 2,
        fillColor: ehUsina ? "#4a3407" : "#3a1440",
        fillOpacity: 0.95,
      }).addTo(grupoCidades);
      m.bindTooltip(ehUsina ? `🏭 ${p.id} · ${p.nome} (usina)` : `${p.id} · ${p.nome}`);
      markerPorId.set(p.id, m);
    }
    enquadrar();
  }

  function latlngsDaOrdem(ordem, fechar) {
    const ll = [];
    for (const id of ordem || []) {
      const p = pontoPorId.get(id);
      if (p) ll.push([-p.y, p.x]);
    }
    if (fechar && ll.length) ll.push(ll[0]);
    return ll;
  }

  // renderRota desenha o ciclo em rosa PLP — pelas estradas reais quando o
  // cache cobre — e anota a posicao de visita no tooltip de cada cidade.
  function renderRota(ordem) {
    if (!map) return;
    if (rotaLinha) map.removeLayer(rotaLinha);
    rotaLinha = null;
    if (!ordem || !ordem.length) return;
    const traco = caminhosDaOrdem(ordem)
      .flat()
      .map((p) => [-p.y, p.x]);
    rotaLinha = L.polyline(traco.length >= 2 ? traco : latlngsDaOrdem(ordem, true), {
      color: "#ee4c9b",
      weight: 3,
      opacity: 0.9,
    }).addTo(map);
    for (let i = 0; i < ordem.length; i++) {
      const m = markerPorId.get(ordem[i]);
      if (!m) continue;
      const nome = nomePorId.get(ordem[i]);
      const pref = ordem[i] === depositoId ? "🏭 " : "";
      m.setTooltipContent(`${pref}${i + 1}º · ${nome}`);
    }
  }

  // renderRotaAntes desenha a rota do construtivo tracejada (cinza), tambem
  // pelas estradas quando o cache cobre.
  function renderRotaAntes(ordem) {
    if (!map) return;
    if (rotaAntesLinha) map.removeLayer(rotaAntesLinha);
    rotaAntesLinha = null;
    if (!ordem || !ordem.length) return;
    const traco = caminhosDaOrdem(ordem)
      .flat()
      .map((p) => [-p.y, p.x]);
    rotaAntesLinha = L.polyline(traco.length >= 2 ? traco : latlngsDaOrdem(ordem, true), {
      color: "#8fa4b5",
      weight: 2.5,
      opacity: 0.5,
      dashArray: "6 5",
    }).addTo(map);
  }

  // renderPercorrido pinta o rastro branco-leite dos trechos ja entregues.
  // pts no espaco {x: lon, y: -lat} do renderer.
  function renderPercorrido(pts) {
    if (!map) return;
    if (percorridoLinha) map.removeLayer(percorridoLinha);
    percorridoLinha = null;
    if (!pts || pts.length < 2) return;
    percorridoLinha = L.polyline(pts.map((p) => [-p.y, p.x]), {
      color: "#f2edd5",
      weight: 5,
      opacity: 0.92,
      lineCap: "round",
    }).addTo(map);
  }

  // pontosDaOrdem devolve os pontos da rota no espaco do renderer.
  function pontosDaOrdem(ordem) {
    const pts = [];
    for (const id of ordem || []) {
      const p = pontoPorId.get(id);
      if (p) pts.push(p);
    }
    return pts;
  }

  function htmlCaminhao() {
    return (
      '<div class="caminhao-real"><svg viewBox="-26 -26 52 52" width="44" height="44">' +
      '<g><circle class="caminhao-halo" r="22"/>' +
      '<rect x="-20" y="-11" width="40" height="22" rx="5" class="caminhao-corpo"/>' +
      '<rect x="-19" y="-9" width="13" height="18" rx="2" class="caminhao-bau"/>' +
      '<rect x="7" y="-9" width="12" height="18" rx="3" class="caminhao-cabine"/>' +
      '<rect x="11" y="-6" width="6" height="12" rx="1.5" class="caminhao-vidro"/>' +
      '<circle cx="19.5" cy="-6" r="1.7" class="caminhao-farol"/>' +
      '<circle cx="19.5" cy="6" r="1.7" class="caminhao-farol"/>' +
      '<text x="-12.5" y="0" text-anchor="middle" dy="0.35em" class="caminhao-simbolo">🥛</text>' +
      "</g></svg></div>"
    );
  }

  // posicionarCaminhao move o marker do caminhao para (x, y) do espaco do
  // renderer e gira o desenho na direcao do movimento.
  function posicionarCaminhao(x, y, anguloGraus) {
    if (!map) return;
    const latlng = [-y, x];
    if (!caminhaoMarker) {
      caminhaoMarker = L.marker(latlng, {
        icon: L.divIcon({ className: "caminhao-real-icone", html: htmlCaminhao(), iconSize: [44, 44], iconAnchor: [22, 22] }),
        interactive: false,
        zIndexOffset: 1000,
      }).addTo(map);
    }
    caminhaoMarker.setLatLng(latlng);
    const el = caminhaoMarker.getElement();
    const rot = el && el.querySelector(".caminhao-real");
    if (rot) rot.style.transform = `rotate(${anguloGraus}deg)`;
  }

  function esconderCaminhao() {
    if (caminhaoMarker && map) {
      map.removeLayer(caminhaoMarker);
      caminhaoMarker = null;
    }
  }

  // limparRotas remove rota, rastro e caminhao, e restaura os tooltips.
  function limparRotas() {
    if (!map) return;
    if (rotaLinha) map.removeLayer(rotaLinha);
    if (rotaAntesLinha) map.removeLayer(rotaAntesLinha);
    if (percorridoLinha) map.removeLayer(percorridoLinha);
    rotaLinha = rotaAntesLinha = percorridoLinha = null;
    esconderCaminhao();
    for (const [id, m] of markerPorId) {
      const nome = nomePorId.get(id);
      m.setTooltipContent(id === depositoId ? `🏭 ${id} · ${nome} (usina)` : `${id} · ${nome}`);
    }
  }

  return {
    disponivel,
    init,
    ativar,
    carregarRotas,
    renderPontos,
    renderRota,
    renderRotaAntes,
    renderPercorrido,
    pontosDaOrdem,
    caminhosDaOrdem,
    kmDaOrdem,
    posicionarCaminhao,
    esconderCaminhao,
    limparRotas,
    velocidadeBase: VEL_BASE,
  };
})();
