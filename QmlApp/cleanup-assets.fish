#!/usr/bin/fish

for f in assets/*.svg
    inkscape --export-filename=$f.out.svg -l $f
    mv -v $f.out.svg $f
end

