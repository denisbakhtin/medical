const { series, parallel, src, dest, watch } = require('gulp');
const sass = require('gulp-sass')(require('sass'));
const notify = require("gulp-notify");
const del = require("del");
const autoprefixer = require("gulp-autoprefixer");
const concat = require("gulp-concat");
const gzip = require("gulp-gzip");

function fonts() {
  del(["public/fonts/**/*"])
  return src('src/fonts/**.*', { encoding: false })
    .pipe(dest('public/fonts'));
}

function scss() {
  return src("src/scss/**/*.scss")
    .pipe(sass({
      outputStyle: "compressed",
      includePaths: [
        "src/scss",
      ]
    }).on("error", notify.onError(function(error) {
      return "Error: " + error.message;
    })))
    .pipe(autoprefixer())
    .pipe(dest("public/css"))
}

function images() {
  del(["public/images/**/*"])
  return src("src/images/**/*", { encoding: false })
    .pipe(dest("public/images"))
}

function ckeditor() {
  del(["public/js/ckeditor/*"])
  return src("src/js/ckeditor/*", { encoding: false })
    .pipe(dest("public/js/ckeditor"))
}

function js() {
  //del(["public/js/**/*"])
  return src([
    "src/js/jquery-2.1.4.min.js",
    "src/js/jquery.actual.min.js",
    "src/js/parsley.min.js",
    "src/js/bootstrap.min.js",
    "src/js/jquery.slimscroll.min.js",
    "src/js/jquery.maskedinput.min.js",
    "src/js/nouislider.min.js",
    "src/js/lightbox.min.js",
    "src/js/siema.min.js",
    "src/js/application.js"])
    .pipe(concat("application.js"))
    .pipe(dest("public/js"))
}

function gzipJs() {
  return src('public/js/*.js')
    .pipe(gzip())
    .pipe(dest('public/js'));
}

function gzipCss() {
  return src('public/css/*.css')
    .pipe(gzip())
    .pipe(dest('public/css'));
}

exports.watch = function() {
  watch(["src/scss/**/*.scss"], scss);
  watch(["src/js/**/*.js"], js);
  watch(["src/images/**/*"], images);
};

exports.default = series(parallel(fonts, scss, images, js, ckeditor), parallel(gzipJs, gzipCss));
