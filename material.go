package rayngo

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type Material struct {
	Diffuse Color
	Specularity uint16
	Textured bool
}

func NewMaterialFromMtl(path string) *Material {
	// For now, assume one material definition per .mtl file. This makes parsing way easy.
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var mat Material
	mat.Textured = false

	lineScanner := bufio.NewScanner(fi)

	for lineScanner.Scan() {
		line := lineScanner.Text()

		var tokens []string
		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			tokens = append(tokens, wordScanner.Text())
		}

		switch tokens[0] {
		// case "newmtl":
		// 	do nothing. Assuming one material per file, so we don't need to track names.
		// case "Ka":
		// 	do nothing. Don't support ambient color yet.
		case "Kd":
			r, _ := strconv.ParseFloat(tokens[1], 64)
			g, _ := strconv.ParseFloat(tokens[2], 64)
			b, _ := strconv.ParseFloat(tokens[3], 64)
			mat.Diffuse = Color{r, g, b, 1.0}
		// case "Ks":
		// 	do nothing. Don't support specular color yet.
		case "Ns":
			spec, _ := strconv.ParseFloat(tokens[1], 64)
			mat.Specularity = uint16(spec)
		// case "d", "Tr":
		// 	No transparency support yet.
		// case "illum":
		// 	No support for mtl illumination models
		}
	}

	return &mat
}

func (m Material) Sample(u, v float64) Color {
	if m.Textured {
		if math.Mod((math.Ceil(u) + math.Ceil(v)), 2) == 0 {
			return Color{0.6, 0.6, 0.5, 1.0}
		}

		return Color{0.75, 0.75, 0.75, 1.0}
	}

	return m.Diffuse
}