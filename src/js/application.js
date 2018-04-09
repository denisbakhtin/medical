$(document).ready(function () {
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
  //need to change filebrowserUploadUrl to '/admin/ckupload' in ckeditor.js
  if ($('#ckeditor').length > 0) {
    CKEDITOR.replace('ckeditor', {
      extraPlugins: 'imageuploader',
      language: 'ru',
      allowedContent: true,
      height: 500
    });
  }

  //captcha-slider
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

  /*
  $('#textarea-comment').on("focusin", function() {
    $('#textarea-comment').css("height", "10em");
  });
  $('#textarea-comment').on("focusout", function() {
    $('#textarea-comment').css("height", "5em");
  });
  */
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

  var withoutpainheight = function () {
    var maxh = 0;
    $('#withoutpain-slide .item .thumbnail').each(function () {
      var item = $(this);
      if (item.actual('height') > maxh) {
        maxh = item.actual('height');
      }
    });
    $('#withoutpain-slide .item .thumbnail').each(function () {
      var item = $(this);
      item.height(maxh);
    });
  }
  withoutpainheight();

  $('#withoutpain-slide .item').each(function () {
    var itemToClone = $(this);
    for (var i = 1; i < 3; i++) {
      itemToClone = itemToClone.next();

      if (!itemToClone.length) {
        break;
        //itemToClone = $(this).siblings(':first');
      }

      itemToClone.children(':first-child').clone()
        .addClass("cloneditem-" + (i))
        .appendTo($(this));
    }
  });

  $('#testimonials-slide .item').each(function () {
    var itemToClone = $(this);
    for (var i = 1; i < 2; i++) {
      itemToClone = itemToClone.next();

      if (!itemToClone.length) {
        itemToClone = $(this).siblings(':first');
      }

      itemToClone.children(':first-child').clone()
        .addClass("cloneditem-" + (i))
        .appendTo($(this));
    }
  });

  $('#phone-input').mask("9 (999) 999-9999");

  $(window).resize(function () {
    withoutpainheight();
  });

});