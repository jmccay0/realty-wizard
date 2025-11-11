import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getService, listProviders, listProviderReviews, type Service, type Provider, type ProviderReview } from '../api';

function ServiceDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [service, setService] = useState<Service | null>(null);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [selectedProvider, setSelectedProvider] = useState<Provider | null>(null);
  const [providerReviews, setProviderReviews] = useState<ProviderReview[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      loadServiceAndProviders();
    }
  }, [id]);

  const loadServiceAndProviders = async () => {
    try {
      const [serviceRes, providersRes] = await Promise.all([
        getService(id!),
        listProviders(id!),
      ]);
      setService(serviceRes.data);
      setProviders(providersRes.data || []); // Handle null response
    } catch (error) {
      console.error('Failed to load service details:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadProviderReviews = async (providerId: string) => {
    try {
      const response = await listProviderReviews(providerId);
      setProviderReviews(response.data || []); // Handle null response
    } catch (error) {
      console.error('Failed to load reviews:', error);
    }
  };

  const handleProviderClick = (provider: Provider) => {
    setSelectedProvider(provider);
    loadProviderReviews(provider.id);
  };

  const formatCost = (service: Service) => {
    if (!service.estimated_cost_min && !service.estimated_cost_max) {
      return 'Contact for pricing';
    }
    if (service.estimated_cost_min && service.estimated_cost_max) {
      return `$${service.estimated_cost_min} - $${service.estimated_cost_max}`;
    }
    return `From $${service.estimated_cost_min || service.estimated_cost_max}`;
  };

  const renderStars = (rating: number) => {
    return '★'.repeat(Math.round(rating)) + '☆'.repeat(5 - Math.round(rating));
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  if (!service) {
    return (
      <div style={{ padding: '40px', textAlign: 'center' }}>
        <h2>Service not found</h2>
        <button className="btn-primary" onClick={() => navigate('/services')}>
          Back to Services
        </button>
      </div>
    );
  }

  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--realwiz-gray-50)', padding: '20px' }}>
      {/* Header */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto 20px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}>
        <button
          onClick={() => navigate('/services')}
          className="btn-secondary"
        >
          ← Back to Services
        </button>
      </div>

      <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
        {/* Service Info */}
        <div className="card" style={{ marginBottom: '24px' }}>
          <h1 style={{ fontSize: '28px', marginBottom: '12px', color: 'var(--realwiz-green)' }}>
            {service.name}
          </h1>
          <p style={{ fontSize: '16px', color: 'var(--realwiz-gray-700)', marginBottom: '16px' }}>
            {service.description}
          </p>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
            {service.typical_timeline && (
              <div>
                <strong>Typical Timeline:</strong>
                <p style={{ marginTop: '4px', color: 'var(--realwiz-gray-700)' }}>{service.typical_timeline}</p>
              </div>
            )}
            <div>
              <strong>Estimated Cost:</strong>
              <p style={{ marginTop: '4px', color: 'var(--realwiz-green)', fontWeight: 'bold' }}>
                {formatCost(service)}
              </p>
            </div>
          </div>
        </div>

        {/* Providers List */}
        <h2 style={{ fontSize: '24px', marginBottom: '16px' }}>
          Available Providers ({providers.length})
        </h2>

        {providers.length === 0 ? (
          <div className="card" style={{ textAlign: 'center', padding: '40px' }}>
            <p style={{ color: 'var(--realwiz-gray-600)' }}>
              No providers available yet for this service. Check back soon!
            </p>
          </div>
        ) : (
          <div style={{
            display: 'grid',
            gridTemplateColumns: selectedProvider ? '1fr 1fr' : '1fr',
            gap: '20px',
          }}>
            {/* Providers Column */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              {providers.map(provider => (
                <div
                  key={provider.id}
                  className="card"
                  style={{
                    cursor: 'pointer',
                    border: selectedProvider?.id === provider.id ? '2px solid var(--realwiz-green)' : undefined,
                  }}
                  onClick={() => handleProviderClick(provider)}
                >
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                    <div style={{ flex: 1 }}>
                      <h3 style={{ fontSize: '18px', marginBottom: '4px' }}>
                        {provider.business_name}
                      </h3>
                      {provider.contact_name && (
                        <p style={{ fontSize: '14px', color: 'var(--realwiz-gray-600)', marginBottom: '8px' }}>
                          Contact: {provider.contact_name}
                        </p>
                      )}
                    </div>
                    <div style={{ textAlign: 'right' }}>
                      {provider.rating_count > 0 && (
                        <>
                          <div style={{ color: 'var(--realwiz-gold)', fontSize: '16px' }}>
                            {renderStars(provider.rating_average)}
                          </div>
                          <div style={{ fontSize: '12px', color: 'var(--realwiz-gray-600)' }}>
                            {provider.rating_average.toFixed(1)} ({provider.rating_count} reviews)
                          </div>
                        </>
                      )}
                    </div>
                  </div>

                  {provider.bio && (
                    <p style={{ fontSize: '14px', color: 'var(--realwiz-gray-700)', marginTop: '8px' }}>
                      {provider.bio}
                    </p>
                  )}

                  <div style={{ marginTop: '12px', display: 'flex', flexWrap: 'wrap', gap: '12px', fontSize: '13px' }}>
                    {provider.years_experience && (
                      <span className="badge badge-info">
                        {provider.years_experience} years exp.
                      </span>
                    )}
                    {provider.license_number && (
                      <span className="badge badge-info">
                        {provider.license_number}
                      </span>
                    )}
                    {provider.insurance_verified && (
                      <span className="badge badge-success">
                        Insured
                      </span>
                    )}
                    <span className={`badge badge-${
                      provider.availability_status === 'available' ? 'success' :
                      provider.availability_status === 'limited' ? 'warning' : 'secondary'
                    }`}>
                      {provider.availability_status}
                    </span>
                  </div>

                  <div style={{ marginTop: '12px', paddingTop: '12px', borderTop: '1px solid var(--realwiz-gray-300)' }}>
                    <p style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)', marginBottom: '4px' }}>
                      📞 {provider.phone}
                    </p>
                    <p style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)' }}>
                      📧 {provider.email}
                    </p>
                    {provider.city && provider.state && (
                      <p style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)', marginTop: '4px' }}>
                        📍 {provider.city}, {provider.state}
                      </p>
                    )}
                  </div>

                  <button
                    className="btn-primary"
                    style={{ marginTop: '12px', width: '100%' }}
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/services/${service.id}/request?provider=${provider.id}`);
                    }}
                  >
                    Request Service
                  </button>
                </div>
              ))}
            </div>

            {/* Provider Detail/Reviews Column */}
            {selectedProvider && (
              <div className="card" style={{ position: 'sticky', top: '20px', maxHeight: '80vh', overflowY: 'auto' }}>
                <h3 style={{ fontSize: '20px', marginBottom: '16px', color: 'var(--realwiz-green)' }}>
                  Reviews for {selectedProvider.business_name}
                </h3>

                {providerReviews.length === 0 ? (
                  <p style={{ color: 'var(--realwiz-gray-600)', textAlign: 'center', padding: '20px' }}>
                    No reviews yet
                  </p>
                ) : (
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    {providerReviews.map(review => (
                      <div
                        key={review.id}
                        style={{
                          padding: '12px',
                          backgroundColor: 'var(--realwiz-gray-50)',
                          borderRadius: '4px',
                        }}
                      >
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                          <div style={{ color: 'var(--realwiz-gold)' }}>
                            {renderStars(review.rating)}
                          </div>
                          <div style={{ fontSize: '12px', color: 'var(--realwiz-gray-600)' }}>
                            {new Date(review.created_at).toLocaleDateString()}
                          </div>
                        </div>
                        {review.review_text && (
                          <p style={{ fontSize: '14px', color: 'var(--realwiz-gray-700)' }}>
                            {review.review_text}
                          </p>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default ServiceDetail;
