import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProjectSummary, generateDocument, downloadDocument } from '../api';
import type { ProjectSummary, Document } from '../types';

const ViewDocuments: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [summary, setSummary] = useState<ProjectSummary | null>(null);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, [id]);

  const loadData = async () => {
    if (!id) return;
    try {
      const response = await getProjectSummary(id);
      setSummary(response.data);
      setDocuments(response.data.documents || []);
    } catch (error) {
      console.error('Failed to load project:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (type: string) => {
    if (!id) return;
    setGenerating(type);
    try {
      const response = await generateDocument(id, type);
      setDocuments([...documents, response.data.document]);
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to generate document');
    } finally {
      setGenerating(null);
    }
  };

  const handleDownload = async (doc: Document) => {
    try {
      const response = await downloadDocument(doc.id);
      const blob = new Blob([response.data], { type: 'application/pdf' });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${doc.type}.pdf`;
      link.click();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to download:', error);
      alert('Failed to download document');
    }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <div className="spinner"></div>
      </div>
    );
  }

  if (!summary) {
    return <div>Project not found</div>;
  }

  const isBuyer = summary.project.user_role === 'buyer';

  const availableTemplates = [
    // Seller-only documents
    { type: 'sellers_disclosure', name: "Seller's Disclosure Notice", formNumber: 'TREC OP-H', condition: !!summary.property && !isBuyer, role: 'seller' },

    // Shared documents (both buyer and seller)
    { type: 'lead_paint', name: 'Lead-Based Paint Disclosure', formNumber: 'OP-L', condition: summary.property && summary.property.year_built < 1978, role: 'both' },
    { type: 'residential_contract', name: 'One to Four Family Residential Contract', formNumber: 'TREC 20-18', condition: !!summary.contract, role: 'both' },
    { type: 'third_party_financing', name: 'Third Party Financing Addendum', formNumber: 'TREC 40-9', condition: summary.contract && summary.contract.buyer_financing, role: 'both' },
    { type: 'seller_lease', name: "Seller's Temporary Residential Lease", formNumber: 'TREC 15-6', condition: summary.contract && summary.contract.seller_stays_post_close, role: 'both' },
    { type: 'hoa_addendum', name: 'Addendum for Property Subject to Mandatory Membership in HOA', formNumber: 'TREC 36-10', condition: summary.property && summary.property.has_hoa, role: 'both' },

    // Buyer-specific documents (placeholder for future)
    // { type: 'buyer_offer', name: 'Buyer Offer Letter', formNumber: 'Custom', condition: isBuyer, role: 'buyer' },
  ];

  const applicableTemplates = availableTemplates.filter(t => t.condition);
  const generatedDocs = documents.filter(d => d.status === 'ready');
  const pendingDocs = applicableTemplates.filter(t => !documents.some(d => d.type === t.type));

  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', padding: '40px 20px' }}>
      {/* Header */}
      <div style={{ marginBottom: '32px' }}>
        <button
          onClick={() => navigate(`/project/${id}`)}
          style={{
            background: 'none',
            border: 'none',
            color: 'var(--realwiz-blue)',
            cursor: 'pointer',
            marginBottom: '16px',
            fontSize: '14px',
          }}
        >
          ← Back to Dashboard
        </button>
        <h1 style={{ fontSize: '28px', marginBottom: '8px' }}>Documents</h1>
        <p style={{ color: 'var(--realwiz-gray-600)' }}>{summary.project.property_address}</p>
      </div>

      {/* Generated Documents */}
      {generatedDocs.length > 0 && (
        <div className="card" style={{ marginBottom: '24px' }}>
          <h2 style={{ fontSize: '22px', marginBottom: '16px' }}>Generated Documents</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {generatedDocs.map((doc) => {
              const template = availableTemplates.find(t => t.type === doc.type);
              return (
                <div
                  key={doc.id}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '16px',
                    background: 'var(--realwiz-gray-50)',
                    borderRadius: '6px',
                    border: '1px solid var(--realwiz-gray-200)',
                  }}
                >
                  <div>
                    <div style={{ fontWeight: 600, marginBottom: '4px' }}>
                      {template?.name || doc.type}
                    </div>
                    <div style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)' }}>
                      {doc.form_number} • Generated {doc.generated_at ? new Date(doc.generated_at).toLocaleDateString() : 'recently'}
                    </div>
                  </div>
                  <button
                    className="btn-secondary"
                    onClick={() => handleDownload(doc)}
                    style={{ minWidth: '120px' }}
                  >
                    Download PDF
                  </button>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Available Documents */}
      {pendingDocs.length > 0 && (
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '16px' }}>Available to Generate</h2>
          <p style={{ color: 'var(--realwiz-gray-600)', marginBottom: '16px', fontSize: '14px' }}>
            These documents apply to your transaction based on the information you've entered.
          </p>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {pendingDocs.map((template) => (
              <div
                key={template.type}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '16px',
                  background: 'var(--realwiz-gray-50)',
                  borderRadius: '6px',
                  border: '1px solid var(--realwiz-gray-200)',
                }}
              >
                <div>
                  <div style={{ fontWeight: 600, marginBottom: '4px' }}>{template.name}</div>
                  <div style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)' }}>
                    {template.formNumber}
                  </div>
                </div>
                <button
                  className="btn-primary"
                  onClick={() => handleGenerate(template.type)}
                  disabled={generating === template.type}
                  style={{ minWidth: '120px' }}
                >
                  {generating === template.type ? 'Generating...' : 'Generate'}
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* No Documents Message */}
      {applicableTemplates.length === 0 && (
        <div className="card">
          <h2 style={{ fontSize: '22px', marginBottom: '16px' }}>No Documents Available Yet</h2>
          <p style={{ color: 'var(--realwiz-gray-600)', marginBottom: '16px' }}>
            Complete the property details and contract information to generate transaction documents.
          </p>
          <button className="btn-primary" onClick={() => navigate(`/project/${id}`)}>
            Go to Dashboard
          </button>
        </div>
      )}
    </div>
  );
};

export default ViewDocuments;
