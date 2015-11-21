$(document).ready(function(){
  //initialize dashboard sidebar slim scroll
  $('#sidebar-dashboard .slimscroll').slimScroll({
    height: '100%',
  });

  //initialize ckeditor
  //need to change filebrowserUploadUrl to '/admin/ckupload' in ckeditor.js
  if ($('#ckeditor').length > 0 ) {
    CKEDITOR.replace('ckeditor', {
      extraPlugins: 'imageuploader',
      language: 'ru',
      allowedContent: true,
      height: 500
    });
  };
});
