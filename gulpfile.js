const { series, parallel, src, dest, watch } = require("gulp");
const sass = require("gulp-sass")(require("sass"));
const notify = require("gulp-notify");
const del = require("del");
const autoprefixer = require("gulp-autoprefixer");
const concat = require("gulp-concat");
const gzip = require("gulp-gzip");

function scss() {
  // return src("src/scss/**/*.scss")
  return src(["src/scss/application.scss"])
    .pipe(
      sass({
        outputStyle: "compressed",
        includePaths: ["src/scss", "node_modules"],
        quietDeps: true, //to stop bootstrap color deprecated warnings
      }).on(
        "error",
        notify.onError(function (error) {
          return "Error: " + error.message;
        }),
      ),
    )
    .pipe(autoprefixer())
    .pipe(dest("public/css"));
}

function images() {
  del(["public/images/**/*"]);
  return src("src/images/**/*", { encoding: false }).pipe(
    dest("public/images"),
  );
}

function ckeditor() {
  del(["public/js/ckeditor/*"]);
  return src("src/js/ckeditor/*", { encoding: false }).pipe(
    dest("public/js/ckeditor"),
  );
}

function js() {
  //del(["public/js/**/*"])
  return src(["src/js/siema.min.js", "src/js/application.js"])
    .pipe(concat("application.js"))
    .pipe(dest("public/js"));
}

function gzipJs() {
  return src("public/js/*.js").pipe(gzip()).pipe(dest("public/js"));
}

function gzipCss() {
  return src("public/css/*.css").pipe(gzip()).pipe(dest("public/css"));
}

exports.watch = function () {
  watch(["src/scss/**/*.scss"], scss);
  watch(["src/js/**/*.js"], js);
  watch(["src/images/**/*"], images);
};

exports.default = series(
  parallel(scss, images, js, ckeditor),
  parallel(gzipJs, gzipCss),
);
