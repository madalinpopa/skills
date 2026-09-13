(function () {
  var root = document.documentElement;

  function storedTheme() {
    try {
      return localStorage.getItem("theme");
    } catch (e) {
      return null;
    }
  }

  function applyTheme(theme) {
    if (theme) {
      root.setAttribute("data-theme", theme);
    } else {
      root.removeAttribute("data-theme");
    }
  }

  applyTheme(storedTheme());

  function currentTheme() {
    var explicit = root.getAttribute("data-theme");
    if (explicit) {
      return explicit;
    }
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }

  document.addEventListener("click", function (event) {
    var toggle = event.target.closest("[data-theme-toggle]");
    if (toggle) {
      var next = currentTheme() === "dark" ? "light" : "dark";
      applyTheme(next);
      try {
        localStorage.setItem("theme", next);
      } catch (e) {}
      return;
    }

    var copy = event.target.closest("[data-copy]");
    if (copy) {
      var source = copy.getAttribute("data-copy");
      var text = source
        ? document.getElementById(source).innerText
        : copy.parentElement.querySelector("pre").innerText;
      text = text
        .split("\n")
        .map(function (line) {
          return line.replace(/^\$ /, "");
        })
        .join("\n")
        .trim();
      navigator.clipboard.writeText(text).then(function () {
        var label = copy.textContent;
        copy.textContent = "copied";
        setTimeout(function () {
          copy.textContent = label;
        }, 1200);
      });
      return;
    }

    var sidebarToggle = event.target.closest("[data-sidebar-toggle]");
    if (sidebarToggle) {
      var sidebar = document.querySelector(".sidebar");
      sidebar.hidden = !sidebar.hidden;
      sidebarToggle.setAttribute("aria-expanded", String(!sidebar.hidden));
    }
  });

  document.querySelectorAll(".code").forEach(function (block) {
    if (block.querySelector(".copy")) {
      return;
    }
    var button = document.createElement("button");
    button.className = "copy";
    button.type = "button";
    button.textContent = "copy";
    button.setAttribute("data-copy", "");
    button.setAttribute("aria-label", "Copy to clipboard");
    block.appendChild(button);
  });

  var sidebarLinks = document.querySelectorAll(".sidebar a[href^='#']");
  if (sidebarLinks.length && "IntersectionObserver" in window) {
    var byId = {};
    sidebarLinks.forEach(function (link) {
      byId[link.getAttribute("href").slice(1)] = link;
    });
    var visible = new Set();
    var observer = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            visible.add(entry.target.id);
          } else {
            visible.delete(entry.target.id);
          }
        });
        var first = null;
        document.querySelectorAll(".content [id]").forEach(function (el) {
          if (!first && visible.has(el.id) && byId[el.id]) {
            first = el.id;
          }
        });
        if (first) {
          sidebarLinks.forEach(function (link) {
            link.classList.toggle("active", link.getAttribute("href") === "#" + first);
          });
        }
      },
      { rootMargin: "-80px 0px -60% 0px", threshold: 0 }
    );
    Object.keys(byId).forEach(function (id) {
      var el = document.getElementById(id);
      if (el) {
        observer.observe(el);
      }
    });
  }

  var term = document.querySelector("[data-terminal]");
  if (term) {
    var script = JSON.parse(term.getAttribute("data-terminal"));
    var reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    var body = term.querySelector(".term-body");
    body.innerHTML = "";

    function esc(s) {
      return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
    }

    function renderStatic() {
      body.innerHTML = script
        .map(function (step) {
          return '<span class="p">$ </span><span class="c">' + esc(step.cmd) + "</span>\n" + step.out;
        })
        .join("\n");
    }

    if (reduce) {
      renderStatic();
      return;
    }

    var typed = "";
    var index = 0;

    function wait(ms) {
      return new Promise(function (resolve) {
        setTimeout(resolve, ms);
      });
    }

    function paint(cmdSoFar, showCaret) {
      body.innerHTML =
        typed + '<span class="p">$ </span><span class="c">' + esc(cmdSoFar) + "</span>" + (showCaret ? '<span class="caret"></span>' : "");
    }

    (async function () {
      paint("", true);
      while (true) {
        await wait(index === 0 ? 500 : 800);
        var step = script[index];
        for (var i = 1; i <= step.cmd.length; i++) {
          paint(step.cmd.slice(0, i), true);
          await wait(26 + Math.random() * 40);
        }
        await wait(260);
        typed += '<span class="p">$ </span><span class="c">' + esc(step.cmd) + "</span>\n" + step.out + "\n\n";
        paint("", true);
        index += 1;
        if (index >= script.length) {
          await wait(7000);
          typed = "";
          index = 0;
          paint("", true);
        }
      }
    })();
  }
})();
