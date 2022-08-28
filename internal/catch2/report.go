package catch2

type Measurement struct {
	Value      Duration `xml:"value,attr"`
	LowerBound Duration `xml:"lowerBound,attr"`
	UpperBound Duration `xml:"upperBound,attr"`
}

type BenchmarkResult struct {
	Name     string      `xml:"name,attr"`
	Duration Duration    `xml:"estimatedDuration,attr"`
	Mean     Measurement `xml:"mean"`
}

type TestCase struct {
	Name     string `xml:"name,attr"`
	Filename string `xml:"filename,attr"`
	Line     int    `xml:"line,attr"`

	BenchmarkResults []BenchmarkResult `xml:"BenchmarkResults"`
}

type XmlReport struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"catch2-version,attr"`

	TestCases []TestCase `xml:"TestCase"`
}
