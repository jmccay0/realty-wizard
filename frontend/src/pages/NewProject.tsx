import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createProject } from '../api';
import { useAuth } from '../contexts/AuthContext';

function NewProject() {
  const navigate = useNavigate();
  const { user } = useAuth();

  const [formData, setFormData] = useState<{
    property_address: string;
    user_role: 'buyer' | 'seller' | '';
    seller_names: string[];
    seller_email: string;
    seller_phone: string;
    buyer_names: string[];
    buyer_email: string;
    buyer_phone: string;
    has_agent: boolean;
    agent_name: string;
  }>({
    property_address: '',
    user_role: '', // 'buyer' or 'seller'
    seller_names: [user?.name || ''],
    seller_email: user?.email || '',
    seller_phone: user?.phone || '',
    buyer_names: [user?.name || ''],
    buyer_email: user?.email || '',
    buyer_phone: user?.phone || '',
    has_agent: false,
    agent_name: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await createProject({
        ...formData,
        user_role: formData.user_role as 'buyer' | 'seller',
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

  const addBuyer = () => {
    setFormData({
      ...formData,
      buyer_names: [...formData.buyer_names, ''],
    });
  };

  const updateBuyer = (index: number, value: string) => {
    const newNames = [...formData.buyer_names];
    newNames[index] = value;
    setFormData({ ...formData, buyer_names: newNames });
  };

  const removeBuyer = (index: number) => {
    const newNames = formData.buyer_names.filter((_, i) => i !== index);
    setFormData({ ...formData, buyer_names: newNames });
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
          {/* Role Selection */}
          <div style={{ marginBottom: '32px', paddingBottom: '24px', borderBottom: '1px solid var(--realwiz-gray-300)' }}>
            <label className="label">What is your role in this transaction? *</label>
            <div style={{ display: 'flex', gap: '16px', marginTop: '12px' }}>
              <label style={{
                flex: 1,
                display: 'flex',
                alignItems: 'center',
                gap: '12px',
                padding: '16px',
                border: `2px solid ${formData.user_role === 'seller' ? 'var(--realwiz-blue)' : 'var(--realwiz-gray-300)'}`,
                borderRadius: '8px',
                cursor: 'pointer',
                backgroundColor: formData.user_role === 'seller' ? 'var(--realwiz-blue-50)' : 'white',
              }}>
                <input
                  type="radio"
                  name="user_role"
                  value="seller"
                  checked={formData.user_role === 'seller'}
                  onChange={e => setFormData({ ...formData, user_role: e.target.value as 'seller', buyer_names: [''], buyer_email: '', buyer_phone: '' })}
                  style={{ width: '20px', height: '20px', cursor: 'pointer' }}
                  required
                />
                <div>
                  <div style={{ fontWeight: 600, marginBottom: '4px' }}>I'm the Seller</div>
                  <div style={{ fontSize: '14px', color: 'var(--realwiz-gray-600)' }}>I'm selling the property</div>
                </div>
              </label>

              <label style={{
                flex: 1,
                display: 'flex',
                alignItems: 'center',
                gap: '12px',
                padding: '16px',
                border: `2px solid ${formData.user_role === 'buyer' ? 'var(--realwiz-blue)' : 'var(--realwiz-gray-300)'}`,
                borderRadius: '8px',
                cursor: 'pointer',
                backgroundColor: formData.user_role === 'buyer' ? 'var(--realwiz-blue-50)' : 'white',
              }}>
                <input
                  type="radio"
                  name="user_role"
                  value="buyer"
                  checked={formData.user_role === 'buyer'}
                  onChange={e => setFormData({ ...formData, user_role: e.target.value as 'buyer', seller_names: [''], seller_email: '', seller_phone: '' })}
                  style={{ width: '20px', height: '20px', cursor: 'pointer' }}
                  required
                />
                <div>
                  <div style={{ fontWeight: 600, marginBottom: '4px' }}>I'm the Buyer</div>
                  <div style={{ fontSize: '14px', color: 'var(--realwiz-gray-600)' }}>I'm purchasing the property</div>
                </div>
              </label>
            </div>
          </div>

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

          {/* Seller Information (shown when user is seller) */}
          {formData.user_role === 'seller' && (
            <>
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
            </>
          )}

          {/* Buyer Information (shown when user is buyer) */}
          {formData.user_role === 'buyer' && (
            <>
              <div style={{ marginBottom: '24px' }}>
                <label className="label">Buyer Name(s) *</label>
                {formData.buyer_names.map((name, index) => (
                  <div key={index} style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
                    <input
                      type="text"
                      className="input-field"
                      value={name}
                      onChange={e => updateBuyer(index, e.target.value)}
                      placeholder="Full legal name"
                      required
                    />
                    {formData.buyer_names.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeBuyer(index)}
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
                  onClick={addBuyer}
                  style={{ marginTop: '8px' }}
                >
                  + Add Another Buyer
                </button>
                <p className="help-text">Legal names as they will appear on the contract</p>
              </div>

              <div style={{ marginBottom: '24px' }}>
                <label className="label">Email *</label>
                <input
                  type="email"
                  className="input-field"
                  value={formData.buyer_email}
                  onChange={e => setFormData({ ...formData, buyer_email: e.target.value })}
                  placeholder="buyer@example.com"
                  required
                />
              </div>

              <div style={{ marginBottom: '24px' }}>
                <label className="label">Phone</label>
                <input
                  type="tel"
                  className="input-field"
                  value={formData.buyer_phone}
                  onChange={e => setFormData({ ...formData, buyer_phone: e.target.value })}
                  placeholder="(512) 555-0123"
                />
              </div>
            </>
          )}

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
