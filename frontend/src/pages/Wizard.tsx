import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProject, createProperty, updateProperty, createDisclosure, updateDisclosure, getProperty, getDisclosure } from '../api';
import type { Project, Property, Disclosure } from '../types';

const STEPS = [
  { id: 1, label: 'Property Details' },
  { id: 2, label: 'Disclosure' },
  { id: 3, label: 'Review' },
];

function Wizard() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState(1);
  const [project, setProject] = useState<Project | null>(null);
  const [property, setProperty] = useState<Partial<Property>>({
    year_built: 0,
    legal_description: '',
    tax_id: '',
    is_homestead: false,
    has_hoa: false,
    has_survey: false,
  });
  const [disclosure, setDisclosure] = useState<Partial<Disclosure>>({
    flooding_history: false,
    in_flood_plain: false,
    flood_insurance_required: false,
    foundation_issues: false,
    roof_age: 0,
    roof_issues: false,
    plumbing_issues: false,
    electrical_issues: false,
    hvac_age: 0,
    hvac_issues: false,
    lead_based_paint: false,
    insurance_claims_history: false,
  });

  useEffect(() => {
    if (id) {
      loadProjectData();
    }
  }, [id]);

  const loadProjectData = async () => {
    try {
      const projectRes = await getProject(id!);
      setProject(projectRes.data);

      // Try to load existing property and disclosure
      try {
        const propertyRes = await getProperty(id!);
        setProperty(propertyRes.data);
      } catch (e) {
        // Property doesn't exist yet
      }

      try {
        const disclosureRes = await getDisclosure(id!);
        setDisclosure(disclosureRes.data);
      } catch (e) {
        // Disclosure doesn't exist yet
      }
    } catch (error) {
      console.error('Failed to load project:', error);
    }
  };

  const handlePropertySubmit = async () => {
    try {
      if (property.id) {
        await updateProperty(id!, property);
      } else {
        await createProperty(id!, property);
      }
      setCurrentStep(2);
    } catch (error) {
      console.error('Failed to save property:', error);
      alert('Failed to save property details');
    }
  };

  const handleDisclosureSubmit = async () => {
    try {
      // Auto-set lead based paint if property built before 1978
      const updatedDisclosure = {
        ...disclosure,
        lead_based_paint: (property.year_built ?? 0) > 0 && (property.year_built ?? 0) < 1978,
      };

      if (disclosure.id) {
        await updateDisclosure(id!, updatedDisclosure);
      } else {
        await createDisclosure(id!, updatedDisclosure);
      }
      setCurrentStep(3);
    } catch (error) {
      console.error('Failed to save disclosure:', error);
      alert('Failed to save disclosure');
    }
  };

  const handleComplete = () => {
    navigate(`/project/${id}`);
  };

  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', padding: '40px 20px' }}>
      {/* Header */}
      <div style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '28px', marginBottom: '8px' }}>
          {project?.property_address || 'Property Setup'}
        </h1>
        <p style={{ color: 'var(--realwiz-gray-600)' }}>
          Complete the following steps to prepare your transaction
        </p>
      </div>

      {/* Progress Indicator */}
      <div className="step-indicator">
        {STEPS.map((step, index) => (
          <div key={step.id} className={`step ${currentStep === step.id ? 'active' : ''} ${currentStep > step.id ? 'completed' : ''}`}>
            {index > 0 && <div className="step-line"></div>}
            <div className="step-circle">{step.id}</div>
            <div className="step-label">{step.label}</div>
          </div>
        ))}
      </div>

      {/* Step Content */}
      {currentStep === 1 && (
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Property Details</h2>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Year Built *</label>
            <input
              type="number"
              className="input-field"
              value={property.year_built || ''}
              onChange={e => setProperty({ ...property, year_built: parseInt(e.target.value) })}
              placeholder="1995"
              required
            />
            <p className="help-text">Pre-1978 properties require Lead-Based Paint Addendum (OP-L)</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Legal Description</label>
            <textarea
              className="input-field"
              value={property.legal_description}
              onChange={e => setProperty({ ...property, legal_description: e.target.value })}
              placeholder="Lot 15, Block 3, Example Subdivision..."
              rows={3}
            />
            <p className="help-text">As shown on the deed or survey</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Tax ID / Parcel Number</label>
            <input
              type="text"
              className="input-field"
              value={property.tax_id}
              onChange={e => setProperty({ ...property, tax_id: e.target.value })}
              placeholder="1234-5678-9012"
            />
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={property.is_homestead}
                onChange={e => setProperty({ ...property, is_homestead: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>This is my homestead (primary residence)</span>
            </label>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={property.has_hoa}
                onChange={e => setProperty({ ...property, has_hoa: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Property is in an HOA/POA</span>
            </label>
          </div>

          {property.has_hoa && (
            <div style={{ paddingLeft: '26px', marginBottom: '24px' }}>
              <div style={{ marginBottom: '16px' }}>
                <label className="label">HOA Name</label>
                <input
                  type="text"
                  className="input-field"
                  value={property.hoa_name || ''}
                  onChange={e => setProperty({ ...property, hoa_name: e.target.value })}
                  placeholder="Example HOA"
                />
              </div>
              <div style={{ marginBottom: '16px' }}>
                <label className="label">HOA Fee Amount</label>
                <input
                  type="number"
                  className="input-field"
                  value={property.hoa_fee || ''}
                  onChange={e => setProperty({ ...property, hoa_fee: parseFloat(e.target.value) })}
                  placeholder="250.00"
                />
              </div>
              <div>
                <label className="label">Fee Frequency</label>
                <select
                  className="input-field"
                  value={property.hoa_frequency || 'monthly'}
                  onChange={e => setProperty({ ...property, hoa_frequency: e.target.value as any })}
                >
                  <option value="monthly">Monthly</option>
                  <option value="quarterly">Quarterly</option>
                  <option value="annual">Annual</option>
                </select>
              </div>
            </div>
          )}

          <div style={{ marginBottom: '24px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={property.has_survey}
                onChange={e => setProperty({ ...property, has_survey: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>I have an existing survey</span>
            </label>
            <p className="help-text" style={{ marginLeft: '26px' }}>
              If yes, you may complete a T-47 Affidavit. If no, buyer typically obtains new survey.
            </p>
          </div>

          <div style={{ display: 'flex', gap: '12px', marginTop: '32px' }}>
            <button className="btn-secondary" onClick={() => navigate(`/project/${id}`)}>
              Save & Exit
            </button>
            <button className="btn-primary" onClick={handlePropertySubmit}>
              Continue
            </button>
          </div>
        </div>
      )}

      {currentStep === 2 && (
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '8px' }}>Seller's Disclosure</h2>
          <p style={{ color: 'var(--realwiz-gray-600)', marginBottom: '24px' }}>
            Texas law requires sellers to disclose known material defects (Property Code §5.008)
          </p>

          <div className="alert alert-warning" style={{ marginBottom: '24px' }}>
            <div>
              <strong>Answer honestly and completely</strong>
              <p style={{ fontSize: '13px', marginTop: '4px' }}>
                Failure to disclose known defects can result in legal liability after closing.
              </p>
            </div>
          </div>

          {/* Flooding */}
          <h3 style={{ fontSize: '16px', marginBottom: '16px', marginTop: '24px' }}>Water & Flooding</h3>

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.flooding_history}
                onChange={e => setDisclosure({ ...disclosure, flooding_history: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Property has experienced flooding</span>
            </label>
          </div>

          {disclosure.flooding_history && (
            <div style={{ marginBottom: '16px', paddingLeft: '26px' }}>
              <label className="label">Flooding Details</label>
              <textarea
                className="input-field"
                value={disclosure.flooding_details || ''}
                onChange={e => setDisclosure({ ...disclosure, flooding_details: e.target.value })}
                placeholder="Describe the flooding events and any remediation..."
                rows={3}
              />
            </div>
          )}

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.in_flood_plain}
                onChange={e => setDisclosure({ ...disclosure, in_flood_plain: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Property is in a flood plain</span>
            </label>
          </div>

          {/* Foundation */}
          <h3 style={{ fontSize: '16px', marginBottom: '16px', marginTop: '24px' }}>Structure</h3>

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.foundation_issues}
                onChange={e => setDisclosure({ ...disclosure, foundation_issues: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Foundation has issues or repairs</span>
            </label>
          </div>

          {disclosure.foundation_issues && (
            <div style={{ marginBottom: '16px', paddingLeft: '26px' }}>
              <label className="label">Foundation Details</label>
              <textarea
                className="input-field"
                value={disclosure.foundation_details || ''}
                onChange={e => setDisclosure({ ...disclosure, foundation_details: e.target.value })}
                placeholder="Describe foundation issues and any repairs made..."
                rows={3}
              />
            </div>
          )}

          {/* Systems */}
          <h3 style={{ fontSize: '16px', marginBottom: '16px', marginTop: '24px' }}>Systems</h3>

          <div style={{ marginBottom: '16px' }}>
            <label className="label">Roof Age (years)</label>
            <input
              type="number"
              className="input-field"
              value={disclosure.roof_age || ''}
              onChange={e => setDisclosure({ ...disclosure, roof_age: parseInt(e.target.value) })}
              placeholder="10"
            />
          </div>

          <div style={{ marginBottom: '16px' }}>
            <label className="label">HVAC Age (years)</label>
            <input
              type="number"
              className="input-field"
              value={disclosure.hvac_age || ''}
              onChange={e => setDisclosure({ ...disclosure, hvac_age: parseInt(e.target.value) })}
              placeholder="5"
            />
          </div>

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.plumbing_issues}
                onChange={e => setDisclosure({ ...disclosure, plumbing_issues: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Known plumbing issues</span>
            </label>
          </div>

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.electrical_issues}
                onChange={e => setDisclosure({ ...disclosure, electrical_issues: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Known electrical issues</span>
            </label>
          </div>

          {/* Insurance Claims */}
          <h3 style={{ fontSize: '16px', marginBottom: '16px', marginTop: '24px' }}>Insurance</h3>

          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <input
                type="checkbox"
                checked={disclosure.insurance_claims_history}
                onChange={e => setDisclosure({ ...disclosure, insurance_claims_history: e.target.checked })}
                style={{ width: '18px', height: '18px' }}
              />
              <span className="label" style={{ marginBottom: 0 }}>Property has insurance claims history</span>
            </label>
          </div>

          {/* Other */}
          <h3 style={{ fontSize: '16px', marginBottom: '16px', marginTop: '24px' }}>Other Defects</h3>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Other Material Defects</label>
            <textarea
              className="input-field"
              value={disclosure.other_material_defects || ''}
              onChange={e => setDisclosure({ ...disclosure, other_material_defects: e.target.value })}
              placeholder="Describe any other known material defects..."
              rows={4}
            />
          </div>

          <div style={{ display: 'flex', gap: '12px', marginTop: '32px' }}>
            <button className="btn-secondary" onClick={() => setCurrentStep(1)}>
              Back
            </button>
            <button className="btn-primary" onClick={handleDisclosureSubmit}>
              Continue
            </button>
          </div>
        </div>
      )}

      {currentStep === 3 && (
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '24px' }}>Setup Complete!</h2>

          <div className="alert alert-success" style={{ marginBottom: '24px' }}>
            <div>
              <strong>Property profile created</strong>
              <p style={{ fontSize: '14px', marginTop: '4px' }}>
                Your property details and disclosure information have been saved.
              </p>
            </div>
          </div>

          <h3 style={{ fontSize: '18px', marginBottom: '16px' }}>Next Steps:</h3>
          <ul style={{ paddingLeft: '20px', lineHeight: '1.8' }}>
            <li>Review your property dashboard</li>
            <li>When you receive an offer, enter contract terms to generate deadlines</li>
            <li>Download and complete required TREC forms</li>
            <li>Track all important deadlines through closing</li>
          </ul>

          <div style={{ marginTop: '32px' }}>
            <button className="btn-primary" onClick={handleComplete}>
              Go to Dashboard
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Wizard;
