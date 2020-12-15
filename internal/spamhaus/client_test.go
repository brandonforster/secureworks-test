package spamhaus

import "testing"

func TestLookup(t *testing.T) {
	var testCases = []struct {
		ip            string
		expectedError bool
		codes         []string
	}{
		{"175.214.225.175", false, []string{"127.0.0.4", "127.0.0.2"}},
		{"46.173.215.108", false, []string{"127.0.0.2", "127.0.0.9"}},
		{"google.com", true, nil},
		{"127.0.0.1", false, nil},
	}

	sut := Client{}

	for _, tt := range testCases {
		testname := tt.ip
		t.Run(testname, func(t *testing.T) {
			actual, err := sut.Lookup(tt.ip)
			if err != nil && !tt.expectedError {
				t.Errorf("unexpected error %s", err)
			} else if err == nil && tt.expectedError {
				t.Errorf("expected error")
			}

			if len(actual) != len(tt.codes) {
				t.Errorf("unexpected number of results, expected %d was %d", len(tt.codes), len(actual))
			}

			expectedMap := make(map[string]int)
			for _, code := range tt.codes {
				expectedMap[code]++
			}

			actualMap := make(map[string]int)
			for _, code := range actual {
				actualMap[code]++
			}

			for expectedKey, expectedVal := range expectedMap {
				if actualMap[expectedKey] != expectedVal {
					t.Errorf("different results than expected for code %s", expectedKey)
				}
			}
		})
	}
}
