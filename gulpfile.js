var gulp = require('gulp');
var gulpif = require('gulp-if');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var minifyCSS = require('gulp-minify-css');
var uncss = require('gulp-uncss');

var debug = false;
gulp.task('scripts', function() {
    gulp.src(['./js/zepto.js', './js/fx.js', './js/event.js','./js/zepto-slide-transition.js', './js/main.js'])
    .pipe(concat('all.js', {newLine: ';'}))
    .pipe(gulpif(!debug,uglify()))
    .pipe(gulp.dest('./static/'));
});

gulp.task('templates', function() {
    gulp.src('templates/layout.jade')
    .pipe(replace('@@timestamp', (new Date()).getTime()))
    .pipe(rename('layout2.jade'))
    .pipe(gulp.dest('templates/'));
});

gulp.task('css', function() {
    gulp.src(['./css/bootstrap.min.css', './css/main.css'])
    .pipe(concat('all.css'))
    .pipe(uncss({
        html: ['./templates/head.html', './templates/index.html', './templates/about.html'],
        ignore: [/btn-sel/, /btn-nonsel/, /has-error/]
     }))
    .pipe(gulpif(!debug, minifyCSS()))
    .pipe(gulp.dest('./static/'));
});

gulp.task('watch', function() {
  gulp.watch('js/*', ['scripts']);
  gulp.watch('css/*', ['css']);
});


gulp.task('default', ['css', 'scripts', /*'watch', 'templates'*/]);
