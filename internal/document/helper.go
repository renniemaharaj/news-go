package document

func ReportFromQuery(s string) Report {
	return Report{
		SearchQuery: s,
	}
}
