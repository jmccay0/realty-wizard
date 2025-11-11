package storage

import (
	"github.com/jessicaandtommymccay/realty-wizard/backend/internal/models"
)

// Storage defines the interface for data persistence
// This clean interface makes testing easy and allows swapping implementations
type Storage interface {
	// Project operations
	CreateProject(project *models.Project) error
	GetProject(id string) (*models.Project, error)
	UpdateProject(project *models.Project) error
	ListProjects() ([]*models.Project, error)

	// Property operations
	CreateProperty(property *models.Property) error
	GetProperty(projectID string) (*models.Property, error)
	UpdateProperty(property *models.Property) error

	// Disclosure operations
	CreateDisclosure(disclosure *models.Disclosure) error
	GetDisclosure(projectID string) (*models.Disclosure, error)
	UpdateDisclosure(disclosure *models.Disclosure) error

	// Contract operations
	CreateContract(contract *models.ContractTerms) error
	GetContract(projectID string) (*models.ContractTerms, error)
	UpdateContract(contract *models.ContractTerms) error

	// Deadline operations
	CreateDeadline(deadline *models.Deadline) error
	GetDeadline(id string) (*models.Deadline, error)
	ListDeadlines(projectID string) ([]*models.Deadline, error)
	UpdateDeadline(deadline *models.Deadline) error
	DeleteDeadline(id string) error

	// Document operations
	CreateDocument(doc *models.Document) error
	GetDocument(id string) (*models.Document, error)
	ListDocuments(projectID string) ([]*models.Document, error)
	UpdateDocument(doc *models.Document) error

	// Service Marketplace operations
	// Service operations
	CreateService(service *models.Service) error
	GetService(id string) (*models.Service, error)
	ListServices() ([]*models.Service, error)
	ListServicesByCategory(category string) ([]*models.Service, error)

	// Provider operations
	CreateProvider(provider *models.Provider) error
	GetProvider(id string) (*models.Provider, error)
	ListProviders(serviceID string) ([]*models.Provider, error)
	UpdateProvider(provider *models.Provider) error

	// Service Request operations
	CreateServiceRequest(request *models.ServiceRequest) error
	GetServiceRequest(id string) (*models.ServiceRequest, error)
	ListServiceRequests(userEmail string) ([]*models.ServiceRequest, error)
	UpdateServiceRequest(request *models.ServiceRequest) error

	// Provider Review operations
	CreateProviderReview(review *models.ProviderReview) error
	ListProviderReviews(providerID string) ([]*models.ProviderReview, error)

	// Utility
	Close() error
}
