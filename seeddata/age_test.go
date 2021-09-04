package seeddata

import (
	"math"
	"testing"

	"github.com/bearbin/go-age"
)

func TestAgeDistribution(t *testing.T) {

	threshold := 0.001
	total := 1000000

	buckets := make([]int, len(populationDistribution))
	for i := 0; i < total; i++ {
		birthday := RandomBirthday()

		age := age.Age(birthday)

		idx := age / 5
		buckets[idx] = buckets[idx] + 1
	}

	distSum := 0
	for _, v := range populationDistribution {
		distSum += v.Count
	}

	for i, v := range buckets {
		part := float64(v) / float64(total)
		distPart := float64(populationDistribution[i].Count) / float64(distSum)

		if math.Abs(part-distPart) > threshold {
			t.Fatalf("distribution might be wrong for bucket [%d] got [%0.3f]%% instead of [%0.3f]%%\n", i, part*100, distPart*100)
		}
	}
}
