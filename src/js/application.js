$(document).ready(function () {
  //siema slider
  if ($('#withoutpain-slider').length > 0) {
    const withoutPainSiema = new Siema({
      selector: '#withoutpain-slider',
      duration: 200,
      easing: 'ease-out',
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
    document.querySelector('.withoutpain-prev').addEventListener('click', () => withoutPainSiema.prev());
    document.querySelector('.withoutpain-next').addEventListener('click', () => withoutPainSiema.next());
    //setInterval(() => withoutPainSiema.next(), 4000);
  }

  if ($('#testimonials-slider').length > 0) {
    const testimonialsSiema = new Siema({
      selector: '#testimonials-slider',
      duration: 200,
      easing: 'ease-out',
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
      onChange: () => {},
    });
    document.querySelector('.testimonials-prev').addEventListener('click', () => testimonialsSiema.prev());
    document.querySelector('.testimonials-next').addEventListener('click', () => testimonialsSiema.next());
    //setInterval(() => withoutPainSiema.next(), 4000);
  }

  //wow animation init
  if (typeof WOW != 'undefined') {
    new WOW().init();
  }

  //initialize dashboard sidebar slim scroll
  if ($('#sidebar-dashboard .slimscroll').length > 0) {
    $('#sidebar-dashboard .slimscroll').slimScroll({
      height: '100%',
    });
  }

  //initialize ckeditor
  //using cdn ckeditor5, see views/dashboard/ckeditor partial
  //need to change filebrowserUploadUrl to '/admin/ckupload' in ckeditor.js
  // if ($('#ckeditor').length > 0) {
  //   CKEDITOR.replace('ckeditor', {
  //     extraPlugins: 'imageuploader',
  //     language: 'ru',
  //     allowedContent: true,
  //     height: 500
  //   });
  // }

  //captcha-slider - optin form
  if ($('#captcha-slider').length > 0) {
    var slider = document.getElementById('captcha-slider');
    noUiSlider.create(slider, {
      start: 0,
      range: {
        'min': 0,
        'max': 100
      }
    });
    slider.noUiSlider.on('update', function () {
      var value = slider.noUiSlider.get();
      if (value == 100) {
        $('#captcha-slider .noUi-handle').addClass('ok');
        $('#captcha-wrapper #ok-sign').addClass('visible');
        $('#captcha').val(btoa(value));
        $('#submit-with-captcha-btn').attr('disabled', null);
      } else {
        $('#captcha-slider .noUi-handle').removeClass('ok');
        $('#captcha-wrapper #ok-sign').removeClass('visible');
        $('#captcha').val(btoa(value));
        $('#submit-with-captcha-btn').attr('disabled', 'disabled');
      }
    });
  }

  //captcha-slider - comments form
  if ($('#captcha-comments-slider').length > 0) {
    var sliderComments = document.getElementById('captcha-comments-slider');
    noUiSlider.create(sliderComments, {
      start: 0,
      range: {
        'min': 0,
        'max': 100
      }
    });
    sliderComments.noUiSlider.on('update', function () {
      var value = sliderComments.noUiSlider.get();
      if (value == 100) {
        $('#captcha-comments-slider .noUi-handle').addClass('ok');
        $('#captcha-comments-wrapper .ok-sign').addClass('visible');
        $('#captcha-comments').val(btoa(value));
        $('#submit-with-captcha-comments-btn').attr('disabled', null);
      } else {
        $('#captcha-comments-slider .noUi-handle').removeClass('ok');
        $('#captcha-comments-wrapper .ok-sign').removeClass('visible');
        $('#captcha-comments').val(btoa(value));
        $('#submit-with-captcha-comments-btn').attr('disabled', 'disabled');
      }
    });
  }


  $('#textarea-comment').one("keyup", function () {
    $('#comment-hidden').show('fast');
  });
  $('#comments-more-btn').on('click', function () {
    $('.comment.extra.hide').first().removeClass('hide');
    $('.comment.extra.hide').first().removeClass('hide');
    $('#comments-less-btn').removeClass('hide');
  });
  $('#comments-less-btn').on('click', function () {
    $('.comment.extra').addClass('hide');
    $('#comments-less-btn').addClass('hide');
  });

  $(window).on("scroll", function () {
    if ($('#navbar-main').length > 0) {
      var iCurScrollPos = $(this).scrollTop();
      var navbarBottom = $('#navbar-main').position().top + $('#navbar-main').outerHeight(true)
      if (iCurScrollPos > navbarBottom) {
        //show cta-popup
        $('#navbar-scroll').css('display', 'block');
      } else {
        //hide cta-popup
        $('#navbar-scroll').css('display', 'none');
      }
    }
  });

  $('#phone-input').mask("9 (999) 999-9999");

  //working with cookies
  document.setCookie = function (cname, cvalue, exdays) {
    const d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    let expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
  }
  document.getCookie = function(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
  }
});
