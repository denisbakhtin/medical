var gulp = require("gulp"),
    sass = require("gulp-sass"),
    autoprefixer = require("gulp-autoprefixer"),
    notify = require("gulp-notify"),
    concat = require("gulp-concat"),
    gzip = require("gulp-gzip"),
    del = require("del")


// Compile SCSS files to CSS
gulp.task("scss", function () {
    //Delete our old css files
    del(["public/css/**/*"])

    //compile hashed css files
    gulp.src("src/scss/**/*.scss")
        .pipe(sass({
            outputStyle: "compressed",
            includePaths: [
                "src/scss",
            ]
        }).on("error", notify.onError(function (error) {
            return "Error: " + error.message;
        })))
        .pipe(autoprefixer({
            browsers: ["last 20 versions"]
        }))
        .pipe(gulp.dest("public/css"))
})

// images
gulp.task("images", function () {
    del(["public/images/**/*"])
    gulp.src("src/images/**/*")
        .pipe(gulp.dest("public/images"))
})

// javascript
gulp.task("js", function () {
    //del(["public/js/**/*"])
    gulp.src("src/js/ckeditor/**/*")
        .pipe(gulp.dest("public/ckeditor"))
    gulp.src([
        "src/js/jquery-2.1.4.min.js",
        "src/js/jquery.actual.min.js",
        "src/js/parsley.min.js",
        "src/js/bootstrap.min.js",
        "src/js/jquery.slimscroll.min.js",
        "src/js/jquery.maskedinput.min.js",
        "src/js/select2.min.js",
        "src/js/nouislider.min.js",
        "src/js/lightbox.min.js",
        "src/js/siema.min.js",
        "src/js/application.js",
    ])
        .pipe(concat("application.js"))
        .pipe(gulp.dest("public/js"))
})

// fonts
gulp.task('icons', function () {
    del(["public/fonts/**/*"])
    gulp.src('src/fonts/**.*')
        .pipe(gulp.dest('public/fonts'));
});

// gzip
gulp.task('gzip', function () {
    gulp.src('public/js/*.js')
        .pipe(gzip())
        .pipe(gulp.dest('public/js'));
    gulp.src('public/css/*.css')
        .pipe(gzip())
        .pipe(gulp.dest('public/css'));
});

// Watch asset folder for changes
gulp.task("watch", ["scss", "images", "js"], function () {
    gulp.watch("src/scss/**/*", ["scss"])
    gulp.watch("src/images/**/*", ["images"])
    gulp.watch("src/js/**/*", ["js"])
})

// Set watch as default task
gulp.task("default", ["icons", "scss", "images", "js", "gzip"])