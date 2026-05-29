document.addEventListener("DOMContentLoaded", function () {
  //siema slider
  const withoutpain = document.getElementById("withoutpain-slider");
  if (withoutpain) {
    const withoutPainSiema = new Siema({
      selector: "#withoutpain-slider",
      duration: 200,
      easing: "ease-out",
      perPage: {
        100: 1,
        576: 2,
        992: 3,
      },
      startIndex: 0,
      draggable: true,
      multipleDrag: true,
      threshold: 20,
      loop: false,
      rtl: false,
      onInit: () => {},
      onChange: () => {},
    });
    document
      .querySelector(".withoutpain-prev")
      .addEventListener("click", () => withoutPainSiema.prev());
    document
      .querySelector(".withoutpain-next")
      .addEventListener("click", () => withoutPainSiema.next());
    //setInterval(() => withoutPainSiema.next(), 4000);
  }

  const testimonials = document.getElementById("testimonials-slider");
  if (testimonials) {
    const testimonialsSiema = new Siema({
      selector: "#testimonials-slider",
      duration: 200,
      easing: "ease-out",
      perPage: {
        100: 1,
        850: 2,
      },
      startIndex: 0,
      draggable: true,
      multipleDrag: true,
      threshold: 20,
      loop: true,
      rtl: false,
      onInit: () => {},
      onChange: () => {
        console.log("on change siema");
      },
    });
    document
      .querySelector(".testimonials-prev")
      .addEventListener("click", () => testimonialsSiema.prev());
    document
      .querySelector(".testimonials-next")
      .addEventListener("click", () => testimonialsSiema.next());
    //setInterval(() => withoutPainSiema.next(), 4000);
  }

  //captcha-slider - optin form
  const slider = document.getElementById("range-captcha");
  if (slider) {
    const captcha = document.getElementById("captcha");
    const btn = document.getElementById("submit-with-captcha-btn");

    slider.oninput = function () {
      if (this.value == 100) {
        captcha.value = btoa(this.value);
        btn.disabled = false;
      } else {
        btn.disabled = true;
      }
    };
  }

  //comment captcha slider
  const cslider = document.getElementById("range-comment-captcha");
  if (cslider) {
    const captcha = document.getElementById("captcha-comment");
    const btn = document.getElementById("submit-with-captcha-comment-btn");

    cslider.oninput = function () {
      if (this.value == 100) {
        captcha.value = btoa(this.value);
        btn.disabled = false;
      } else {
        btn.disabled = true;
      }
    };
  }

  const areaComment = document.getElementById("textarea-comment");
  if (areaComment) {
    areaComment.addEventListener("keyup", (event) => {
      document.getElementById("comment-hidden").style.display = "block";
    });
  }

  // $("#comments-more-btn").on("click", function () {
  //   $(".comment.extra.hide").first().removeClass("hide");
  //   $(".comment.extra.hide").first().removeClass("hide");
  //   $("#comments-less-btn").removeClass("hide");
  // });
  // $("#comments-less-btn").on("click", function () {
  //   $(".comment.extra").addClass("hide");
  //   $("#comments-less-btn").addClass("hide");
  // });

  window.addEventListener("scroll", (event) => {
    const iCurScrollPos = document.documentElement.scrollTop;
    const navbar = document.getElementById("navbar-main");
    const navbarTop = navbar.getBoundingClientRect().top;
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const navbarBottom = scrollTop + navbarTop + navbar.offsetHeight;
    const navbarScroll = document.getElementById("navbar-scroll");
    if (iCurScrollPos > navbarBottom) {
      navbarScroll.style.display = "block";
    } else {
      navbarScroll.style.display = "none";
    }
  });

  // $("#phone-input").mask("9 (999) 999-9999");
  // >>>> vanilla masked phone input: +7 (999) 999-99-99
  const phoneInput = document.getElementById("phone-input");
  if (phoneInput) {
    phoneInput.addEventListener("input", function (e) {
      let value = this.value;

      // Удаляем все нецифровые символы
      let digits = value.replace(/\D/g, "");

      // Если первая цифра 8, заменяем на 7 для формата +7
      if (digits.startsWith("8") && digits.length > 1) {
        digits = "7" + digits.substring(1);
      }

      let formattedValue = "";

      if (digits) {
        formattedValue = "+";

        // Добавляем первую цифру (7)
        if (digits.length >= 1) {
          formattedValue += digits[0];
        }

        // Добавляем скобку и следующие 3 цифры
        if (digits.length >= 2) {
          formattedValue += " (";
          formattedValue += digits.substring(1, 4);
        }

        // Закрываем скобку и добавляем следующие 3 цифры с дефисом
        if (digits.length >= 5) {
          formattedValue += ") ";
          formattedValue += digits.substring(4, 7);
        }

        // Добавляем дефис и следующие 2 цифры
        if (digits.length >= 8) {
          formattedValue += "-";
          formattedValue += digits.substring(7, 9);
        }

        // Добавляем последний дефис и 2 цифры
        if (digits.length >= 10) {
          formattedValue += "-";
          formattedValue += digits.substring(9, 11);
        }
      }

      this.value = formattedValue;
    });

    // Блокировка ввода букв и спецсимволов
    phoneInput.addEventListener("keydown", function (e) {
      const allowedKeys = [
        "Backspace",
        "Delete",
        "ArrowLeft",
        "ArrowRight",
        "Tab",
        "Enter",
        "Escape",
      ];

      // Разрешаем служебные клавиши
      if (allowedKeys.includes(e.key)) {
        return;
      }

      // Блокируем всё, что не является цифрой
      if (!/\d/.test(e.key) && e.key !== "8") {
        e.preventDefault();
      }
    });
  }
  // <<< masked phone input

  // filter infos
  document.filterInfos = function (published) {
    published = String(published);
    const filters = document.querySelectorAll("#filters > li");
    filters?.forEach((el) => {
      el.classList.remove("active");
    });
    const hidden = document.querySelectorAll(".table > tbody > tr.hidden");
    hidden?.forEach((el) => {
      el.classList.remove("hidden");
    });
    document.setCookie("infoFilter", published, 365);
    if (published == "true") {
      const pub = document.getElementById("filter-published");
      pub?.classList.add("active");
      const unpubs = document.querySelectorAll(
        ".table > tbody > tr[data-published=false]",
      );
      unpubs?.forEach((el) => {
        el.classList.add("hidden");
      });
    } else if (published == "false") {
      const unpub = document.getElementById("filter-unpublished");
      unpub?.classList.add("active");
      const pubs = document.querySelectorAll(
        ".table > tbody > tr[data-published=true]",
      );
      pubs?.forEach((el) => {
        el.classList.add("hidden");
      });
    } else {
      document.getElementById("filter-all")?.classList.add("active");
    }
  };

  //working with cookies
  document.setCookie = function (cname, cvalue, exdays) {
    const d = new Date();
    d.setTime(d.getTime() + exdays * 24 * 60 * 60 * 1000);
    let expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
  };
  document.getCookie = function (cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(";");
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == " ") {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  };
});
