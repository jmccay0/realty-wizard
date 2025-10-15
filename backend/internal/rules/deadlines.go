package rules

import (
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/google/uuid"
)

// DeadlineEngine calculates all deadlines based on contract terms
type DeadlineEngine struct {
	holidays map[string]bool // Texas-specific holidays as "YYYY-MM-DD" strings
}

// NewDeadlineEngine creates a new deadline calculation engine
func NewDeadlineEngine() *DeadlineEngine {
	currentYear := time.Now().Year()
	holidayMap := make(map[string]bool)

	// Add holidays for current year, previous year, and next year
	for year := currentYear - 1; year <= currentYear + 1; year++ {
		holidays := getTexasHolidays(year)
		for _, h := range holidays {
			key := h.Format("2006-01-02")
			holidayMap[key] = true
		}
	}

	return &DeadlineEngine{
		holidays: holidayMap,
	}
}

// CalculateDeadlines generates all deadlines from contract terms
func (e *DeadlineEngine) CalculateDeadlines(contract *models.ContractTerms, property *models.Property) []*models.Deadline {
	var deadlines []*models.Deadline

	// Option fee & earnest money due within 3 business days
	optionEarnestDue := e.addBusinessDays(contract.EffectiveDate, 3)
	deadlines = append(deadlines, &models.Deadline{
		ID:            uuid.New().String(),
		ProjectID:     contract.ProjectID,
		Type:          "option_fee",
		Description:   "Option fee due to title/escrow",
		DueDate:       optionEarnestDue,
		IsBusinessDay: true,
		Priority:      "critical",
		Completed:     false,
	})

	deadlines = append(deadlines, &models.Deadline{
		ID:            uuid.New().String(),
		ProjectID:     contract.ProjectID,
		Type:          "earnest_money",
		Description:   "Earnest money due to escrow",
		DueDate:       optionEarnestDue,
		IsBusinessDay: true,
		Priority:      "critical",
		Completed:     false,
	})

	// Option period end (calendar days from effective date)
	if contract.OptionPeriodDays > 0 {
		optionEnd := contract.EffectiveDate.AddDate(0, 0, contract.OptionPeriodDays)
		deadlines = append(deadlines, &models.Deadline{
			ID:            uuid.New().String(),
			ProjectID:     contract.ProjectID,
			Type:          "inspection_period",
			Description:   "End of option/inspection period",
			DueDate:       optionEnd,
			IsBusinessDay: false,
			Priority:      "high",
			Completed:     false,
			Notes:         "Buyer can terminate for any reason before this date",
		})

		// Remind about repair requests 2 days before option ends
		repairRequestReminder := optionEnd.AddDate(0, 0, -2)
		if repairRequestReminder.After(contract.EffectiveDate) {
			deadlines = append(deadlines, &models.Deadline{
				ID:            uuid.New().String(),
				ProjectID:     contract.ProjectID,
				Type:          "repair_request_reminder",
				Description:   "Expect buyer's repair requests",
				DueDate:       repairRequestReminder,
				IsBusinessDay: false,
				Priority:      "normal",
				Completed:     false,
				Notes:         "Buyers typically submit repair requests 24-48hrs before option period ends",
			})
		}
	}

	// Title commitment due (default 20 days if not specified)
	commitmentDays := contract.TitleCommitmentDays
	if commitmentDays == 0 {
		commitmentDays = 20
	}
	titleCommitmentDue := contract.EffectiveDate.AddDate(0, 0, commitmentDays)
	deadlines = append(deadlines, &models.Deadline{
		ID:            uuid.New().String(),
		ProjectID:     contract.ProjectID,
		Type:          "title_commitment",
		Description:   "Title commitment expected from title company",
		DueDate:       titleCommitmentDue,
		IsBusinessDay: false,
		Priority:      "high",
		Completed:     false,
		Notes:         "Review for exceptions and requirements",
	})

	// Survey/T-47 deadline (if property has survey and needs T-47)
	if property != nil && property.HasSurvey {
		// T-47 typically due a few days before closing
		t47Due := contract.ClosingDate.AddDate(0, 0, -7)
		if t47Due.After(contract.EffectiveDate) {
			deadlines = append(deadlines, &models.Deadline{
				ID:            uuid.New().String(),
				ProjectID:     contract.ProjectID,
				Type:          "t47_affidavit",
				Description:   "T-47 Residential Real Property Affidavit due",
				DueDate:       t47Due,
				IsBusinessDay: false,
				Priority:      "high",
				Completed:     false,
				Notes:         "Complete and deliver to title company",
			})
		}
	}

	// Financing milestones (if buyer is financing)
	if contract.BuyerFinancing {
		// Loan approval typically due halfway to closing or 21 days
		loanApprovalDays := (int(contract.ClosingDate.Sub(contract.EffectiveDate).Hours()) / 24) / 2
		if loanApprovalDays > 21 {
			loanApprovalDays = 21
		}
		loanApprovalDue := contract.EffectiveDate.AddDate(0, 0, loanApprovalDays)

		deadlines = append(deadlines, &models.Deadline{
			ID:            uuid.New().String(),
			ProjectID:     contract.ProjectID,
			Type:          "loan_approval",
			Description:   "Buyer's loan approval expected",
			DueDate:       loanApprovalDue,
			IsBusinessDay: false,
			Priority:      "high",
			Completed:     false,
			Notes:         "Confirm buyer's financing is on track",
		})

		// Appraisal typically 10-14 days after effective date
		appraisalDue := contract.EffectiveDate.AddDate(0, 0, 12)
		deadlines = append(deadlines, &models.Deadline{
			ID:            uuid.New().String(),
			ProjectID:     contract.ProjectID,
			Type:          "appraisal",
			Description:   "Property appraisal expected",
			DueDate:       appraisalDue,
			IsBusinessDay: false,
			Priority:      "high",
			Completed:     false,
			Notes:         "Appraisal must meet or exceed sales price",
		})
	}

	// Closing prep reminder (5 days before)
	closingPrepReminder := contract.ClosingDate.AddDate(0, 0, -5)
	if closingPrepReminder.After(contract.EffectiveDate) {
		deadlines = append(deadlines, &models.Deadline{
			ID:            uuid.New().String(),
			ProjectID:     contract.ProjectID,
			Type:          "closing_prep",
			Description:   "Begin closing preparation",
			DueDate:       closingPrepReminder,
			IsBusinessDay: false,
			Priority:      "normal",
			Completed:     false,
			Notes:         "Coordinate utilities, keys, garage openers, HOA docs, etc.",
		})
	}

	// Closing date
	deadlines = append(deadlines, &models.Deadline{
		ID:            uuid.New().String(),
		ProjectID:     contract.ProjectID,
		Type:          "closing",
		Description:   "Closing date",
		DueDate:       contract.ClosingDate,
		IsBusinessDay: false,
		Priority:      "critical",
		Completed:     false,
		Notes:         "Final walkthrough typically morning of closing",
	})

	// Seller's temporary lease end (if applicable)
	if contract.SellerStaysPostClose && contract.SellerLeaseDays > 0 {
		leaseEnd := contract.ClosingDate.AddDate(0, 0, contract.SellerLeaseDays)
		deadlines = append(deadlines, &models.Deadline{
			ID:            uuid.New().String(),
			ProjectID:     contract.ProjectID,
			Type:          "lease_end",
			Description:   "Seller's temporary lease ends - move out date",
			DueDate:       leaseEnd,
			IsBusinessDay: false,
			Priority:      "critical",
			Completed:     false,
			Notes:         "Final move-out and property delivery to buyer",
		})
	}

	return deadlines
}

