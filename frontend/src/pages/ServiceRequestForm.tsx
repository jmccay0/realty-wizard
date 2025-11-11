import { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { getService, getProvider, createServiceRequest, type Service, type Provider } from '../api';

function ServiceRequestForm() {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const providerId = searchParams.get('provider');

  const [service, setService] = useState<Service | null>(null);
  const [provider, setProvider] = useState<Provider | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [formData, setFormData] = useState({
    user_name: '',
    user_email: '',
    user_phone: '',
    property_address: '',
    requested_date: '',
    preferred_time: '',
    notes: '',
  });

  useEffect(() => {
    loadData();
  }, [id, providerId]);

  const loadData = async () => {
    try {
      const serviceRes = await getService(id!);
      setService(serviceRes.data);

      if (providerId) {
        const providerRes = await getProvider(providerId);
        setProvider(providerRes.data);
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      await createServiceRequest({
        service_id: id!,
        provider_id: providerId || undefined,
        user_name: formData.user_name,
        user_email: formData.user_email,
        user_phone: formData.user_phone || undefined,
        property_address: formData.property_address || undefined,
        requested_date: formData.requested_date || undefined,
        preferred_time: formData.preferred_time || undefined,
        notes: formData.notes || undefined,
        status: 'pending',
      });

      alert('Service request submitted successfully! The provider will contact you soon.');
      navigate(`/services/${id}`);
    } catch (error) {
      console.error('Failed to submit request:', error);
      alert('Failed to submit request. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
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
      <div style={{ maxWidth: '800px', margin: '0 auto' }}>
        {/* Header */}
        <div style={{ marginBottom: '20px' }}>
          <button
            onClick={() => navigate(`/services/${id}`)}
            className="btn-secondary"
          >
            ← Back
          </button>
        </div>

        <div className="card">
          <h1 style={{ fontSize: '28px', marginBottom: '8px' }}>
            Request Service
          </h1>
          <p style={{ color: 'var(--realwiz-gray-600)', marginBottom: '24px' }}>
            {service.name}
            {provider && (
              <span style={{ display: 'block', marginTop: '4px', fontWeight: 'bold', color: 'var(--realwiz-green)' }}>
                Provider: {provider.business_name}
              </span>
            )}
          </p>

          <form onSubmit={handleSubmit}>
            {/* Your Information */}
            <div style={{ marginBottom: '24px' }}>
              <h3 style={{ fontSize: '18px', marginBottom: '12px', color: 'var(--realwiz-green)' }}>
                Your Information
              </h3>

              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                  Your Name <span style={{ color: 'red' }}>*</span>
                </label>
                <input
                  type="text"
                  name="user_name"
                  value={formData.user_name}
                  onChange={handleChange}
                  required
                  style={{
                    width: '100%',
                    padding: '8px 12px',
                    fontSize: '14px',
                    border: '1px solid var(--realwiz-gray-300)',
                    borderRadius: '4px',
                  }}
                />
              </div>

              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                  Email <span style={{ color: 'red' }}>*</span>
                </label>
                <input
                  type="email"
                  name="user_email"
                  value={formData.user_email}
                  onChange={handleChange}
                  required
                  style={{
                    width: '100%',
                    padding: '8px 12px',
                    fontSize: '14px',
                    border: '1px solid var(--realwiz-gray-300)',
                    borderRadius: '4px',
                  }}
                />
              </div>

              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                  Phone
                </label>
                <input
                  type="tel"
                  name="user_phone"
                  value={formData.user_phone}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: '8px 12px',
                    fontSize: '14px',
                    border: '1px solid var(--realwiz-gray-300)',
                    borderRadius: '4px',
                  }}
                />
              </div>
            </div>

            {/* Service Details */}
            <div style={{ marginBottom: '24px' }}>
              <h3 style={{ fontSize: '18px', marginBottom: '12px', color: 'var(--realwiz-green)' }}>
                Service Details
              </h3>

              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                  Property Address
                </label>
                <input
                  type="text"
                  name="property_address"
                  value={formData.property_address}
                  onChange={handleChange}
                  placeholder="123 Main St, Austin, TX 78701"
                  style={{
                    width: '100%',
                    padding: '8px 12px',
                    fontSize: '14px',
                    border: '1px solid var(--realwiz-gray-300)',
                    borderRadius: '4px',
                  }}
                />
              </div>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '16px' }}>
                <div>
                  <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                    Preferred Date
                  </label>
                  <input
                    type="date"
                    name="requested_date"
                    value={formData.requested_date}
                    onChange={handleChange}
                    min={new Date().toISOString().split('T')[0]}
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      fontSize: '14px',
                      border: '1px solid var(--realwiz-gray-300)',
                      borderRadius: '4px',
                    }}
                  />
                </div>

                <div>
                  <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                    Preferred Time
                  </label>
                  <input
                    type="text"
                    name="preferred_time"
                    value={formData.preferred_time}
                    onChange={handleChange}
                    placeholder="Morning / Afternoon / Evening"
                    style={{
                      width: '100%',
                      padding: '8px 12px',
                      fontSize: '14px',
                      border: '1px solid var(--realwiz-gray-300)',
                      borderRadius: '4px',
                    }}
                  />
                </div>
              </div>

              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', marginBottom: '4px', fontWeight: 'bold' }}>
                  Additional Notes
                </label>
                <textarea
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  rows={4}
                  placeholder="Any specific details, concerns, or questions for the provider..."
                  style={{
                    width: '100%',
                    padding: '8px 12px',
                    fontSize: '14px',
                    border: '1px solid var(--realwiz-gray-300)',
                    borderRadius: '4px',
                    fontFamily: 'inherit',
                  }}
                />
              </div>
            </div>

            {/* Submit */}
            <div className="alert alert-info" style={{ marginBottom: '16px' }}>
              <p style={{ margin: 0, fontSize: '14px' }}>
                By submitting this request, you consent to sharing your contact information with the service provider.
                The provider will contact you directly to confirm availability and pricing.
              </p>
            </div>

            <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
              <button
                type="button"
                className="btn-secondary"
                onClick={() => navigate(`/services/${id}`)}
                disabled={submitting}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="btn-primary"
                disabled={submitting}
              >
                {submitting ? 'Submitting...' : 'Submit Request'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

export default ServiceRequestForm;
