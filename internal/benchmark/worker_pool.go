package benchmark

import "sync"

func RunWorkers(
	servers []string,
	workers int,
) []Result {

	jobs := make(chan string)

	results := make(chan Result)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for ip := range jobs {

				results <- TestDNS(ip)

			}

		}()
	}

	go func() {

		for _, s := range servers {

			jobs <- s

		}

		close(jobs)

	}()

	go func() {

		wg.Wait()

		close(results)

	}()

	var out []Result

	for r := range results {

		out = append(out, r)

	}

	return out
}
