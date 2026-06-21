package benchmark

import "sort"

func SortResults(
	results []Result,
) {

	sort.Slice(results,
		func(i, j int) bool {

			return results[i].Latency <
				results[j].Latency
		},
	)
}
