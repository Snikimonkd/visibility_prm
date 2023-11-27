package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"main/graph"
	"main/model"
)

const (
	Width  = 100
	Height = 100
	M      = 1000
)

func goodPoint(obstacles []model.Polygon, p model.Point) bool {
	for i := 0; i < len(obstacles); i++ {
		if obstacles[i].Inside(p) {
			return false
		}
	}

	return true
}

func findPath(
	start, end model.Point,
	obstacles []model.Polygon,
) ([]model.Point, float64, []model.Point, []model.Point) {
	guard := []model.Point{start, end}
	connections := []model.Point{}

	ntry := 0

	for {
		if ntry >= M {
			fmt.Println("ntry reached M -> break")
			break
		}

		path, l := graph.GetPath(guard, connections, obstacles)
		if len(path) != 0 {
			fmt.Println("len of path: ", l, "path: ", path)
			return path, l, guard, connections
		}

		newPoint := model.NewRandPoint(float64(Width), float64(Height))
		for !goodPoint(obstacles, newPoint) {
			newPoint = model.NewRandPoint(float64(Width), float64(Height))
		}

		count := 0
		for i := 0; i < len(guard); i++ {
			seg := model.Segment{A: guard[i], B: newPoint}

			flag := true
			for j := 0; j < len(obstacles); j++ {
				if obstacles[j].IntersectsWithSegment(seg) {
					flag = false
					break
				}
			}
			if flag {
				count += 1
			}
		}

		if count == 0 {
			newPoint.G = true
			guard = append(guard, newPoint)
			ntry = 0
			// fmt.Println("append to guard: ", newPoint)
			continue
		}
		if count == 1 {
			ntry += 1
			// fmt.Println("skip: ", newPoint)
			continue
		}
		if count >= 2 {
			connections = append(connections, newPoint)
			ntry += 1
			// fmt.Println("append to connections: ", newPoint)
			continue
		}
	}

	return nil, 0, guard, connections
}

func gnuplot(
	path []model.Point,
	obstacles []model.Polygon,
	guard []model.Point,
	connections []model.Point,
) {
	out, err := os.Create("./plot/out.plot")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return
	}
	defer out.Close()

	_, err = out.WriteString(fmt.Sprintf(`
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
    `))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	obst, err := os.Create("./plot/obst.plot")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return
	}
	defer obst.Close()

	for i := 0; i < len(obstacles); i++ {
		for j := 0; j < len(obstacles[i].Points); j++ {
			_, err := obst.WriteString(
				fmt.Sprintf("%f %f\n",
					obstacles[i].Points[j].X,
					obstacles[i].Points[j].Y,
				))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		_, err = obst.WriteString(
			fmt.Sprintf("%f %f\n",
				obstacles[i].Points[0].X,
				obstacles[i].Points[0].Y,
			))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err := obst.WriteString("\n")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	g, err := os.Create("./plot/guard_points.plot")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return
	}
	defer g.Close()

	for i := 0; i < len(guard); i++ {
		_, err = g.WriteString(
			fmt.Sprintf("%f %f\n",
				guard[i].X,
				guard[i].Y,
			))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	c, err := os.Create("./plot/conn_points.plot")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return
	}
	defer g.Close()

	for i := 0; i < len(connections); i++ {
		_, err = c.WriteString(
			fmt.Sprintf("%f %f\n",
				connections[i].X,
				connections[i].Y,
			))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	pathLines, err := os.Create("./plot/path_lines.plot")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return
	}
	defer pathLines.Close()

	for i := 0; i < len(path); i++ {
		_, err = pathLines.WriteString(
			fmt.Sprintf("%f %f\n",
				path[i].X,
				path[i].Y,
			))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func main() {
	file, err := os.Open("concave.json")
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		os.Exit(-1)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("can't read file: %v", err.Error())
		os.Exit(-1)
	}

	var res []model.Object
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Printf("can't unmarshal json: %v", err.Error())
		os.Exit(-1)
	}

	var start, end model.Point
	var obstacles []model.Polygon

	for _, obj := range res {
		switch obj.Type {
		case "info":
		case "startPoint":
			start, err = obj.ToPoint()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			start.G = true
		case "endPoint":
			end, err = obj.ToPoint()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			end.G = true
		case "polygon":
			poly, err := obj.ToPolygon()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			obstacles = append(obstacles, poly)
		default:
			fmt.Println("wtf", obj)
			os.Exit(-1)
		}
	}

	path, _, g, c := findPath(start, end, obstacles)
	gnuplot(path, obstacles, g, c)
}
