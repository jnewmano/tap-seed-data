package seeddata

import (
	"bufio"
	"embed"
	"io"
	"strconv"
	"strings"
	"sync"
)

const (
	lastNameDataset        = "lastnames.txt"
	femaleFirstNameDataset = "firstnames-female.txt"
	maleFirstNameDataset   = "firstnames-male.txt"
)

//go:embed namedata
var nameFiles embed.FS

type dataset struct {
	name string
	data []nameFileRow
}

var datasets struct {
	sync.Mutex
	datasets []dataset
}

func RandomName() (firstName string, lastName string) {
	if randomGenerator.Float64() < 0.5 {
		return RandomMaleName()
	}
	return RandomFemaleName()
}

func RandomFemaleName() (firstName string, lastName string) {
	firstName = randomFemaleFirstName()
	lastName = randomLastName()
	return
}

func RandomMaleName() (firstName string, lastName string) {
	firstName = randomMaleFirstName()
	lastName = randomLastName()
	return
}

func randomMaleFirstName() string {
	return pick(maleFirstNameDataset)
}

func randomFemaleFirstName() string {
	return pick(femaleFirstNameDataset)
}

func randomLastName() string {
	return pick(lastNameDataset)
}

func pick(set string) string {

	ds, err := getDataSet(set)
	if err != nil {
		panic(err)
	}
	r := randomGenerator.Float64()

	last := ds.data[0]
	// TODO: linear array search is inefficient
	for _, v := range ds.data {
		if v.cdf >= r {
			break
		}
		last = v
	}

	return last.name
}

func getDataSet(set string) (dataset, error) {
	datasets.Lock()
	defer datasets.Unlock()

	// find a matching dataset
	var ds dataset
	for _, v := range datasets.datasets {
		if v.name == set {
			ds = v
			break
		}
	}

	if len(ds.data) == 0 {
		// need to load the set in
		file, err := nameFiles.Open("namedata/" + set)
		if err != nil {
			return dataset{}, err
		}
		defer file.Close()

		rows, err := loadFileContents(file)
		if err != nil {
			return dataset{}, err
		}
		ds = dataset{
			name: set,
			data: rows,
		}
		datasets.datasets = append(datasets.datasets, ds)
	}
	return ds, nil
}

type nameFileRow struct {
	name        string
	probability float64
	cdf         float64
}

func loadFileContents(f io.Reader) ([]nameFileRow, error) {
	r := bufio.NewReader(f)

	all := []nameFileRow{}
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// split the line
		name := strings.Title(strings.ToLower(strings.TrimSpace(string(l[0:15]))))
		rateString := strings.TrimSpace(string(l[15:21]))
		rate, err := strconv.ParseFloat(rateString, 64)
		if err != nil {
			return nil, err
		}

		cdfString := strings.TrimSpace(string(l[21:28]))
		cdf, err := strconv.ParseFloat(cdfString, 64)
		if err != nil {
			return nil, err
		}

		row := nameFileRow{
			name:        name,
			probability: rate / 100,
			cdf:         cdf / 100,
		}

		all = append(all, row)
	}

	// rescale the cdf so that the CDF covers (0,1]
	// our name list doesn't include all names, we
	// only have a truncated dataset
	max := all[len(all)-1].cdf
	for i, v := range all {
		v.cdf = v.cdf / max
		all[i] = v
	}

	// TODO: use a data structure that isn't a slice
	// to make name lookups faster
	return all, nil
}
