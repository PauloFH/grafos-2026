// layout.js — posicionamento dos cruzamentos no mapa.
//
// O dataset nao traz coordenadas, entao geramos um traçado de "bairro" com um
// algoritmo force-directed (Fruchterman-Reingold simplificado). E DETERMINISTICO:
// a posicao inicial vem do indice do vertice (circulo), sem Math.random, para
// que o mesmo grafo gere sempre o mesmo mapa.

function computeLayout(vertices, arestas, width, height) {
  const n = vertices.length;
  const pos = {};
  if (n === 0) return pos;

  const cx = width / 2;
  const cy = height / 2;
  const raio = Math.min(width, height) * 0.4;

  // Init determinístico: distribui em círculo pela ordem do vértice.
  vertices.forEach((v, i) => {
    const ang = (2 * Math.PI * i) / n;
    pos[v] = { x: cx + raio * Math.cos(ang), y: cy + raio * Math.sin(ang) };
  });

  if (n === 1) return pos;

  // Pares não-direcionados (ignora sentido) para a força de mola.
  const pares = arestas.map((e) => [e.de, e.para]);

  const area = width * height;
  const k = Math.sqrt(area / n) * 0.55; // distância ideal entre vértices
  const ITER = 500;
  let temp = width * 0.12;
  const cool = temp / (ITER + 1);

  for (let it = 0; it < ITER; it++) {
    const disp = {};
    for (const v of vertices) disp[v] = { x: 0, y: 0 };

    // Repulsão entre todos os pares (Coulomb): k^2 / d.
    for (let i = 0; i < n; i++) {
      for (let j = i + 1; j < n; j++) {
        const a = vertices[i];
        const b = vertices[j];
        let dx = pos[a].x - pos[b].x;
        let dy = pos[a].y - pos[b].y;
        let d = Math.hypot(dx, dy) || 0.01;
        const f = (k * k) / d;
        const ux = dx / d;
        const uy = dy / d;
        disp[a].x += ux * f;
        disp[a].y += uy * f;
        disp[b].x -= ux * f;
        disp[b].y -= uy * f;
      }
    }

    // Atração ao longo das arestas (mola): d^2 / k.
    for (const [a, b] of pares) {
      let dx = pos[a].x - pos[b].x;
      let dy = pos[a].y - pos[b].y;
      let d = Math.hypot(dx, dy) || 0.01;
      const f = (d * d) / k;
      const ux = dx / d;
      const uy = dy / d;
      disp[a].x -= ux * f;
      disp[a].y -= uy * f;
      disp[b].x += ux * f;
      disp[b].y += uy * f;
    }

    // Gravidade fraca ao centro (evita deriva de componentes frouxas).
    for (const v of vertices) {
      disp[v].x += (cx - pos[v].x) * 0.012;
      disp[v].y += (cy - pos[v].y) * 0.012;
    }

    // Aplica deslocamento limitado pela "temperatura" (resfriamento).
    for (const v of vertices) {
      const dd = Math.hypot(disp[v].x, disp[v].y) || 0.01;
      const lim = Math.min(dd, temp);
      pos[v].x += (disp[v].x / dd) * lim;
      pos[v].y += (disp[v].y / dd) * lim;
    }
    temp -= cool;
  }

  return normaliza(pos, vertices, width, height, 70);
}

// normaliza reescala/centraliza o layout para caber na viewBox com padding.
function normaliza(pos, vertices, width, height, pad) {
  let minx = Infinity;
  let miny = Infinity;
  let maxx = -Infinity;
  let maxy = -Infinity;
  for (const v of vertices) {
    minx = Math.min(minx, pos[v].x);
    miny = Math.min(miny, pos[v].y);
    maxx = Math.max(maxx, pos[v].x);
    maxy = Math.max(maxy, pos[v].y);
  }
  const sx = (width - 2 * pad) / ((maxx - minx) || 1);
  const sy = (height - 2 * pad) / ((maxy - miny) || 1);
  const s = Math.min(sx, sy);
  const offx = pad + (width - 2 * pad - (maxx - minx) * s) / 2;
  const offy = pad + (height - 2 * pad - (maxy - miny) * s) / 2;
  for (const v of vertices) {
    pos[v] = {
      x: offx + (pos[v].x - minx) * s,
      y: offy + (pos[v].y - miny) * s,
    };
  }
  return pos;
}
