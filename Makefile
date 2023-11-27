run:
	go run main.go
	gnuplot ./plot/out.plot
	open out.png
