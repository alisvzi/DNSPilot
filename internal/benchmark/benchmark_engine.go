package benchmark

func FindFastestDNS(
	servers []string,
) (*Result, []Result) {

	results := RunWorkers(
		servers,
		50,
	)

	SortResults(results)

	if len(results) == 0 {

		return nil, nil
	}

	return &results[0], results
}
