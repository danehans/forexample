package forexample

import (
	"net/http"
	"os"
	"testing"

	"istio.io/fortio/fhttp"
	"istio.io/fortio/periodic"
)

func TestDestination(t *testing.T) {
	tests := []struct {
		dest string
	}{
		{
			dest: "http://localhost:8080/echo/",
		},
		{
			dest: "http://35.230.102.204/echo/",
		},
		{
			dest: "http://35.230.62.254:8080/echo/",
		},
		{
			dest: "http://35.227.220.111/echo/",
		},
		{
			dest: "http://35.230.50.21/fortio1/echo/",
		},
	}

	for _, tc := range tests {
		// Run 100 qps for exactly 300 requests for 5 threads to the URL destination.
		opts := fhttp.HTTPRunnerOptions{
			RunnerOptions: periodic.RunnerOptions{
				QPS:        100,
				Exactly:    300,
				NumThreads: 5,
				Out:        os.Stderr,
			},
			HTTPOptions: fhttp.HTTPOptions{
				URL: tc.dest,
			},
		}

		// Run http load test.
		res, err := fhttp.RunHTTPTest(&opts)
		if err != nil {
			t.Errorf("Generating traffic via Fortio failed: %v", err)
		}

		totalReqs := res.DurationHistogram.Count
		succReqs := float64(res.RetCodes[http.StatusOK])
		badReqs := res.RetCodes[http.StatusBadRequest]
		actualDuration := res.ActualDuration.Seconds()

		t.Logf("Successfully sent request(s) to destination: %s; calculating results...", tc.dest)
		t.Logf("Fortio Summary: %d reqs (%f rps, %f 200s (%f rps), %d 400s)",
			totalReqs, res.ActualQPS, succReqs, succReqs/actualDuration, badReqs)
	}
}
