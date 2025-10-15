import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createProject } from '../api';

function NewProject() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    property_address: '',
    seller_names: [''],
    seller_email: '',
    seller_phone: '',
    has_agent: false,
    agent_name: '',
    title_company: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await createProject({
        ...formData,
        seller_names: formData.seller_names.filter(name => name.trim() !== ''),
      });

      // Navigate to the wizard for the new project
      navigate(`/project/${response.data.id}/wizard`);
    } catch (error) {
      console.error('Failed to create project:', error);
      alert('Failed to create project. Please try again.');
    }
  };

  const addSeller = () => {
    setFormData({
      ...formData,
      seller_names: [...formData.seller_names, ''],
    });
  };

  const updateSeller = (index: number, value: string) => {
    const newNames = [...formData.seller_names];
    newNames[index] = value;
    setFormData({ ...formData, seller_names: newNames });
  };

  const removeSeller = (index: number) => {
    const newNames = formData.seller_names.filter((_, i) => i !== index);
    setFormData({ ...formData, seller_names: newNames });
  };

  return (
    <div style={{ maxWidth: '700px', margin: '0 auto', padding: '40px 20px' }}>
      <div style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '28px', marginBottom: '8px' }}>Start New Transaction</h1>
        <p style={{ color: 'var(--realwiz-gray-600)' }}>
          Let's get started with some basic information
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="card">
          {/* Property Address */}
          <div style={{ marginBottom: '24px' }}>
            <label className="label">Property Address *</label>
            <input
              type="text"
              className="input-field"
              value={formData.property_address}
              onChange={e => setFormData({ ...formData, property_address: e.target.value })}
              placeholder="123 Main St, Austin, TX 78701"
              required
            />
            <p className="help-text">The address of the property being sold</p>
          </div>

          {/* Seller Names */}
          <div style={{ marginBottom: '24px' }}>
            <label className="label">Seller Name(s) *</label>
            {formData.seller_names.map((name, index) => (
              <div key={index} style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
                <input
                  type="text"
                  className="input-field"
                  value={name}
                  onChange={e => updateSeller(index, e.target.value)}
                  placeholder="Full legal name"
                  required
                />
                {formData.seller_names.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeSeller(index)}
                    style={{
                      padding: '8px 16px',
                      background: 'var(--realwiz-gray-200)',
                      border: 'none',
                      borderRadius: '6px',
                      cursor: 'pointer',
                    }}
                  >
                    Remove
                  </button>
                )}
              </div>
            ))}
            <button
              type="button"
              className="btn-secondary"
              onClick={addSeller}
              style={{ marginTop: '8px' }}
            >
              + Add Another Seller
            </button>
            <p className="help-text">Legal names as they appear on the deed</p>
          </div>

          {/* Contact Info */}
          <div style={{ marginBottom: '24px' }}>
            <label className="label">Email *</label>
            <input
              type="email"
              className="input-field"
              value={formData.seller_email}
              onChange={e => setFormData({ ...formData, seller_email: e.target.value })}
              placeholder="seller@example.com"
              required
            />
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Phone</label>
            <input
              type="tel"
              className="input-field"
              value={formData.seller_phone}
              onChange={e => setFormData({ ...formData, seller_phone: e.target.value })}
              placeholder="(512) 555-0123"
            />
          </div>

          {/* Agent Info */}
          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
              <input
                type="checkbox"
                checked={formData.has_agent}
                onChange={e => setFormData({ ...formData, has_agent: e.target.checked })}
                style={{ width: '18px', height: '18px', cursor: 'pointer' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Working with a real estate agent</span>
            </label>
          </div>

          {formData.has_agent && (
            <div style={{ marginBottom: '24px', paddingLeft: '26px' }}>
              <label className="label">Agent Name</label>
              <input
                type="text"
                className="input-field"
                value={formData.agent_name}
                onChange={e => setFormData({ ...formData, agent_name: e.target.value })}
                placeholder="Agent's full name"
              />
            </div>
          )}

          {/* Title Company */}
          <div style={{ marginBottom: '0' }}>
            <label className="label">Preferred Title Company</label>
            <input
              type="text"
              className="input-field"
              value={formData.title_company}
              onChange={e => setFormData({ ...formData, title_company: e.target.value })}
              placeholder="Optional - can be added later"
            />
            <p className="help-text">The title company that will handle escrow and closing</p>
          </div>
        </div>

        {/* Action Buttons */}
        <div style={{ display: 'flex', gap: '12px', marginTop: '24px' }}>
          <button
            type="button"
            className="btn-secondary"
            onClick={() => navigate('/')}
          >
            Cancel
          </button>
          <button type="submit" className="btn-primary">
            Continue to Property Details
          </button>
        </div>
      </form>
    </div>
  );
}

export default NewProject;
