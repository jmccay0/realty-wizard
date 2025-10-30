import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProjectSummary } from '../api';
import type { ProjectSummary, Deadline } from '../types';

function Dashboard() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [summary, setSummary] = useState<ProjectSummary | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      loadSummary();
    }
  }, [id]);

  const loadSummary = async () => {
    try {
      const response = await getProjectSummary(id!);
      setSummary(response.data);
    } catch (error) {
      console.error('Failed to load project summary:', error);
    } finally {
      setLoading(false);
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

  const upcomingDeadlines = summary.deadlines
    ?.filter(d => !d.completed && new Date(d.due_date) >= new Date())
    .sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime())
    .slice(0, 5) || [];

  const completedDeadlines = summary.deadlines?.filter(d => d.completed) || [];

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const getDaysUntil = (dateString: string) => {
    const days = Math.ceil((new Date(dateString).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24));
    if (days === 0) return 'Today';
    if (days === 1) return 'Tomorrow';
    if (days < 0) return `${Math.abs(days)} days ago`;
    return `in ${days} days`;
  };

  const getPriorityBadge = (priority: string) => {
    switch (priority) {
      case 'critical': return 'badge-critical';
      case 'high': return 'badge-warning';
      default: return 'badge-info';
    }
  };

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '40px 20px' }}>
      {/* Header */}
      <div style={{ marginBottom: '32px' }}>
        <button
          onClick={() => navigate('/')}
          style={{
            background: 'none',
            border: 'none',
            color: 'var(--realwiz-blue)',
            cursor: 'pointer',
            marginBottom: '16px',
            fontSize: '14px',
          }}
        >
          ← Back to all transactions
        </button>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '8px' }}>
          <h1 style={{ fontSize: '28px', margin: 0 }}>
            {summary.project.property_address}
          </h1>
          <span style={{
            display: 'inline-block',
            padding: '6px 16px',
            borderRadius: '20px',
            fontSize: '14px',
            fontWeight: 600,
            backgroundColor: summary.project.user_role === 'buyer' ? 'var(--realwiz-green-50)' : 'var(--realwiz-blue-50)',
            color: summary.project.user_role === 'buyer' ? 'var(--realwiz-green)' : 'var(--realwiz-blue)',
            border: `2px solid ${summary.project.user_role === 'buyer' ? 'var(--realwiz-green)' : 'var(--realwiz-blue)'}`,
          }}>
            {summary.project.user_role === 'buyer' ? '🏠 Buyer' : '📋 Seller'}
          </span>
        </div>
        <div style={{ display: 'flex', gap: '16px', alignItems: 'center' }}>
          <span className={`badge badge-${
            summary.project.status === 'setup' ? 'info' :
            summary.project.status === 'listing_prep' ? 'warning' :
            summary.project.status === 'under_contract' ? 'warning' :
            'success'
          }`}>
            {summary.project.status.replace('_', ' ')}
          </span>
          <span style={{ color: 'var(--realwiz-gray-600)', fontSize: '14px' }}>
            Sellers: {summary.project.seller_names.join(', ')}
          </span>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px', marginBottom: '32px' }}>
        {/* Quick Actions */}
        <div className="card">
          <h2 style={{ fontSize: '18px', marginBottom: '16px' }}>Quick Actions</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {!summary.property && (
              <button
                className="btn-primary"
                onClick={() => navigate(`/project/${id}/wizard`)}
              >
                Complete Property Setup
              </button>
            )}
            {summary.property && !summary.contract && (
              <button
                className="btn-primary"
                onClick={() => navigate(`/project/${id}/contract/new`)}
              >
                Enter Contract Terms
              </button>
            )}
            <button className="btn-secondary" onClick={() => navigate(`/project/${id}/wizard`)}>
              Edit Property Details
            </button>
            <button className="btn-secondary" onClick={() => navigate(`/project/${id}/documents`)}>
              View Documents
            </button>
          </div>
        </div>

        {/* Property Info */}
        <div className="card">
          <h2 style={{ fontSize: '18px', marginBottom: '16px' }}>Property Information</h2>
          {summary.property ? (
            <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
              <div><strong>Year Built:</strong> {summary.property.year_built}</div>
              <div><strong>Homestead:</strong> {summary.property.is_homestead ? 'Yes' : 'No'}</div>
              <div><strong>HOA:</strong> {summary.property.has_hoa ? `Yes - ${summary.property.hoa_name}` : 'No'}</div>
              <div><strong>Survey:</strong> {summary.property.has_survey ? 'Have existing survey' : 'Buyer to obtain'}</div>
              {summary.property.year_built < 1978 && (
                <div className="alert alert-warning" style={{ marginTop: '12px', padding: '12px' }}>
                  <strong>Lead-Based Paint Disclosure Required</strong>
                  <p style={{ fontSize: '12px', marginTop: '4px' }}>
                    Pre-1978 property requires OP-L addendum
                  </p>
                </div>
              )}
            </div>
          ) : (
            <p style={{ color: 'var(--realwiz-gray-600)' }}>
              No property details yet. Complete the setup wizard to add information.
            </p>
          )}
        </div>
      </div>

      {/* Deadlines Section */}
      {summary.contract && (
        <div className="card">
          <h2 style={{ fontSize: '20px', marginBottom: '16px' }}>Upcoming Deadlines</h2>

          {upcomingDeadlines.length === 0 ? (
            <p style={{ color: 'var(--realwiz-gray-600)' }}>No upcoming deadlines</p>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              {upcomingDeadlines.map((deadline: Deadline) => (
                <div
                  key={deadline.id}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'start',
                    padding: '16px',
                    background: 'var(--realwiz-gray-50)',
                    borderRadius: '6px',
                    borderLeft: `4px solid ${
                      deadline.priority === 'critical' ? 'var(--realwiz-red)' :
                      deadline.priority === 'high' ? 'var(--realwiz-orange)' :
                      'var(--realwiz-blue)'
                    }`,
                  }}
                >
                  <div style={{ flex: 1 }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                      <strong>{deadline.description}</strong>
                      <span className={`badge ${getPriorityBadge(deadline.priority)}`}>
                        {deadline.priority}
                      </span>
                    </div>
                    {deadline.notes && (
                      <p style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)', marginTop: '4px' }}>
                        {deadline.notes}
                      </p>
                    )}
                  </div>
                  <div style={{ textAlign: 'right', minWidth: '120px' }}>
                    <div style={{ fontWeight: 600 }}>{formatDate(deadline.due_date)}</div>
                    <div style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)' }}>
                      {getDaysUntil(deadline.due_date)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {completedDeadlines.length > 0 && (
            <div style={{ marginTop: '24px' }}>
              <h3 style={{ fontSize: '16px', marginBottom: '12px', color: 'var(--realwiz-gray-600)' }}>
                Completed ({completedDeadlines.length})
              </h3>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                {completedDeadlines.map((deadline: Deadline) => (
                  <div
                    key={deadline.id}
                    style={{
                      display: 'flex',
                      justifyContent: 'space-between',
                      padding: '12px',
                      background: 'var(--realwiz-gray-100)',
                      borderRadius: '6px',
                      opacity: 0.7,
                    }}
                  >
                    <div style={{ textDecoration: 'line-through' }}>{deadline.description}</div>
                    <div style={{ fontSize: '13px', color: 'var(--realwiz-gray-600)' }}>
                      {formatDate(deadline.due_date)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Contract Summary */}
      {summary.contract && (
        <div className="card" style={{ marginTop: '24px' }}>
          <h2 style={{ fontSize: '20px', marginBottom: '16px' }}>Contract Summary</h2>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', fontSize: '14px' }}>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Sales Price</div>
              <div style={{ fontSize: '18px', fontWeight: 600 }}>
                ${summary.contract.sales_price.toLocaleString()}
              </div>
            </div>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Closing Date</div>
              <div style={{ fontSize: '18px', fontWeight: 600 }}>
                {formatDate(summary.contract.closing_date)}
              </div>
            </div>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Option Fee</div>
              <div>${summary.contract.option_fee.toLocaleString()}</div>
            </div>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Option Period</div>
              <div>{summary.contract.option_period_days} days</div>
            </div>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Earnest Money</div>
              <div>${summary.contract.earnest_money.toLocaleString()}</div>
            </div>
            <div>
              <div style={{ color: 'var(--realwiz-gray-600)', fontSize: '12px' }}>Financing</div>
              <div>{summary.contract.buyer_financing ? 'Yes' : 'Cash'}</div>
            </div>
          </div>
        </div>
      )}

      {/* Disclosure Summary */}
      {summary.disclosure && (
        <div className="card" style={{ marginTop: '24px' }}>
          <h2 style={{ fontSize: '20px', marginBottom: '16px' }}>Disclosure Status</h2>
          <div className="alert alert-success">
            <div>
              <strong>Seller's Disclosure completed</strong>
              <p style={{ fontSize: '13px', marginTop: '4px' }}>
                Last updated {formatDate(summary.disclosure.updated_at)}
              </p>
            </div>
          </div>
          {summary.disclosure.lead_based_paint && (
            <div className="alert alert-warning" style={{ marginTop: '12px' }}>
              <strong>Lead-Based Paint Disclosure Required</strong>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default Dashboard;
