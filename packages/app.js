let fs = require("fs");
let browserify = require("browserify");
let moduleName = "Buffer"
browserify({ ignoreMissing: true, standalone: moduleName })
    .transform(
        "babelify", {
            presets: ["babel-preset-es2015"],
            plugins:["babel-plugin-transform-remove-console"],
            compact:true
        }
    )
    .require(require.resolve(moduleName),{ entry: true })
    .bundle()
    .pipe(fs.createWriteStream(`${moduleName}.js`));