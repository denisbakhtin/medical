$(document).ready(function(){
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
  if ($('#ckeditor').length > 0 ) {
    CKEDITOR.replace('ckeditor', {
      extraPlugins: 'imageuploader',
      language: 'ru',
      allowedContent: true,
      height: 500
    });
  }

  //captcha-slider
  if ($('#captcha-slider').length > 0 ) {
    var slider = document.getElementById('captcha-slider');
    noUiSlider.create(slider, {
      start: 0,
      range: {
        'min': 0,
        'max': 100
      }
    });
    slider.noUiSlider.on('update', function(){
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

  $('#textarea-comment').on("focusin", function() {
    $('#textarea-comment').css("height", "10em");
  });
  $('#textarea-comment').on("focusout", function() {
    $('#textarea-comment').css("height", "5em");
  });
  $('#textarea-comment').one("keyup", function() {
    $('#comment-hidden').show('fast');
  });

  $(window).on("scroll", function() {
    if ($('#navbar-main').length > 0) {
      var iCurScrollPos = $(this).scrollTop();
      var navbarBottom = $('#navbar-main').position().top+$('#navbar-main').outerHeight(true)
      if (iCurScrollPos > navbarBottom) {
        //show cta-popup
        $('#navbar-cta').css('display', 'block');
      } else {
        //hide cta-popup
        $('#navbar-cta').css('display', 'none');
      }
    }
  });
});
