package rules

import (
	"testing"
	"time"

	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAddBusinessDays(t *testing.T) {
	engine := NewDeadlineEngine()

	tests := []struct {
		name     string
		start    time.Time
		days     int
		expected time.Time
	}{
		{
			name:     "3 business days from Monday",
			start:    time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
			days:     3,
			expected: time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC), // Thursday
		},
		{
			name:     "3 business days from Friday (crosses weekend)",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Friday
			days:     3,
			expected: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), // Wednesday
		},
		{
			name:     "1 business day from Friday",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC), // Friday
			days:     1,
			expected: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), // Monday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.addBusinessDays(tt.start, tt.days)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBusinessDay(t *testing.T) {
	engine := NewDeadlineEngine()

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{
			name:     "Monday is business day",
			date:     time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "Saturday is not business day",
			date:     time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Sunday is not business day",
			date:     time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "New Year's Day is not business day",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "Christmas is not business day",
			date:     time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.isBusinessDay(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateDeadlines(t *testing.T) {
	engine := NewDeadlineEngine()

	effectiveDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday
	closingDate := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)

	contract := &models.ContractTerms{
		ProjectID:           "test-project-1",
		EffectiveDate:       effectiveDate,
		SalesPrice:          450000,
		ClosingDate:         closingDate,
		OptionFee:           500,
		OptionPeriodDays:    10,
		EarnestMoney:        5000,
		TitleCommitmentDays: 20,
		BuyerFinancing:      true,
	}

	property := &models.Property{
		ProjectID: "test-project-1",
		YearBuilt: 2015,
		HasSurvey: true,
	}

	deadlines := engine.CalculateDeadlines(contract, property)

	// Verify we have the expected deadlines
	assert.NotEmpty(t, deadlines, "Should generate deadlines")

	// Count deadline types
	deadlineTypes := make(map[string]int)
	for _, d := range deadlines {
		deadlineTypes[d.Type]++
	}

	// Critical deadlines that should always exist
	assert.Equal(t, 1, deadlineTypes["option_fee"], "Should have option fee deadline")
	assert.Equal(t, 1, deadlineTypes["earnest_money"], "Should have earnest money deadline")
	assert.Equal(t, 1, deadlineTypes["inspection_period"], "Should have inspection period deadline")
	assert.Equal(t, 1, deadlineTypes["title_commitment"], "Should have title commitment deadline")
	assert.Equal(t, 1, deadlineTypes["closing"], "Should have closing deadline")

	// Financing-specific deadlines (only when buyer is financing)
	assert.Equal(t, 1, deadlineTypes["loan_approval"], "Should have loan approval deadline")
	assert.Equal(t, 1, deadlineTypes["appraisal"], "Should have appraisal deadline")

	// Survey-specific deadline
	assert.Equal(t, 1, deadlineTypes["t47_affidavit"], "Should have T-47 deadline when property has survey")

	// Verify option/earnest money is 3 business days from effective date
	for _, d := range deadlines {
		if d.Type == "option_fee" {
			expectedDue := engine.addBusinessDays(effectiveDate, 3)
			assert.Equal(t, expectedDue, d.DueDate, "Option fee should be due 3 business days after effective date")
			assert.True(t, d.IsBusinessDay, "Option fee deadline should be marked as business day calculation")
			assert.Equal(t, "critical", d.Priority, "Option fee should be critical priority")
		}
	}
}

func TestCalculateDeadlinesWithoutFinancing(t *testing.T) {
	engine := NewDeadlineEngine()

	contract := &models.ContractTerms{
		ProjectID:           "test-project-2",
		EffectiveDate:       time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		ClosingDate:         time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		BuyerFinancing:      false, // Cash deal
		TitleCommitmentDays: 20,
	}

	deadlines := engine.CalculateDeadlines(contract, nil)

	// Count financing-related deadlines
	financingDeadlines := 0
	for _, d := range deadlines {
		if d.Type == "loan_approval" || d.Type == "appraisal" {
			financingDeadlines++
		}
	}

	assert.Equal(t, 0, financingDeadlines, "Cash deal should not have financing deadlines")
}

func TestCalculateDeadlinesWithSellerLease(t *testing.T) {
	engine := NewDeadlineEngine()

	closingDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	contract := &models.ContractTerms{
		ProjectID:            "test-project-3",
		EffectiveDate:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		ClosingDate:          closingDate,
		SellerStaysPostClose: true,
		SellerLeaseDays:      30,
		TitleCommitmentDays:  20,
	}

	deadlines := engine.CalculateDeadlines(contract, nil)

	// Find lease end deadline
	var leaseEndDeadline *models.Deadline
	for _, d := range deadlines {
		if d.Type == "lease_end" {
			leaseEndDeadline = d
			break
		}
	}

	assert.NotNil(t, leaseEndDeadline, "Should have lease end deadline when seller stays post-close")
	expectedLeaseEnd := closingDate.AddDate(0, 0, 30)
	assert.Equal(t, expectedLeaseEnd, leaseEndDeadline.DueDate, "Lease end should be 30 days after closing")
	assert.Equal(t, "critical", leaseEndDeadline.Priority, "Lease end should be critical priority")
}

func TestNthWeekday(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		weekday  time.Weekday
		n        int
		expected time.Time
	}{
		{
			name:     "MLK Day 2024 (3rd Monday in January)",
			year:     2024,
			month:    time.January,
			weekday:  time.Monday,
			n:        3,
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Thanksgiving 2024 (4th Thursday in November)",
			year:     2024,
			month:    time.November,
			weekday:  time.Thursday,
			n:        4,
			expected: time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nthWeekday(tt.year, tt.month, tt.weekday, tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLastWeekday(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		weekday  time.Weekday
		expected time.Time
	}{
		{
			name:     "Memorial Day 2024 (last Monday in May)",
			year:     2024,
			month:    time.May,
			weekday:  time.Monday,
			expected: time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Last Friday in December 2024",
			year:     2024,
			month:    time.December,
			weekday:  time.Friday,
			expected: time.Date(2024, 12, 27, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastWeekday(tt.year, tt.month, tt.weekday)
			assert.Equal(t, tt.expected, result)
		})
	}
}
