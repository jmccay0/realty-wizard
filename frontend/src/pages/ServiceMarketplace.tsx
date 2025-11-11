import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { listServices, type Service } from '../api';

const CATEGORY_LABELS: Record<string, string> = {
  inspections: 'Inspections & Assessments',
  contracts_title_legal: 'Contracts, Title & Legal',
  repairs_trades: 'Repairs & Trades',
};

const SERVICE_IMAGES: Record<string, string> = {
  'Home Inspector': 'https://images.unsplash.com/photo-1560518883-ce09059eeffa?w=400&h=300&fit=crop',
  'Pest Inspector': 'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=400&h=300&fit=crop',
  'Roof Inspector': 'https://images.unsplash.com/photo-1625525675475-37addb0f3677?w=400&h=300&fit=crop',
  'Sewer Line Inspector': 'https://images.unsplash.com/photo-1607472586893-edb57bdc0e39?w=400&h=300&fit=crop',
  'HVAC Inspector': 'https://images.unsplash.com/photo-1607400201515-c2c41c07d307?w=400&h=300&fit=crop',
  'Pool/Spa Inspector': 'https://images.unsplash.com/photo-1576013551627-0cc20b96c2a7?w=400&h=300&fit=crop',
  'Environmental Inspector': 'https://images.unsplash.com/photo-1532996122724-e3c354a0b15b?w=400&h=300&fit=crop',
  'Foundation Inspector': 'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=400&h=300&fit=crop',
  'Land Surveyor': 'https://images.unsplash.com/photo-1589939705384-5185137a7f0f?w=400&h=300&fit=crop',
  'Real Estate Attorney': 'https://images.unsplash.com/photo-1589829545856-d10d557cf95f?w=400&h=300&fit=crop',
  'Title Agent/Company': 'https://images.unsplash.com/photo-1450101499163-c8848c66ca85?w=400&h=300&fit=crop',
  'Mobile Notary': 'https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400&h=300&fit=crop',
  'Appraiser': 'https://images.unsplash.com/photo-1560520653-9e0e4c89eb11?w=400&h=300&fit=crop',
  'Mortgage Loan Organizer': 'https://images.unsplash.com/photo-1554224154-26032ffc0d07?w=400&h=300&fit=crop',
  'Document Runner': 'https://images.unsplash.com/photo-1565728744382-61accd4aa148?w=400&h=300&fit=crop',
  'Due Diligence Coordinator': 'https://images.unsplash.com/photo-1551836022-d5d88e9218df?w=400&h=300&fit=crop',
  'Transaction Coordinator': 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=400&h=300&fit=crop',
  'Insurance Verification Specialist': 'https://images.unsplash.com/photo-1450101499163-c8848c66ca85?w=400&h=300&fit=crop',
  'Licensed Electrician': 'https://images.unsplash.com/photo-1621905251189-08b45d6a269e?w=400&h=300&fit=crop',
  'Licensed Plumber': 'https://images.unsplash.com/photo-1607472586893-edb57bdc0e39?w=400&h=300&fit=crop',
  'General Contractor': 'https://images.unsplash.com/photo-1504307651254-35680f356dfd?w=400&h=300&fit=crop',
  'Roofing Contractor': 'https://images.unsplash.com/photo-1565008576549-57569a49371d?w=400&h=300&fit=crop',
  'HVAC Technician': 'https://images.unsplash.com/photo-1607400201889-565b1ee75f8e?w=400&h=300&fit=crop',
  'Pest Control Company': 'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=400&h=300&fit=crop',
  'Appliance Repair': 'https://images.unsplash.com/photo-1556911220-bff31c812dba?w=400&h=300&fit=crop',
  'Mold Remediation': 'https://images.unsplash.com/photo-1628177142898-93e36e4e3a50?w=400&h=300&fit=crop',
  'Painting Contractor': 'https://images.unsplash.com/photo-1562259949-e8e7689d7828?w=400&h=300&fit=crop',
  'Flooring Specialist': 'https://images.unsplash.com/photo-1615874959474-d609969a20ed?w=400&h=300&fit=crop',
};

