package root

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func Runner(cmd *cobra.Command, _ []string) (er error) {
	return cmd.Help()
}

func PersistentPreRunner(_ *cobra.Command, _ []string) (er error) {
	fromDate, errFrom := time.Parse(DateFormatString, FlagFromDateValue)

	if errFrom != nil {
		return fmt.Errorf("failed to parse flag '-%s' '--%s': %v", FlagFromDateShort, FlagFromDate, FlagFromDateValue)
	}

	toDate, errTo := time.Parse(DateFormatString, FlagToDateValue)
	if errTo != nil {
		return fmt.Errorf("failed to parse flag '-%s' '--%s': %v", FlagToDateShort, FlagToDate, FlagToDateValue)
	}

	if fromDate.After(toDate) {
		return fmt.Errorf("'-%s' '--%s': %v should before '-%s' '--%s': %v", FlagFromDateShort, FlagFromDate, FlagFromDateValue, FlagToDateShort, FlagToDate, FlagToDateValue)
	}

	if FlagTotalDataSampleToGenerateValue <= 0 {
		return fmt.Errorf("'-%s' '--%s': %v should bigger than 0", FlagTotalDataSampleToGenerateShort, FlagTotalDataSampleToGenerate, FlagTotalDataSampleToGenerateValue)
	}

	if FlagPercentageMatchSampleToGenerateValue < 0 || FlagPercentageMatchSampleToGenerateValue > 100 {
		return fmt.Errorf("'-%s' '--%s': %v should between 0 and 100", FlagPercentageMatchSampleToGenerateShort, FlagPercentageMatchSampleToGenerate, FlagPercentageMatchSampleToGenerateValue)
	}

	return nil
}
