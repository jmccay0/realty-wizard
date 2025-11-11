package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/storage"
)

func main() {
	// Check for environment variable to include sample providers
	includeSampleProviders := os.Getenv("SEED_SAMPLE_PROVIDERS") != "false"

	// Initialize storage
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	dbPath := filepath.Join(dataDir, "realty-wizard.db")
	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	log.Println("Seeding services...")
	if includeSampleProviders {
		log.Println("(Sample providers will be included)")
	} else {
		log.Println("(Sample providers will NOT be included)")
	}

	now := time.Now()

	// Inspections Category
	inspectionServices := []models.Service{
		{
			ID:               uuid.New().String(),
			Name:             "Home Inspector",
			Category:         "inspections",
			Description:      "General home inspection covering structure, systems, and safety",
			TypicalTimeline:  "2-3 hours on-site, report within 24 hours",
			EstimatedCostMin: 300,
			EstimatedCostMax: 600,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Pest Inspector",
			Category:         "inspections",
			Description:      "Inspection for termites, wood-destroying insects, and pest damage",
			TypicalTimeline:  "1-2 hours on-site, report within 24 hours",
			EstimatedCostMin: 100,
			EstimatedCostMax: 200,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Roof Inspector",
			Category:         "inspections",
			Description:      "Specialized inspection of roof condition, materials, and lifespan",
			TypicalTimeline:  "1-2 hours on-site, report same day",
			EstimatedCostMin: 150,
			EstimatedCostMax: 350,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Sewer Line Inspector",
			Category:         "inspections",
			Description:      "Camera inspection of sewer lines for damage or blockages",
			TypicalTimeline:  "1-2 hours on-site, report same day",
			EstimatedCostMin: 200,
			EstimatedCostMax: 400,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "HVAC Inspector",
			Category:         "inspections",
			Description:      "Detailed inspection of heating and cooling systems",
			TypicalTimeline:  "1-2 hours on-site, report within 24 hours",
			EstimatedCostMin: 150,
			EstimatedCostMax: 300,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Pool/Spa Inspector",
			Category:         "inspections",
			Description:      "Inspection of pool, spa, and related equipment",
			TypicalTimeline:  "1-2 hours on-site, report same day",
			EstimatedCostMin: 200,
			EstimatedCostMax: 400,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Environmental Inspector",
			Category:         "inspections",
			Description:      "Testing for radon, mold, asbestos, or lead paint",
			TypicalTimeline:  "Varies by test, results in 3-7 days",
			EstimatedCostMin: 200,
			EstimatedCostMax: 800,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Foundation Inspector",
			Category:         "inspections",
			Description:      "Specialized structural engineer inspection of foundation",
			TypicalTimeline:  "2-3 hours on-site, report within 2-3 days",
			EstimatedCostMin: 400,
			EstimatedCostMax: 800,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Land Surveyor",
			Category:         "inspections",
			Description:      "Official property boundary survey",
			TypicalTimeline:  "1-2 weeks from request to delivery",
			EstimatedCostMin: 400,
			EstimatedCostMax: 1000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	// Contracts/Title/Legal Category
	legalServices := []models.Service{
		{
			ID:               uuid.New().String(),
			Name:             "Real Estate Attorney",
			Category:         "contracts_title_legal",
			Description:      "Legal review of contracts and closing documents",
			TypicalTimeline:  "2-5 business days for review",
			EstimatedCostMin: 500,
			EstimatedCostMax: 2000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Title Agent/Company",
			Category:         "contracts_title_legal",
			Description:      "Title search, insurance, and closing services",
			TypicalTimeline:  "7-14 days for full title work",
			EstimatedCostMin: 1000,
			EstimatedCostMax: 3000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Mobile Notary",
			Category:         "contracts_title_legal",
			Description:      "Notarization of documents at your location",
			TypicalTimeline:  "Same day or next day service",
			EstimatedCostMin: 50,
			EstimatedCostMax: 150,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Appraiser",
			Category:         "contracts_title_legal",
			Description:      "Professional property valuation for financing",
			TypicalTimeline:  "3-7 days from inspection to report",
			EstimatedCostMin: 400,
			EstimatedCostMax: 700,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Mortgage Loan Organizer",
			Category:         "contracts_title_legal",
			Description:      "Assistance with loan application and document prep",
			TypicalTimeline:  "2-4 weeks for full loan process",
			EstimatedCostMin: 300,
			EstimatedCostMax: 1000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Document Runner",
			Category:         "contracts_title_legal",
			Description:      "Pick up and deliver time-sensitive documents",
			TypicalTimeline:  "Same day service",
			EstimatedCostMin: 50,
			EstimatedCostMax: 200,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Due Diligence Coordinator",
			Category:         "contracts_title_legal",
			Description:      "Coordinate inspections and deadline tracking",
			TypicalTimeline:  "Throughout transaction (30-45 days)",
			EstimatedCostMin: 500,
			EstimatedCostMax: 1500,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Transaction Coordinator",
			Category:         "contracts_title_legal",
			Description:      "Full transaction management from contract to close",
			TypicalTimeline:  "Throughout transaction (30-45 days)",
			EstimatedCostMin: 300,
			EstimatedCostMax: 800,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Insurance Verification Specialist",
			Category:         "contracts_title_legal",
			Description:      "Obtain insurance quotes and verify coverage requirements",
			TypicalTimeline:  "2-5 business days",
			EstimatedCostMin: 100,
			EstimatedCostMax: 300,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	// Repairs/Trades Category
	tradeServices := []models.Service{
		{
			ID:               uuid.New().String(),
			Name:             "Licensed Electrician",
			Category:         "repairs_trades",
			Description:      "Electrical repairs, upgrades, and safety work",
			TypicalTimeline:  "1-7 days depending on scope",
			EstimatedCostMin: 150,
			EstimatedCostMax: 2000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Licensed Plumber",
			Category:         "repairs_trades",
			Description:      "Plumbing repairs, fixture installation, water heater work",
			TypicalTimeline:  "1-7 days depending on scope",
			EstimatedCostMin: 150,
			EstimatedCostMax: 2000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "General Contractor",
			Category:         "repairs_trades",
			Description:      "General repairs, handyman work, and project management",
			TypicalTimeline:  "Varies by project size",
			EstimatedCostMin: 200,
			EstimatedCostMax: 5000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Roofing Contractor",
			Category:         "repairs_trades",
			Description:      "Roof repairs, replacement, and emergency patching",
			TypicalTimeline:  "1-14 days depending on scope",
			EstimatedCostMin: 500,
			EstimatedCostMax: 15000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "HVAC Technician",
			Category:         "repairs_trades",
			Description:      "AC/heating repair, maintenance, and replacement",
			TypicalTimeline:  "1-7 days depending on scope",
			EstimatedCostMin: 150,
			EstimatedCostMax: 8000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Pest Control Company",
			Category:         "repairs_trades",
			Description:      "Pest treatment, termite treatment, and prevention",
			TypicalTimeline:  "1-7 days, ongoing monitoring available",
			EstimatedCostMin: 150,
			EstimatedCostMax: 2500,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Appliance Repair",
			Category:         "repairs_trades",
			Description:      "Repair of major appliances (dishwasher, oven, etc.)",
			TypicalTimeline:  "1-5 days",
			EstimatedCostMin: 100,
			EstimatedCostMax: 500,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Mold Remediation",
			Category:         "repairs_trades",
			Description:      "Mold removal, treatment, and prevention",
			TypicalTimeline:  "3-14 days depending on severity",
			EstimatedCostMin: 500,
			EstimatedCostMax: 6000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Painting Contractor",
			Category:         "repairs_trades",
			Description:      "Interior/exterior painting for staging or repairs",
			TypicalTimeline:  "2-7 days depending on scope",
			EstimatedCostMin: 300,
			EstimatedCostMax: 5000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New().String(),
			Name:             "Flooring Specialist",
			Category:         "repairs_trades",
			Description:      "Floor repair, refinishing, or replacement",
			TypicalTimeline:  "3-10 days depending on scope",
			EstimatedCostMin: 500,
			EstimatedCostMax: 8000,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	// Insert all services
	allServices := append(append(inspectionServices, legalServices...), tradeServices...)
	for _, service := range allServices {
		if err := store.CreateService(&service); err != nil {
			log.Printf("Failed to create service %s: %v", service.Name, err)
		} else {
			log.Printf("Created service: %s", service.Name)
		}
	}

	log.Printf("Successfully seeded %d services!", len(allServices))

	// Only seed providers if flag is set
	if !includeSampleProviders {
		log.Println("\nSkipping sample providers (SEED_SAMPLE_PROVIDERS=false)")
		log.Println("Seed complete!")
		return
	}

	// Create some sample providers for the first few services
	log.Println("\nSeeding sample providers...")

	sampleProviders := []models.Provider{
		{
			ID:                 uuid.New().String(),
			ServiceID:          inspectionServices[0].ID, // Home Inspector
			BusinessName:       "Lone Star Home Inspections",
			ContactName:        "Mike Johnson",
			Email:              "mike@lonestarinspections.com",
			Phone:              "(512) 555-0100",
			Address:            "123 Main St",
			City:               "Austin",
			State:              "TX",
			Zip:                "78701",
			Bio:                "20+ years experience in home inspections. Licensed, insured, and TREC certified.",
			YearsExperience:    20,
			LicenseNumber:      "TREC #12345",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.8,
			RatingCount:        142,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          inspectionServices[0].ID, // Home Inspector
			BusinessName:       "Texas Property Inspectors",
			ContactName:        "Sarah Williams",
			Email:              "sarah@txpropertyinspectors.com",
			Phone:              "(512) 555-0101",
			Address:            "456 Oak Ave",
			City:               "Austin",
			State:              "TX",
			Zip:                "78702",
			Bio:                "Certified home inspector specializing in older homes and historical properties.",
			YearsExperience:    15,
			LicenseNumber:      "TREC #12346",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.9,
			RatingCount:        98,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[0].ID, // Real Estate Attorney
			BusinessName:       "Miller & Associates Law",
			ContactName:        "Jennifer Miller",
			Email:              "jmiller@millerlaw.com",
			Phone:              "(512) 555-0200",
			Address:            "789 Congress Ave, Suite 500",
			City:               "Austin",
			State:              "TX",
			Zip:                "78701",
			Bio:                "Real estate law firm with 30 years serving Austin area. Specializing in residential transactions.",
			YearsExperience:    30,
			LicenseNumber:      "State Bar #54321",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.7,
			RatingCount:        76,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          tradeServices[0].ID, // Electrician
			BusinessName:       "Austin Electric Pros",
			ContactName:        "David Chen",
			Email:              "david@austinelectricpros.com",
			Phone:              "(512) 555-0300",
			Address:            "321 Industrial Blvd",
			City:               "Austin",
			State:              "TX",
			Zip:                "78745",
			Bio:                "Licensed master electrician. 24/7 emergency service. Specializing in residential repairs and code compliance.",
			YearsExperience:    18,
			LicenseNumber:      "Master Electrician #78901",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.9,
			RatingCount:        203,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          tradeServices[1].ID, // Plumber
			BusinessName:       "Reliable Plumbing Services",
			ContactName:        "Tom Rodriguez",
			Email:              "tom@reliableplumbing.com",
			Phone:              "(512) 555-0301",
			Address:            "654 Service Dr",
			City:               "Austin",
			State:              "TX",
			Zip:                "78748",
			Bio:                "Family-owned plumbing company. Licensed, bonded, insured. Same-day service available.",
			YearsExperience:    25,
			LicenseNumber:      "Master Plumber #45678",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.6,
			RatingCount:        167,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// Additional Pest Inspectors
		{
			ID:                 uuid.New().String(),
			ServiceID:          inspectionServices[1].ID, // Pest Inspector
			BusinessName:       "Texas Termite & Pest Control",
			ContactName:        "James Patterson",
			Email:              "james@texastermite.com",
			Phone:              "(512) 555-0400",
			City:               "Austin",
			State:              "TX",
			Bio:                "Specializing in termite inspections and WDI reports for real estate transactions.",
			YearsExperience:    15,
			LicenseNumber:      "TDA License #8901",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.7,
			RatingCount:        89,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          inspectionServices[1].ID, // Pest Inspector
			BusinessName:       "All Clear Pest Inspections",
			ContactName:        "Maria Garcia",
			Email:              "maria@allclearpest.com",
			Phone:              "(512) 555-0401",
			City:               "Round Rock",
			State:              "TX",
			Bio:                "Fast turnaround on pest inspection reports. Same-day service available.",
			YearsExperience:    10,
			LicenseNumber:      "TDA License #8902",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.5,
			RatingCount:        56,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// Roof Inspectors
		{
			ID:                 uuid.New().String(),
			ServiceID:          inspectionServices[2].ID, // Roof Inspector
			BusinessName:       "Central Texas Roof Inspections",
			ContactName:        "Bob Anderson",
			Email:              "bob@centraltexasroof.com",
			Phone:              "(512) 555-0500",
			City:               "Austin",
			State:              "TX",
			Bio:                "Certified roof inspector with drone technology for detailed assessments.",
			YearsExperience:    18,
			LicenseNumber:      "TX Roofing #5678",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.9,
			RatingCount:        134,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// Title Agents
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[1].ID, // Title Agent/Company
			BusinessName:       "Austin Title Company",
			ContactName:        "Rebecca Thompson",
			Email:              "rebecca@austintitle.com",
			Phone:              "(512) 555-0600",
			Address:            "100 Congress Ave, Suite 800",
			City:               "Austin",
			State:              "TX",
			Zip:                "78701",
			Bio:                "Full-service title company handling residential and commercial closings. In business since 1985.",
			YearsExperience:    35,
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.8,
			RatingCount:        245,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[1].ID, // Title Agent/Company
			BusinessName:       "Lone Star Title Services",
			ContactName:        "Mark Stevens",
			Email:              "mark@lonestartitle.com",
			Phone:              "(512) 555-0601",
			City:               "Cedar Park",
			State:              "TX",
			Bio:                "Fast, reliable title services with mobile closing available.",
			YearsExperience:    20,
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.6,
			RatingCount:        178,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// Mobile Notaries
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[2].ID, // Mobile Notary
			BusinessName:       "24/7 Mobile Notary Austin",
			ContactName:        "Linda Martinez",
			Email:              "linda@247notary.com",
			Phone:              "(512) 555-0700",
			City:               "Austin",
			State:              "TX",
			Bio:                "Available 24/7 for emergency closings. Will travel anywhere in Travis County.",
			YearsExperience:    8,
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      5.0,
			RatingCount:        92,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// HVAC Technicians
		{
			ID:                 uuid.New().String(),
			ServiceID:          tradeServices[4].ID, // HVAC Technician
			BusinessName:       "Cool Breeze HVAC",
			ContactName:        "Kevin White",
			Email:              "kevin@coolbreezehvac.com",
			Phone:              "(512) 555-0800",
			City:               "Austin",
			State:              "TX",
			Bio:                "AC and heating repair, replacement, and maintenance. Emergency service available.",
			YearsExperience:    22,
			LicenseNumber:      "TACL #12345",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.7,
			RatingCount:        156,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          tradeServices[4].ID, // HVAC Technician
			BusinessName:       "Texas Air Specialists",
			ContactName:        "Amanda Brown",
			Email:              "amanda@texasairspecialists.com",
			Phone:              "(512) 555-0801",
			City:               "Pflugerville",
			State:              "TX",
			Bio:                "Specializing in energy-efficient HVAC systems and repairs.",
			YearsExperience:    15,
			LicenseNumber:      "TACL #12346",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.8,
			RatingCount:        134,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// General Contractors
		{
			ID:                 uuid.New().String(),
			ServiceID:          tradeServices[2].ID, // General Contractor
			BusinessName:       "Austin Home Renovations",
			ContactName:        "Steve Jackson",
			Email:              "steve@austinhomereno.com",
			Phone:              "(512) 555-0900",
			City:               "Austin",
			State:              "TX",
			Bio:                "Full-service general contractor for repairs and renovations. Licensed and insured.",
			YearsExperience:    28,
			LicenseNumber:      "TX Contractor #9876",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.6,
			RatingCount:        201,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		// Appraisers
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[3].ID, // Appraiser
			BusinessName:       "Central Texas Appraisal Services",
			ContactName:        "Robert Taylor",
			Email:              "robert@centraltexasappraisal.com",
			Phone:              "(512) 555-1000",
			City:               "Austin",
			State:              "TX",
			Bio:                "State-certified residential appraiser. Fast turnaround for purchase and refinance appraisals.",
			YearsExperience:    16,
			LicenseNumber:      "TX Appraiser #54321",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.7,
			RatingCount:        112,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New().String(),
			ServiceID:          legalServices[3].ID, // Appraiser
			BusinessName:       "Hill Country Appraisals",
			ContactName:        "Nancy Wilson",
			Email:              "nancy@hillcountryappraisals.com",
			Phone:              "(512) 555-1001",
			City:               "Dripping Springs",
			State:              "TX",
			Bio:                "Specializing in residential appraisals throughout the Austin metro area.",
			YearsExperience:    12,
			LicenseNumber:      "TX Appraiser #54322",
			InsuranceVerified:  true,
			AvailabilityStatus: "available",
			RatingAverage:      4.9,
			RatingCount:        87,
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	for _, provider := range sampleProviders {
		if err := store.CreateProvider(&provider); err != nil {
			log.Printf("Failed to create provider %s: %v", provider.BusinessName, err)
		} else {
			log.Printf("Created provider: %s", provider.BusinessName)
		}
	}

	log.Printf("Successfully seeded %d sample providers!", len(sampleProviders))
	log.Println("\nSeed complete!")
}
