
    set terminal png
    set size square
    set xrange [0:100]
    set yrange [0:100]
    set grid
    set output 'out.png'
    set palette model RGB defined (0 "black",1 "green")
    plot './plot/obst.plot' using 1:2 with polygons notitle fillstyle transparent solid 0.5, \
         './plot/guard_points.plot' using 1:2 with points notitle linewidth 1 linecolor rgb 'black' pointtype 7, \
         './plot/conn_points.plot' using 1:2 with points notitle linewidth 1 linecolor rgb 'black' pointtype 6, \
         './plot/path_lines.plot' using 1:2 with lines notitle
    