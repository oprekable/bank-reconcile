package service

import (
	"github.com/oprekable/bank-reconcile/internal/app/service/process"
	"github.com/oprekable/bank-reconcile/internal/app/service/sample"
)

type Services struct {
	SvcSample  sample.Service
	SvcProcess process.Service
}