function ServiceMarketplace() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const navigate = useNavigate();

  useEffect(() => {
    loadServices();
  }, [selectedCategory]);

  const loadServices = async () => {
    try {
      const category = selectedCategory === 'all' ? undefined : selectedCategory;
      const response = await listServices(category);
      setServices(response.data);
    } catch (error) {
      console.error('Failed to load services:', error);
    } finally {
      setLoading(false);
    }
  };

  // Group services by category
  const servicesByCategory: Record<string, Service[]> = {};
  services.forEach(service => {
    if (!servicesByCategory[service.category]) {
      servicesByCategory[service.category] = [];
    }
    servicesByCategory[service.category].push(service);
  });

  const formatCost = (service: Service) => {
    if (!service.estimated_cost_min && !service.estimated_cost_max) {
      return 'Contact for pricing';
    }
    if (service.estimated_cost_min && service.estimated_cost_max) {
      return `$${service.estimated_cost_min.toLocaleString()} - $${service.estimated_cost_max.toLocaleString()}`;
    }
    return `From $${(service.estimated_cost_min || service.estimated_cost_max)?.toLocaleString()}`;
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f8f9fa', padding: '20px' }}>
      {/* Header */}
      <div style={{
        maxWidth: '1400px',
        margin: '0 auto 30px',
        background: 'white',
        padding: '24px',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <h1 style={{ fontSize: '32px', margin: 0, color: '#2c3e50' }}>
            Service Marketplace
          </h1>
          <button
            onClick={() => navigate('/')}
            className="btn-secondary"
          >
            ← Back to Home
          </button>
        </div>
        <p style={{ margin: 0, color: '#6c757d', fontSize: '16px' }}>
          Find trusted professionals for your real estate transaction. Click any service to view available providers.
        </p>
      </div>

      <div style={{ maxWidth: '1400px', margin: '0 auto' }}>
        {/* Category Filter */}
        <div style={{ marginBottom: '30px', display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
          <button
            className={selectedCategory === 'all' ? 'btn-primary' : 'btn-secondary'}
            onClick={() => setSelectedCategory('all')}
            style={{ fontSize: '15px', padding: '10px 20px' }}
          >
            All Services ({services.length})
          </button>
          {Object.keys(CATEGORY_LABELS).map(cat => (
            <button
              key={cat}
              className={selectedCategory === cat ? 'btn-primary' : 'btn-secondary'}
              onClick={() => setSelectedCategory(cat)}
              style={{ fontSize: '15px', padding: '10px 20px' }}
            >
              {CATEGORY_LABELS[cat]}
            </button>
          ))}
        </div>

        {/* Services Grid */}
        {loading ? (
          <div style={{ display: 'flex', justifyContent: 'center', padding: '60px' }}>
            <div className="spinner"></div>
          </div>
        ) : services.length === 0 ? (
          <div style={{
            background: 'white',
            padding: '60px',
            textAlign: 'center',
            borderRadius: '8px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          }}>
            <p style={{ color: '#6c757d', fontSize: '18px' }}>No services found in this category.</p>
          </div>
        ) : (
          <div>
            {Object.keys(servicesByCategory).sort().map(category => (
              <div key={category} style={{ marginBottom: '40px' }}>
                <h2 style={{
                  fontSize: '24px',
                  marginBottom: '20px',
                  color: '#2c3e50',
                  borderLeft: '4px solid var(--realwiz-green)',
                  paddingLeft: '16px',
                }}>
                  {CATEGORY_LABELS[category] || category}
                </h2>
                <div style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(auto-fill, minmax(340px, 1fr))',
                  gap: '20px',
                }}>
                  {servicesByCategory[category].map(service => (
                    <div
                      key={service.id}
                      onClick={() => navigate(`/services/${service.id}`)}
                      style={{
                        background: 'white',
                        borderRadius: '8px',
                        overflow: 'hidden',
                        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                        cursor: 'pointer',
                        transition: 'all 0.3s ease',
                        border: '1px solid #e9ecef',
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'translateY(-4px)';
                        e.currentTarget.style.boxShadow = '0 8px 20px rgba(0,0,0,0.15)';
                        e.currentTarget.style.borderColor = 'var(--realwiz-green)';
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'translateY(0)';
                        e.currentTarget.style.boxShadow = '0 2px 8px rgba(0,0,0,0.1)';
                        e.currentTarget.style.borderColor = '#e9ecef';
                      }}
                    >
                      {/* Image Header */}
                      <div style={{
                        position: 'relative',
                        height: '200px',
                        overflow: 'hidden',
                      }}>
                        <img
                          src={SERVICE_IMAGES[service.name] || 'https://images.unsplash.com/photo-1560518883-ce09059eeffa?w=400&h=300&fit=crop'}
                          alt={service.name}
                          style={{
                            width: '100%',
                            height: '100%',
                            objectFit: 'cover',
                          }}
                        />
                        <div style={{
                          position: 'absolute',
                          bottom: 0,
                          left: 0,
                          right: 0,
                          background: 'linear-gradient(to top, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0.4) 100%)',
                          padding: '16px',
                        }}>
                          <h3 style={{
                            fontSize: '20px',
                            margin: 0,
                            color: 'white',
                            fontWeight: '600',
                            textShadow: '0 2px 4px rgba(0,0,0,0.5)',
                          }}>
                            {service.name}
                          </h3>
                        </div>
                      </div>

                      {/* Content */}
                      <div style={{ padding: '20px' }}>
                        <p style={{
                          fontSize: '14px',
                          color: '#6c757d',
                          marginBottom: '16px',
                          minHeight: '60px',
                          lineHeight: '1.5',
                        }}>
                          {service.description}
                        </p>

                        {/* Details */}
                        <div style={{ marginBottom: '16px' }}>
                          {service.typical_timeline && (
                            <div style={{
                              fontSize: '13px',
                              color: '#495057',
                              marginBottom: '8px',
                              display: 'flex',
                              alignItems: 'center',
                              gap: '8px',
                            }}>
                              <span style={{ fontSize: '16px' }}>⏱️</span>
                              <span>{service.typical_timeline}</span>
                            </div>
                          )}
                          <div style={{
                            fontSize: '13px',
                            color: '#495057',
                            display: 'flex',
                            alignItems: 'center',
                            gap: '8px',
                          }}>
                            <span style={{ fontSize: '16px' }}>💵</span>
                            <span style={{ fontWeight: '600', color: 'var(--realwiz-green)' }}>
                              {formatCost(service)}
                            </span>
                          </div>
                        </div>

                        {/* Call to Action */}
                        <div style={{
                          paddingTop: '16px',
                          borderTop: '1px solid #e9ecef',
                        }}>
                          <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'space-between',
                            color: 'var(--realwiz-green)',
                            fontWeight: '600',
                            fontSize: '14px',
                          }}>
                            <span>View Providers</span>
                            <span style={{ fontSize: '18px' }}>→</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default ServiceMarketplace;
