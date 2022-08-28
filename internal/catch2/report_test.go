package catch2_test

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lesomnus/catch2-compare/internal/catch2"
)

const data = `
<?xml version="1.0" encoding="UTF-8"?>
<Catch2TestRun name="test_build" rng-seed="3413175810" catch2-version="3.1.0">
  <TestCase name="Build Point Cloud" filename="/workspaces/hesai/tests/src/build.cpp" line="18">
    <BenchmarkResults name="Thor Implementation" samples="100" resamples="100000" iterations="1" clockResolution="16.7063" estimatedDuration="2.0186e+07">
      <!-- All values in nano seconds -->
      <mean value="162750" lowerBound="149438" upperBound="187959" ci="0.95"/>
      <standardDeviation value="90465.3" lowerBound="54039.8" upperBound="143714" ci="0.95"/>
      <outliers variance="0.989504" lowMild="0" lowSevere="0" highMild="1" highSevere="8"/>
    </BenchmarkResults>
    <BenchmarkResults name="Hesai SDK" samples="100" resamples="100000" iterations="1" clockResolution="16.7063" estimatedDuration="1.924e+08">
      <!-- All values in nano seconds -->
      <mean value="1.85689e+06" lowerBound="1.83899e+06" upperBound="1.8798e+06" ci="0.95"/>
      <standardDeviation value="103265" lowerBound="84804.2" upperBound="129199" ci="0.95"/>
      <outliers variance="0.534717" lowMild="0" lowSevere="0" highMild="8" highSevere="2"/>
    </BenchmarkResults>
    <OverallResult success="true"/>
  </TestCase>
  <OverallResults successes="0" failures="0" expectedFailures="0"/>
  <OverallResultsCases successes="1" failures="0" expectedFailures="0"/>
</Catch2TestRun>
`

func TestUnmarshal(t *testing.T) {
	require := require.New(t)
	expected := catch2.XmlReport{
		Name:    "test_build",
		Version: "3.1.0",
		TestCases: []catch2.TestCase{{
			Name:     "Build Point Cloud",
			Filename: "/workspaces/hesai/tests/src/build.cpp",
			Line:     18,
			BenchmarkResults: []catch2.BenchmarkResult{
				{
					Name:     "Thor Implementation",
					Duration: catch2.Duration{2.0186e+07},
					Mean: catch2.Measurement{
						Value:      catch2.Duration{162750},
						LowerBound: catch2.Duration{149438},
						UpperBound: catch2.Duration{187959},
					},
				},
				{
					Name:     "Hesai SDK",
					Duration: catch2.Duration{1.924e+08},
					Mean: catch2.Measurement{
						Value:      catch2.Duration{1.85689e+06},
						LowerBound: catch2.Duration{1.83899e+06},
						UpperBound: catch2.Duration{1.8798e+06},
					},
				},
			},
		}},
	}

	var actual catch2.XmlReport

	err := xml.Unmarshal([]byte(data), &actual)
	require.NoError(err)

	require.Equal(expected, actual)
}
