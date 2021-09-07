package seeddata

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var randomGenerator *rand.Rand

func init() {
	// setup a random number generator specifically for the birthday generator
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type dataset struct {
	name string
	data []entryRow
}

var datasets struct {
	sync.Mutex
	datasets []dataset
}

func pick(set string) string {

	ds, err := getDataSet(set)
	if err != nil {
		panic(err)
	}

	r := randomGenerator.Float64()

	current := ds.data[0]
	// TODO: linear array search is inefficient
	for _, v := range ds.data {
		current = v
		if r < v.cdf {
			break
		}
	}

	return current.value
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
		return dataset{}, fmt.Errorf("no data found")
	}

	return ds, nil
}

func loadPopulationData() {

	total := 0
	for _, v := range populationDistribution {
		total += v.Count
	}

	var all []entryRow
	cumulative := 0
	for i, v := range populationDistribution {
		cumulative += v.Count
		row := entryRow{
			value:       strconv.Itoa(i),
			probability: float64(v.Count) / float64(total),
			cdf:         float64(cumulative) / float64(total),
		}

		all = append(all, row)
	}

	// rescale the cdf so that the CDF covers (0,1]
	max := all[len(all)-1].cdf
	for i, v := range all {
		v.cdf = v.cdf / max
		all[i] = v
	}

	ds := dataset{
		name: "ages",
		data: all,
	}

	datasets.datasets = append(datasets.datasets, ds)
}

type entryRow struct {
	value       string
	probability float64
	cdf         float64
}

func loadFile(fn string) error {
	// need to load the set in
	file, err := nameFiles.Open("namedata/" + fn)
	if err != nil {
		return err
	}
	defer file.Close()

	rows, err := loadFileContents(file)
	if err != nil {
		return err
	}
	ds := dataset{
		name: fn,
		data: rows,
	}
	datasets.datasets = append(datasets.datasets, ds)
	return nil
}

func loadFileContents(f io.Reader) ([]entryRow, error) {
	r := bufio.NewReader(f)

	all := []entryRow{}
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

		row := entryRow{
			value:       name,
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
