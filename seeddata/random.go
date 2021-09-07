package seeddata

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
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
	// find a matching dataset
	for _, v := range datasets.datasets {
		if v.name == set {
			return v, nil
		}
	}

	return dataset{}, fmt.Errorf("no data found")
}

type entryRow struct {
	value string
	cdf   float64
}

func addDataset(name string, data []entryRow) {

	// TODO: use a data structure that isn't a slice
	// to make name lookups faster

	// normalize the cdf so that the CDF covers (0,1]
	max := data[len(data)-1].cdf
	for i, v := range data {
		v.cdf = v.cdf / max
		data[i] = v
	}

	ds := dataset{
		name: name,
		data: data,
	}

	datasets.datasets = append(datasets.datasets, ds)
}

func loadPopulationData() {

	var all []entryRow
	cumulative := 0
	for i, v := range populationDistribution {
		cumulative += v.Count
		row := entryRow{
			value: strconv.Itoa(i),
			cdf:   float64(cumulative), // addDataset will normalize the range of the CDF to (0,1]
		}

		all = append(all, row)
	}

	addDataset("ages", all)
}

func loadNameFile(fn string) error {

	file, err := nameFiles.Open("namedata/" + fn)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)

	all := []entryRow{}
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// split the line
		name := strings.Title(strings.ToLower(strings.TrimSpace(string(l[0:15]))))

		cdfString := strings.TrimSpace(string(l[21:28]))
		cdf, err := strconv.ParseFloat(cdfString, 64)
		if err != nil {
			return err
		}

		row := entryRow{
			value: name,
			cdf:   cdf,
		}

		all = append(all, row)
	}

	addDataset(fn, all)

	return nil
}
