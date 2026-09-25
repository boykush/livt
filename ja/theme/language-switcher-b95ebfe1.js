// Each translation is its own book, built into {lang}/ beneath the English one
// (the docs:build task), so the same page in another language is the same path
// under another root.
(() => {
  const labels = { en: "English", ja: "日本語" };
  const buttons = document.querySelector("#mdbook-menu-bar .right-buttons");
  const lang = document.documentElement.lang;
  if (!buttons || !(lang in labels)) return;

  const root = new URL(path_to_root, location.href);
  const englishRoot = lang === "en" ? root : new URL("../", root);
  const page = location.pathname.slice(root.pathname.length);

  for (const [code, label] of Object.entries(labels)) {
    if (code === lang) continue;
    const link = document.createElement("a");
    link.className = "icon-button";
    link.href = new URL(code === "en" ? page : `${code}/${page}`, englishRoot).href;
    link.hreflang = code;
    link.lang = code;
    link.textContent = label;
    buttons.prepend(link);
  }
})();
