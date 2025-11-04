import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { listProjects } from '../api';
import type { Project } from '../types';

function Home() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      const response = await listProjects();
      setProjects(response.data);
    } catch (error) {
      console.error('Failed to load projects:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '40px 20px' }}>
      {/* Header */}
      <div style={{ marginBottom: '40px' }}>
        <h1 style={{ fontSize: '32px', marginBottom: '8px' }}>Realty Wizard</h1>
        <p style={{ color: 'var(--realwiz-gray-600)', fontSize: '18px' }}>
          Your guide through Texas real estate transactions
        </p>
      </div>

      {/* Legal disclaimer banner */}
      <div className="alert alert-info">
        <div>
          <strong>Information Only - Not Legal Advice</strong>
          <p style={{ marginTop: '4px', fontSize: '14px' }}>
            This tool provides general information about Texas real estate transactions.
            For form use and interpretation, consult a licensed Texas real estate professional or attorney.
          </p>
        </div>
      </div>

      {/* Action Buttons */}
      <div style={{ marginBottom: '32px', display: 'flex', gap: '12px' }}>
        <button
          className="btn-primary"
          onClick={() => navigate('/new')}
        >
          Start New Transaction
        </button>
        <button
          className="btn-secondary"
          onClick={() => navigate('/learn/trec-contract')}
        >
          📚 Learn About Texas Contracts
        </button>
      </div>

      {/* Projects List */}
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
          <div className="spinner"></div>
        </div>
      ) : projects.length === 0 ? (
        <div className="card" style={{ textAlign: 'center', padding: '60px 20px' }}>
          <h3 style={{ marginBottom: '8px' }}>No transactions yet</h3>
          <p style={{ color: 'var(--realwiz-gray-600)' }}>
            Get started by creating your first transaction
          </p>
        </div>
      ) : (
        <div>
          <h2 style={{ fontSize: '24px', marginBottom: '16px' }}>Your Transactions</h2>
          {projects.map(project => (
            <div
              key={project.id}
              className="card"
              style={{ cursor: 'pointer' }}
              onClick={() => navigate(`/project/${project.id}`)}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start' }}>
                <div>
                  <h3 style={{ fontSize: '18px', marginBottom: '8px' }}>
                    {project.property_address}
                  </h3>
                  <p style={{ color: 'var(--realwiz-gray-600)', fontSize: '14px' }}>
                    Sellers: {project.seller_names.join(', ')}
                  </p>
                  <p style={{ color: 'var(--realwiz-gray-600)', fontSize: '13px', marginTop: '4px' }}>
                    Created {new Date(project.created_at).toLocaleDateString()}
                  </p>
                </div>
                <span className={`badge badge-${
                  project.status === 'setup' ? 'info' :
                  project.status === 'listing_prep' ? 'warning' :
                  project.status === 'under_contract' ? 'warning' :
                  'success'
                }`}>
                  {project.status.replace('_', ' ')}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Home;
