package err

import (
	"github.com/lyraproj/issue/issue"
	"runtime"
)

func NoArgs(code issue.Code) issue.Reported {
	_, file, line, _ := runtime.Caller(1)
	return issue.NewReported(code, issue.SeverityError, issue.NoArgs, issue.NewLocation(file, line, 0))
}

func WithArgs(code issue.Code, args issue.H) issue.Reported {
	_, file, line, _ := runtime.Caller(1)
	return issue.NewReported(code, issue.SeverityError, args, issue.NewLocation(file, line, 0))
}
