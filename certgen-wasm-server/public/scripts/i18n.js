// Client-side i18n for static Astro pages.
// Supported languages: en (default), es, pt.
(function () {
  const translations = {
    en: {
      "lang.label": "Language",
      "hero.title": "CertGen (offline) - generate your CA + server cert",
      "hero.hint.before1": "Fill",
      "hero.hint.before2": "and",
      "hero.hint.after1": "Then download a",
      "hero.hint.after2": "with",
      "form.kicker": "Form",
      "form.title": "Certificate data",
      "form.subtitle": "Validation runs in browser and in WASM backend to avoid generation errors.",
      "field.domain": "Domain",
      "field.ip": "IP",
      "field.org": "Org",
      "field.domain.help": "DNS name used by the browser (must match exactly).",
      "field.ip.help": "Server IP address (added to SAN).",
      "field.org.help": "Organization name written inside the certificate.",
      "button.generate": "Generate & Download ZIP",
      "side.title": "Which file goes where?",
      "side.subtitle": "These files are generated inside",
      "file.servercrt": "Server certificate (public part).",
      "file.serverkey": "Server private key (secret).",
      "file.cacrt": "CA that signed the server cert (chain).",
      "file.cakey": "CA private key (VERY secret, do not upload).",
      "file.fullchain": "Server cert + CA cert in one file (useful for some systems).",
      "validity.title": "Certificate validity",
      "validity.ca": "CA certificate (ca.crt): 10 years",
      "validity.server": "Server certificate (server.crt): 5 years",
      "validity.note": "These validity periods are enforced by the WASM backend during generation.",
      "vital.title": "Generic file mapping example",
      "mapping.note":
        "Use this mapping in platforms that ask for Certificate / Key / Chain fields (Asterisk, FreePBX, IssabelPBX, VitalPBX).",
      "asterisk.title": "Asterisk example (general)",
      "asterisk.hint":
        "Asterisk usually reads cert/key from configured paths (for example WSS/WebRTC). In general use fullchain.crt or server.crt + ca.crt, and server.key.",
      "alert.checkFields": "Please fix the fields highlighted in red.",
      "alert.generating": "Generating certificates in WASM...",
      "alert.validationFailed": "Backend validation failed. Please review the fields.",
      "alert.zipError": "Could not generate the ZIP.",
      "alert.downloaded": "Done. Downloaded {filename}",
      "alert.swNotReady": "Service Worker/WASM not ready. Open http://localhost:8000 and reload once.",
      "err.invalidDomain": "Invalid domain (example: pbx.local)",
      "err.invalidIP": "Invalid IP (example: 192.168.1.10)",
      "err.orgRequired": "Org is required",
      "err.orgTooLong": "Org is too long (max 128 chars)"
    },
    es: {
      "lang.label": "Idioma",
      "hero.title": "CertGen (offline) - genera tu CA + certificado de servidor",
      "hero.hint.before1": "Completa",
      "hero.hint.before2": "y",
      "hero.hint.after1": "Luego descargas un",
      "hero.hint.after2": "con",
      "form.kicker": "Formulario",
      "form.title": "Datos del certificado",
      "form.subtitle": "Validamos en el navegador y en el backend WASM para evitar errores de generación.",
      "field.domain": "Dominio",
      "field.ip": "IP",
      "field.org": "Org",
      "field.domain.help": "Nombre DNS usado por el navegador (debe coincidir exactamente).",
      "field.ip.help": "IP del servidor (se agrega al SAN).",
      "field.org.help": "Nombre de la organización dentro del certificado.",
      "button.generate": "Generar y descargar ZIP",
      "side.title": "¿Qué archivo va dónde?",
      "side.subtitle": "Estos archivos se generan dentro de",
      "file.servercrt": "Certificado del servidor (parte pública).",
      "file.serverkey": "Clave privada del servidor (secreta).",
      "file.cacrt": "CA que firmó el certificado del servidor (cadena).",
      "file.cakey": "Clave privada de la CA (MUY secreta, no subir).",
      "file.fullchain": "Certificado del servidor + CA en un solo archivo (útil en algunos sistemas).",
      "validity.title": "Vigencia de certificados",
      "validity.ca": "Certificado CA (ca.crt): 10 años",
      "validity.server": "Certificado de servidor (server.crt): 5 años",
      "validity.note": "Estos periodos de vigencia se aplican y validan en el backend WASM durante la generación.",
      "vital.title": "Ejemplo genérico de correspondencia",
      "mapping.note":
        "Usa esta correspondencia en plataformas que piden campos Certificate / Key / Chain (Asterisk, FreePBX, IssabelPBX, VitalPBX).",
      "asterisk.title": "Ejemplo Asterisk (general)",
      "asterisk.hint":
        "Asterisk suele leer certificado/clave desde rutas configuradas (por ejemplo WSS/WebRTC). En general usa fullchain.crt o server.crt + ca.crt, y server.key.",
      "alert.checkFields": "Revisa los campos marcados en rojo.",
      "alert.generating": "Generando certificados en WASM...",
      "alert.validationFailed": "Falló la validación del backend. Corrige los campos.",
      "alert.zipError": "No se pudo generar el ZIP.",
      "alert.downloaded": "Listo. Se descargó {filename}",
      "alert.swNotReady": "Service Worker/WASM no listo. Abre http://localhost:8000 y recarga una vez.",
      "err.invalidDomain": "Dominio inválido (ejemplo: pbx.local)",
      "err.invalidIP": "IP inválida (ejemplo: 192.168.1.10)",
      "err.orgRequired": "Org es requerido",
      "err.orgTooLong": "Org es demasiado largo (máx. 128 caracteres)"
    },
    pt: {
      "lang.label": "Idioma",
      "hero.title": "CertGen (offline) - gere sua CA + certificado do servidor",
      "hero.hint.before1": "Preencha",
      "hero.hint.before2": "e",
      "hero.hint.after1": "Depois baixe um",
      "hero.hint.after2": "com",
      "form.kicker": "Formulário",
      "form.title": "Dados do certificado",
      "form.subtitle": "A validação roda no navegador e no backend WASM para evitar erros de geração.",
      "field.domain": "Domínio",
      "field.ip": "IP",
      "field.org": "Org",
      "field.domain.help": "Nome DNS usado no navegador (deve corresponder exatamente).",
      "field.ip.help": "IP do servidor (adicionado ao SAN).",
      "field.org.help": "Nome da organização dentro do certificado.",
      "button.generate": "Gerar e baixar ZIP",
      "side.title": "Qual arquivo vai onde?",
      "side.subtitle": "Esses arquivos são gerados dentro de",
      "file.servercrt": "Certificado do servidor (parte pública).",
      "file.serverkey": "Chave privada do servidor (secreta).",
      "file.cacrt": "CA que assinou o certificado do servidor (cadeia).",
      "file.cakey": "Chave privada da CA (MUITO secreta, não enviar).",
      "file.fullchain": "Certificado do servidor + CA em um arquivo (útil em alguns sistemas).",
      "validity.title": "Validade dos certificados",
      "validity.ca": "Certificado CA (ca.crt): 10 anos",
      "validity.server": "Certificado do servidor (server.crt): 5 anos",
      "validity.note": "Esses períodos de validade são aplicados e validados no backend WASM durante a geração.",
      "vital.title": "Exemplo genérico de mapeamento",
      "mapping.note":
        "Use este mapeamento em plataformas que pedem campos Certificate / Key / Chain (Asterisk, FreePBX, IssabelPBX, VitalPBX).",
      "asterisk.title": "Exemplo Asterisk (geral)",
      "asterisk.hint":
        "O Asterisk normalmente lê certificado/chave de caminhos configurados (por exemplo WSS/WebRTC). Em geral use fullchain.crt ou server.crt + ca.crt, e server.key.",
      "alert.checkFields": "Revise os campos marcados em vermelho.",
      "alert.generating": "Gerando certificados em WASM...",
      "alert.validationFailed": "A validação do backend falhou. Corrija os campos.",
      "alert.zipError": "Não foi possível gerar o ZIP.",
      "alert.downloaded": "Pronto. Download de {filename}",
      "alert.swNotReady": "Service Worker/WASM não pronto. Abra http://localhost:8000 e recarregue uma vez.",
      "err.invalidDomain": "Domínio inválido (exemplo: pbx.local)",
      "err.invalidIP": "IP inválido (exemplo: 192.168.1.10)",
      "err.orgRequired": "Org é obrigatório",
      "err.orgTooLong": "Org é muito longo (máx. 128 caracteres)"
    }
  };

  const DEFAULT_LANG = "en";

  function interpolate(text, vars) {
    if (!vars) return text;
    return text.replace(/\{(\w+)\}/g, (_, k) => (k in vars ? String(vars[k]) : `{${k}}`));
  }

  function getLang() {
    const saved = localStorage.getItem("app_lang");
    if (saved && translations[saved]) return saved;
    return DEFAULT_LANG;
  }

  function t(key, vars) {
    const lang = getLang();
    const text = translations[lang]?.[key] ?? translations[DEFAULT_LANG]?.[key] ?? key;
    return interpolate(text, vars);
  }

  function applyI18n() {
    const lang = getLang();
    document.documentElement.lang = lang;

    document.querySelectorAll("[data-i18n]").forEach((el) => {
      const key = el.getAttribute("data-i18n");
      if (!key) return;
      el.textContent = t(key);
    });

    document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
      const key = el.getAttribute("data-i18n-placeholder");
      if (!key) return;
      el.setAttribute("placeholder", t(key));
    });

    const select = document.getElementById("langSelect");
    if (select && select.value !== lang) select.value = lang;
  }

  function setLang(lang) {
    if (!translations[lang]) return;
    localStorage.setItem("app_lang", lang);
    applyI18n();
    window.dispatchEvent(new CustomEvent("app:lang-changed", { detail: { lang } }));
  }

  function initLanguageSelector() {
    const select = document.getElementById("langSelect");
    if (!select) return;
    select.addEventListener("change", (event) => {
      setLang(event.target.value);
    });
  }

  window.I18N = { t, getLang, setLang, apply: applyI18n };

  window.addEventListener("DOMContentLoaded", () => {
    applyI18n();
    initLanguageSelector();
  });
})();

