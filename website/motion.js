(() => {
  const d = document, h = d.querySelector("header");
  if (h) { const f = () => h.classList.toggle("scrolled", scrollY > 8); addEventListener("scroll", f, { passive: true }); f(); }
  if (matchMedia("(prefers-reduced-motion: reduce)").matches || !("IntersectionObserver" in window)) return;
  const io = new IntersectionObserver(es => es.forEach(e => { if (e.isIntersecting) { e.target.classList.add("in"); io.unobserve(e.target); if (e.target.dataset.n) count(e.target); } }), { threshold: .12, rootMargin: "0px 0px -6% 0px" });
  // stagger groups: direct children of each section wrap (except hero h1), plus card grids
  d.querySelectorAll("section .wrap").forEach(w => [...w.children].forEach((c, i) => { if (c.tagName === "H1") return; c.classList.add("rv"); c.style.setProperty("--i", Math.min(i, 6)); }));
  d.querySelectorAll("#why .card, #who .card, .faq details, details").forEach((c, i, a) => { c.classList.add("rv"); c.style.setProperty("--i", [...c.parentNode.children].indexOf(c) % 6); });
  // headline word reveal (keeps <br>)
  const h1 = d.querySelector("h1");
  if (h1) { let i = 0; h1.querySelectorAll("*").length <= 2 && [...h1.childNodes].forEach(n => { if (n.nodeType === 3) { const s = d.createElement("span"); s.innerHTML = n.textContent.split(/(\s+)/).map(t => t.trim() ? `<span class="w" style="--i:${i++}">${t}</span>` : t).join(""); n.replaceWith(...s.childNodes); } }); }
  // count-ups for plain numbers marked data-count, e.g. <b data-count>12</b>
  const count = el => { const end = +el.dataset.n, t0 = performance.now(), dur = 1200; const step = t => { const p = Math.min(1, (t - t0) / dur); el.textContent = Math.round(end * (1 - Math.pow(1 - p, 3))).toLocaleString("en"); if (p < 1) requestAnimationFrame(step); }; requestAnimationFrame(step); };
  d.querySelectorAll("[data-count]").forEach(el => { const n = parseInt(el.textContent.replace(/\D/g, ""), 10); if (n) { el.dataset.n = n; io.observe(el); } });
  d.documentElement.classList.add("m");
  d.querySelectorAll(".rv").forEach(el => io.observe(el));
  // chart draw-in (first time it's visible, and after region switches)
  const ch = d.getElementById("chart");
  if (ch) {
    let seen = false;
    const run = () => { ch.querySelectorAll("path").forEach(p => p.style.setProperty("--l", Math.ceil(p.getTotalLength()))); ch.classList.remove("draw"); void ch.getBBox(); ch.classList.add("draw"); };
    new IntersectionObserver((es, o) => es.forEach(e => { if (e.isIntersecting) { seen = true; run(); o.disconnect(); } }), { threshold: .3 }).observe(ch);
    new MutationObserver(() => seen && run()).observe(ch, { childList: true });
  }
})();
