(() => {
  const R = {
    "europe-west": { n: "Europe West", iana: "Europe/Lisbon", base: 0, rule: "eu" },
    "europe-central": { n: "Europe Central", iana: "Europe/Berlin", base: 1, rule: "eu" },
    "europe-east": { n: "Europe East", iana: "Europe/Kiev", base: 2, rule: "eu" },
    "us-eastern": { n: "US Eastern", iana: "America/New_York", base: -5, rule: "us" },
    "us-central": { n: "US Central", iana: "America/Chicago", base: -6, rule: "us" },
    "us-mountain": { n: "US Mountain", iana: "America/Denver", base: -7, rule: "us" },
    "us-pacific": { n: "US Pacific", iana: "America/Los_Angeles", base: -8, rule: "us" },
    "australia-sydney": { n: "Australia Sydney", iana: "Australia/Sydney", base: 10, rule: "au" },
    "new-zealand": { n: "New Zealand", iana: "Pacific/Auckland", base: 12, rule: "nz" },
    "chile": { n: "Chile", iana: "America/Santiago", base: -4, rule: "cl" },
    "argentina": { n: "Argentina", iana: "America/Argentina/Buenos_Aires", base: -3, rule: "none" }
  };
  const nth = (y, m, k) => { const d = new Date(Date.UTC(y, m, 1)); let c = 0; for (;;) { if (d.getUTCDay() === 0 && ++c === k) return d; d.setUTCDate(d.getUTCDate() + 1); } };
  const last = (y, m) => { const d = new Date(Date.UTC(y, m + 1, 0)); while (d.getUTCDay()) d.setUTCDate(d.getUTCDate() - 1); return d; };
  const dst = (y, r) => r === "eu" ? [last(y, 2), last(y, 9)] : r === "us" ? [nth(y, 2, 2), nth(y, 10, 1)] :
    r === "au" ? [nth(y, 9, 1), nth(y + 1, 3, 1)] : r === "nz" ? [last(y, 8), nth(y + 1, 3, 1)] :
    r === "cl" ? [nth(y, 8, 1), nth(y + 1, 3, 1)] : null;
  // Mirrors CalculateSmoothOffset in main.go
  function smooth(t, key) {
    const c = R[key]; if (c.rule === "none") return c.base;
    let y = t.getUTCFullYear();
    if (t < dst(y, c.rule)[0] && t < dst(y - 1, c.rule)[1]) y--;
    const [s, e] = dst(y, c.rule);
    if (t >= s && t < e) return c.base + 1 - (t - s) / (e - s);
    let ws = e, we = dst(y + 1, c.rule)[0];
    if (t < s) { ws = dst(y - 1, c.rule)[1]; we = s; }
    return c.base + (t - ws) / (we - ws);
  }
  function official(t, key) {
    const p = new Intl.DateTimeFormat("en-US", { timeZone: R[key].iana, timeZoneName: "longOffset" }).formatToParts(t).find(x => x.type === "timeZoneName").value;
    const m = p.match(/GMT([+-])(\d{2}):?(\d{2})?/); if (!m) return 0;
    return (m[1] === "-" ? -1 : 1) * (+m[2] + (+(m[3] || 0)) / 60);
  }
  const fmtT = t => t.toISOString().substr(11, 8);
  const fmtO = h => (h >= 0 ? "+" : "−") + Math.abs(h).toFixed(4) + " h";
  const $ = s => document.querySelector(s), $$ = s => document.querySelectorAll(s);
  let cur = "europe-central";

  function tick() {
    const now = new Date(), s = smooth(now, cur), o = official(now, cur);
    $("#h-official").textContent = fmtT(new Date(+now + o * 3600e3));
    $("#h-smooth").textContent = fmtT(new Date(+now + s * 3600e3));
    $("#h-offi").textContent = "UTC" + fmtO(o).replace(".0000 h", "");
    $("#h-offs").textContent = "UTC" + fmtO(s);
    $("#h-tz").textContent = R[cur].iana;
    const d = Math.round((s - o) * 60);
    $("#h-diff").textContent = d === 0 ? "In sync with official time" : `${Math.abs(d)} min ${d < 0 ? "behind" : "ahead of"} official time`;
    $$("[data-off]").forEach(el => el.textContent = fmtO(smooth(now, el.dataset.off)));
    $$("[data-offi]").forEach(el => el.textContent = fmtO(official(now, el.dataset.offi)).replace(".0000", ""));
  }

  function chart() {
    const W = 800, H = 260, P = { l: 56, r: 12, t: 12, b: 28 }, now = new Date(), start = new Date(+now - 30 * 864e5), span = 365 * 864e5, N = 365;
    const pts = [], off = [];
    for (let i = 0; i <= N; i++) { const t = new Date(+start + span * i / N); pts.push(smooth(t, cur)); off.push(official(t, cur)); }
    const all = pts.concat(off), lo = Math.floor(Math.min(...all) * 2) / 2 - 0.25, hi = Math.ceil(Math.max(...all) * 2) / 2 + 0.25;
    const x = i => P.l + (W - P.l - P.r) * i / N, y = v => P.t + (H - P.t - P.b) * (hi - v) / (hi - lo);
    const path = a => a.map((v, i) => (i ? "L" : "M") + x(i).toFixed(1) + " " + y(v).toFixed(1)).join("");
    let g = "";
    for (let v = Math.ceil(lo * 2) / 2; v <= hi; v += 0.5) g += `<line x1="${P.l}" x2="${W - P.r}" y1="${y(v)}" y2="${y(v)}" stroke="#232327"/><text x="${P.l - 8}" y="${y(v) + 4}" fill="#8b8b93" font-size="11" text-anchor="end" font-family="JetBrains Mono,monospace">${v >= 0 ? "+" : ""}${v.toFixed(1)}h</text>`;
    for (let m = 0; m < 13; m++) {
      const d = new Date(Date.UTC(start.getUTCFullYear(), start.getUTCMonth() + m, 1)); const i = (d - start) / span * N;
      if (i < 0 || i > N) continue;
      g += `<text x="${x(i)}" y="${H - 8}" fill="#8b8b93" font-size="11" text-anchor="middle" font-family="Inter,sans-serif">${d.toLocaleString("en", { month: "short", timeZone: "UTC" })}</text>`;
    }
    const ti = 30 / 365 * N;
    g += `<line x1="${x(ti)}" x2="${x(ti)}" y1="${P.t}" y2="${H - P.b}" stroke="#ededef" stroke-dasharray="3 3" opacity=".6"/>`;
    g += `<path d="${path(off)}" fill="none" stroke="#8b8b93" stroke-width="1.5"/><path d="${path(pts)}" fill="none" stroke="#2dd4bf" stroke-width="2.25"/>`;
    $("#chart").innerHTML = g;
    $("#c-region").textContent = R[cur].n;
  }

  function setRegion(k) {
    cur = k; $("#region").value = k; $("#h-addr").textContent = k + ".smoothtime.io";
    $$(".srv").forEach(e => e.textContent = k);
    tick(); chart();
  }

  function copy(text, btn) {
    navigator.clipboard && navigator.clipboard.writeText(text).then(() => {
      btn.classList.add("acc"); const sp = btn.querySelector("span"); if (sp) sp.textContent = "Copied";
      setTimeout(() => { btn.classList.remove("acc"); if (sp) sp.textContent = "Copy"; }, 1500);
    });
  }
  $("#h-copy").onclick = e => copy($("#h-addr").textContent, e.currentTarget);
  $$("[data-copy]").forEach(b => b.onclick = () => copy(b.dataset.copy, b));
  $$("pre").forEach(p => { const b = document.createElement("button"); b.className = "icon-btn"; b.setAttribute("aria-label", "Copy"); b.innerHTML = $("#h-copy svg").outerHTML; b.onclick = () => copy(p.querySelector("code").innerText.split("\n").filter(l => !l.startsWith("#")).join("\n"), b); p.appendChild(b); });
  $$(".tab").forEach(t => t.onclick = () => { $$(".tab").forEach(x => x.setAttribute("aria-selected", x === t)); $$(".pane").forEach(p => p.hidden = p.dataset.p !== t.dataset.t); });
  $("#region").onchange = e => setRegion(e.target.value);


  const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const map = { "America/New_York": "us-eastern", "America/Toronto": "us-eastern", "America/Chicago": "us-central", "America/Denver": "us-mountain", "America/Los_Angeles": "us-pacific",
    "Europe/London": "europe-west", "Europe/Dublin": "europe-west", "Europe/Lisbon": "europe-west", "Europe/Bucharest": "europe-east", "Europe/Helsinki": "europe-east", "Europe/Riga": "europe-east",
    "Europe/Tallinn": "europe-east", "Europe/Vilnius": "europe-east", "Europe/Sofia": "europe-east", "Europe/Athens": "europe-east", "Europe/Kyiv": "europe-east", "Europe/Kiev": "europe-east",
    "Australia/Sydney": "australia-sydney", "Australia/Melbourne": "australia-sydney", "Pacific/Auckland": "new-zealand", "America/Santiago": "chile", "America/Argentina/Buenos_Aires": "argentina" };
  setRegion(map[tz] || "europe-central");
  setInterval(tick, 1000);
})();
