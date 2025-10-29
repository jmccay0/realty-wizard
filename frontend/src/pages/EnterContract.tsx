import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProject, createContract } from '../api';
import type { Project, ContractTerms } from '../types';

const EnterContract: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    effectiveDate: '',
    salesPrice: '',
    closingDate: '',
    optionFee: '',
    optionPeriodDays: '10',
    earnestMoney: '',
    titleCommitmentDays: '20',
    buyerFinancing: false,
    sellerStaysPostClose: false,
    sellerLeaseDays: '',
    excludedItems: '',
    includedPersonalItems: '',
  });

  useEffect(() => {
    const loadProject = async () => {
      if (!id) return;
      try {
        const response = await getProject(id);
        setProject(response.data);
      } catch (err) {
        setError('Failed to load project');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    loadProject();
  }, [id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    const checked = (e.target as HTMLInputElement).checked;

    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;

    setSubmitting(true);
    setError(null);

    try {
      // Parse string arrays from comma-separated input
      const excludedItems = formData.excludedItems
        ? formData.excludedItems.split(',').map(item => item.trim()).filter(Boolean)
        : [];

      const includedPersonalItems = formData.includedPersonalItems
        ? formData.includedPersonalItems.split(',').map(item => item.trim()).filter(Boolean)
        : [];

      const contractData: Partial<ContractTerms> = {
        project_id: id,
        effective_date: new Date(formData.effectiveDate).toISOString(),
        sales_price: parseFloat(formData.salesPrice),
        closing_date: new Date(formData.closingDate).toISOString(),
        option_fee: parseFloat(formData.optionFee),
        option_period_days: parseInt(formData.optionPeriodDays),
        earnest_money: parseFloat(formData.earnestMoney),
        title_commitment_days: parseInt(formData.titleCommitmentDays),
        buyer_financing: formData.buyerFinancing,
        seller_stays_post_close: formData.sellerStaysPostClose,
        excluded_items: excludedItems,
        included_personal_items: includedPersonalItems,
      };

      // Only include seller_lease_days if seller stays post-close
      if (formData.sellerStaysPostClose && formData.sellerLeaseDays) {
        contractData.seller_lease_days = parseInt(formData.sellerLeaseDays);
      }

      await createContract(id, contractData);

      // Navigate back to dashboard
      navigate(`/project/${id}`);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create contract');
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading project...</div>;
  }

  if (!project) {
    return <div className="error">Project not found</div>;
  }

  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', padding: '40px 20px' }}>
      {/* Header */}
      <div style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '28px', marginBottom: '8px' }}>Enter Contract Terms</h1>
        <p style={{ color: 'var(--realwiz-gray-600)' }}>{project.property_address}</p>
      </div>

      {error && (
        <div className="alert alert-error" style={{ marginBottom: '24px' }}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Key Dates & Amounts</h2>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Effective Date <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="date"
              className="input-field"
              name="effectiveDate"
              value={formData.effectiveDate}
              onChange={handleChange}
              required
            />
            <p className="help-text">The date the contract becomes effective (typically after final acceptance)</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Sales Price <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="number"
              className="input-field"
              name="salesPrice"
              value={formData.salesPrice}
              onChange={handleChange}
              step="0.01"
              min="0"
              placeholder="450000"
              required
            />
            <p className="help-text">Total purchase price in dollars</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Closing Date <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="date"
              className="input-field"
              name="closingDate"
              value={formData.closingDate}
              onChange={handleChange}
              required
            />
            <p className="help-text">Target closing date (can be extended by mutual agreement)</p>
          </div>
        </div>

        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Option Period & Earnest Money</h2>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Option Fee <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="number"
              className="input-field"
              name="optionFee"
              value={formData.optionFee}
              onChange={handleChange}
              step="0.01"
              min="0"
              placeholder="500"
              required
            />
            <p className="help-text">Fee paid for unrestricted right to terminate contract during option period</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Option Period (Days) <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="number"
              className="input-field"
              name="optionPeriodDays"
              value={formData.optionPeriodDays}
              onChange={handleChange}
              min="0"
              max="365"
              placeholder="10"
              required
            />
            <p className="help-text">Number of days buyer has to terminate for any reason (typically 7-10 days)</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Earnest Money <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="number"
              className="input-field"
              name="earnestMoney"
              value={formData.earnestMoney}
              onChange={handleChange}
              step="0.01"
              min="0"
              placeholder="5000"
              required
            />
            <p className="help-text">Good faith deposit held in escrow (typically 1-2% of sales price)</p>
          </div>
        </div>

        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Title & Financing</h2>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">
              Title Commitment Days <span style={{ color: 'var(--realwiz-red)' }}>*</span>
            </label>
            <input
              type="number"
              className="input-field"
              name="titleCommitmentDays"
              value={formData.titleCommitmentDays}
              onChange={handleChange}
              min="1"
              max="90"
              placeholder="20"
              required
            />
            <p className="help-text">Days within which seller must furnish title commitment (typically 20 days)</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                name="buyerFinancing"
                checked={formData.buyerFinancing}
                onChange={handleChange}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Buyer will use financing (Third Party Financing Addendum)</span>
            </label>
            <p className="help-text" style={{ marginLeft: '26px' }}>Check if buyer needs a loan to purchase the property</p>
          </div>
        </div>

        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Seller Occupancy After Closing</h2>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                name="sellerStaysPostClose"
                checked={formData.sellerStaysPostClose}
                onChange={handleChange}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Seller will remain in property after closing</span>
            </label>
            <p className="help-text" style={{ marginLeft: '26px' }}>Check if using Seller's Temporary Residential Lease</p>
          </div>

          {formData.sellerStaysPostClose && (
            <div style={{ paddingLeft: '26px', marginBottom: '24px' }}>
              <label className="label">Number of Days Seller Will Lease Back</label>
              <input
                type="number"
                className="input-field"
                name="sellerLeaseDays"
                value={formData.sellerLeaseDays}
                onChange={handleChange}
                min="1"
                max="90"
                placeholder="30"
              />
              <p className="help-text">Maximum 90 days per standard TREC lease</p>
            </div>
          )}
        </div>

        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Items & Accessories (Optional)</h2>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Excluded Items</label>
            <textarea
              className="input-field"
              name="excludedItems"
              value={formData.excludedItems}
              onChange={handleChange}
              rows={3}
              placeholder="Enter items separated by commas (e.g., chandelier in dining room, custom curtains, antique mirror)"
            />
            <p className="help-text">Items that would normally convey but are being excluded</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Included Personal Property</label>
            <textarea
              className="input-field"
              name="includedPersonalItems"
              value={formData.includedPersonalItems}
              onChange={handleChange}
              rows={3}
              placeholder="Enter items separated by commas (e.g., refrigerator, washer, dryer, patio furniture)"
            />
            <p className="help-text">Personal property that does not typically convey but is being included</p>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '12px', marginTop: '32px' }}>
          <button
            type="button"
            onClick={() => navigate(`/project/${id}`)}
            className="btn-secondary"
            disabled={submitting}
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn-primary"
            disabled={submitting}
          >
            {submitting ? 'Creating Contract...' : 'Create Contract & Generate Deadlines'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default EnterContract;
