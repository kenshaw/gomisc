package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/knq/imdb"
)

const (
	nfoDesc = `<movie>
  <set>%s</set>
</movie>
http://www.imdb.com/title/%s/
`
)

type info struct {
	Title, ImdbID string
}

var (
	setFlag     = flag.String("set", "", "name of the set to use")
	pathFlag    = flag.String("path", "", "directory path to walk")
	fileExtFlag = flag.String("fileExt", ".avi", "file extension")
)

func main() {
	flag.Parse()

	m := map[string]info{
		"Wabbit Trouble":                   {"Wabbit Twouble", "tt0034368"},
		"The Scarlet Pumpernickle":         {"The Scarlet Pumpernickel", "tt0042928"},
		"Duck Dodgers in the 24.5 Century": {"Duck Dodgers in the 24½th Century", "tt0045709"},
		"The Awful Orphan":                 {"Awful Orphan", "tt0041142"},
		"Puddy Tat Trouble":                {"Putty Tat Trouble", "tt0043943"},
		"Broomstick Bunny":                 {"Broom-Stick Bunny", "tt0049032"},
		"Ready, Set, Zoom":                 {"Ready.. Set.. Zoom!", "tt0048544"},
		"Scrambled Arches":                 {"Scrambled Aches", "tt0050942"},
		"All Abir-r-r-d":                   {"All a Bir-r-r-d", "tt0042191"},
		"Back Alley Uproar":                {"Back Alley Oproar", "tt0040143"},
		"The Three Little Bops":            {"Three Little Bops", "tt0051078"},
		"A Hare Grows in Manhatten":        {"A Hare Grows in Manhattan", "tt0039448"},
		"An Egg Scrambler":                 {"An Egg Scramble", "tt0042431"},
		"Daffy Duck and Egghead":           {"Daffy Duck & Egghead", "tt0030034"},
		"Gonzales' Tomales":                {"Gonzales' Tamales", "tt0050449"},
	}

	re := regexp.MustCompile(`^[0-9]+\s*-\s*`)
	err := filepath.Walk(*pathFlag, func(filepath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filepath, *fileExtFlag) {
			nfoPath := strings.TrimSuffix(filepath, *fileExtFlag) + ".nfo"
			name := strings.TrimSpace(strings.TrimSuffix(path.Base(filepath), *fileExtFlag))
			if re.MatchString(name) {
				name = re.ReplaceAllString(name, "")
			}

			movie := info{"", ""}

			if mi, ok := m[name]; ok {
				movie.Title = mi.Title
				movie.ImdbID = mi.ImdbID
			} else {
				res, err := imdb.MovieByTitle(name, "")
				if err == nil {
					movie.Title = res.Title
					movie.ImdbID = res.ImdbID
				}
			}

			if movie.Title != "" {
				err = ioutil.WriteFile(nfoPath, []byte(fmt.Sprintf(nfoDesc, *setFlag, movie.ImdbID)), 0644)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				fmt.Printf("could not find: %s ( %s )\n", name, nfoPath)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
