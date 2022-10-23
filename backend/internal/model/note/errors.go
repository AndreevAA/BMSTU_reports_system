package report

type CanNotCreateReportErr struct{}

func (a *CanNotCreateReportErr) Error() string {
	return "can't create report"
}

type CanNotAssignReportErr struct{}

func (a *CanNotAssignReportErr) Error() string {
	return "can't assign report"
}

type ReportNotFoundErr struct{}

func (a *ReportNotFoundErr) Error() string {
	return "report does not exist or does not belong to user"
}
