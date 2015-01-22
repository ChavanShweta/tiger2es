package main

import (
	"archive/zip"
	shp "github.com/jonas-p/go-shp"
	geojson "github.com/kpawlik/geojson"
	"io"
	"log"
	"os"
	//"reflect"
)

func UnzipFile(zipFile string) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer rc.Close()

		f, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func ReadShapefile(filename string) {
	getLines(filename)

}

func getLines(filename string) {
	values := make(map[string]interface{})
	file, err := shp.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for file.Next() {
		_, shape := file.Shape()
		line := shape.(*shp.PolyLine)
		geometry := getLineString(line)
		feature := geojson.NewFeature(geometry, values, nil)
		log.Print(geojson.Marshal(feature))

	}
}

func getLineString(line *shp.PolyLine) *geojson.LineString {
	points := line.Points
	coordinates := geojson.Coordinates{}
	for _, point := range points {
		c := geojson.Coordinate{geojson.Coord(point.X), geojson.Coord(point.Y)}
		coordinates = append(coordinates, c)
	}
	linestring := geojson.NewLineString(coordinates)
	return linestring
}