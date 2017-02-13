package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/vg"
	uuid "github.com/satori/go.uuid"
)

const (
	inDir  = "/pfs/distributions/"
	outDir = "/pfs/out/"
)

func main() {

	// Walk over files in the input directory.
	if err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {

		// Do nothing if we encounter a directory.
		if info.IsDir() {
			return nil
		}

		// Otherwise, open the file.
		f, err := os.Open(filepath.Join(inDir, info.Name()))
		if err != nil {
			return err
		}

		// Create a 1D histogram value.
		hist := hbook.NewH1D(20, 0, +15)

		// Scan over the lines of the file.
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {

			// Read in the line as text.
			rawVal := scanner.Text()

			// Try to convert the line to a float64.
			val, err := strconv.ParseFloat(rawVal, 64)
			if err != nil {
				return err
			}

			// Add the value to the 1D histogram.
			hist.Fill(val, 1)
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		// normalize histogram
		area := 0.0
		for _, bin := range hist.Binning().Bins() {
			area += bin.SumW() * bin.XWidth()
		}
		hist.Scale(1 / area)

		// Make a plot and set its title.
		p, err := hplot.New()
		if err != nil {
			return err
		}
		p.Title.Text = info.Name()
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"

		// Create a histogram of our values drawn
		// from the standard normal.
		h, err := hplot.NewH1D(hist)
		if err != nil {
			return err
		}
		h.Infos.Style = hplot.HInfoSummary
		p.Add(h)

		// Draw a grid.
		p.Add(hplot.NewGrid())

		// Save the plot to a PNG file.
		id := uuid.NewV4()
		outFile := filepath.Join(outDir, id.String())
		if err := p.Save(6*vg.Inch, -1, outFile+".png"); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
