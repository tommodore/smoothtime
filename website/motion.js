(() => {
  const d = document, h = d.querySelector("header"), pr = d.querySelector(".progress");
  const sc = () => { h && h.classList.toggle("scrolled", scrollY > 8); if (pr) { const m = d.documentElement.scrollHeight - innerHeight; pr.style.transform = `scaleX(${m > 0 ? scrollY / m : 0})`; } };
  addEventListener("scroll", sc, { passive: true }); sc();
  const u = d.getElementById("utc"); if (u) { const t = () => u.textContent = new Date().toISOString().slice(11, 19); t(); setInterval(t, 1000); }
  if (matchMedia("(prefers-reduced-motion: reduce)").matches || !("IntersectionObserver" in window)) return;
  const io = new IntersectionObserver(es => es.forEach(e => { if (e.isIntersecting) { e.target.classList.add("in"); e.target.closest("section")?.classList.add("seen"); io.unobserve(e.target); } }), { threshold: .12, rootMargin: "0px 0px -6% 0px" });
  d.querySelectorAll("main section:not(.hero) .wrap").forEach(w => [...w.children].forEach((c, i) => { c.classList.add("rv"); c.style.setProperty("--i", Math.min(i, 6)); }));
  d.querySelectorAll("#why .card, #who .card, details").forEach(c => { c.classList.add("rv"); c.style.setProperty("--i", [...c.parentNode.children].indexOf(c) % 6); });
  const ms = () => d.querySelectorAll(".eyebrow.rv").forEach(l => { const c = l.cloneNode(true); c.className = "eyebrow"; c.style.cssText = "position:absolute;visibility:hidden;width:auto"; l.parentNode.appendChild(c); l.style.setProperty("--w", c.getBoundingClientRect().width + 14 + "px"); l.style.setProperty("--n", c.textContent.length + 12); c.remove(); });
  ms(); d.fonts && d.fonts.ready.then(ms);
  d.documentElement.classList.add("m");
  d.querySelectorAll(".rv").forEach(el => io.observe(el));
  const ch = d.getElementById("chart");
  if (ch) {
    let seen = false;
    const run = () => { ch.querySelectorAll("path").forEach(p => p.style.setProperty("--l", Math.ceil(p.getTotalLength()))); ch.classList.remove("draw"); void ch.getBBox(); ch.classList.add("draw"); };
    new IntersectionObserver((es, o) => es.forEach(e => { if (e.isIntersecting) { seen = true; run(); o.disconnect(); } }), { threshold: .3 }).observe(ch);
    new MutationObserver(() => seen && run()).observe(ch, { childList: true });
  }
})();