// addBusinessDays adds business days (Mon-Fri, excluding holidays) to a date
func (e *DeadlineEngine) addBusinessDays(start time.Time, days int) time.Time {
	current := start
	added := 0

	for added < days {
		current = current.AddDate(0, 0, 1)
		if e.isBusinessDay(current) {
			added++
		}
	}

	return current
}

// isBusinessDay checks if a date is a business day (not weekend, not holiday)
func (e *DeadlineEngine) isBusinessDay(date time.Time) bool {
	// Check weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return false
	}

	// Check holidays
	key := date.Format("2006-01-02")
	if e.holidays[key] {
		return false
	}

	return true
}

// getTexasHolidays returns Texas state holidays for a given year
func getTexasHolidays(year int) []time.Time {
	return []time.Time{
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year's Day
		// MLK Day (3rd Monday in January)
		nthWeekday(year, time.January, time.Monday, 3),
		// Presidents Day (3rd Monday in February)
		nthWeekday(year, time.February, time.Monday, 3),
		// Memorial Day (last Monday in May)
		lastWeekday(year, time.May, time.Monday),
		time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
		// Labor Day (1st Monday in September)
		nthWeekday(year, time.September, time.Monday, 1),
		// Thanksgiving (4th Thursday in November)
		nthWeekday(year, time.November, time.Thursday, 4),
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
	}
}

// nthWeekday finds the nth occurrence of a weekday in a month
func nthWeekday(year int, month time.Month, weekday time.Weekday, n int) time.Time {
	first := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Find first occurrence of the target weekday
	daysUntilTarget := int(weekday - first.Weekday())
	if daysUntilTarget < 0 {
		daysUntilTarget += 7
	}

	firstOccurrence := first.AddDate(0, 0, daysUntilTarget)
	return firstOccurrence.AddDate(0, 0, 7*(n-1))
}

// lastWeekday finds the last occurrence of a weekday in a month
func lastWeekday(year int, month time.Month, weekday time.Weekday) time.Time {
	// Start from the last day of the month
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	// Walk backwards to find the weekday
	for lastDay.Weekday() != weekday {
		lastDay = lastDay.AddDate(0, 0, -1)
	}

	return lastDay
}
