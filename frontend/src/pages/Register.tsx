import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

function Register() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
    phone: '',
    defaultRole: 'seller',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validate password match
    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    // Validate password length
    if (formData.password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }

    setLoading(true);

    try {
      await register(
        formData.email,
        formData.password,
        formData.name,
        formData.phone,
        formData.defaultRole
      );
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '80px auto', padding: '0 20px' }}>
      <div className="card">
        <h1 style={{ fontSize: '28px', marginBottom: '8px', textAlign: 'center' }}>
          Create Account
        </h1>
        <p style={{ color: 'var(--realwiz-gray-600)', textAlign: 'center', marginBottom: '32px' }}>
          Sign up to get started
        </p>

        {error && (
          <div className="alert alert-error" style={{ marginBottom: '24px' }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '20px' }}>
            <label className="label">Full Name</label>
            <input
              type="text"
              name="name"
              className="input-field"
              value={formData.name}
              onChange={handleChange}
              placeholder="John Doe"
              required
              autoFocus
            />
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label className="label">Email</label>
            <input
              type="email"
              name="email"
              className="input-field"
              value={formData.email}
              onChange={handleChange}
              placeholder="you@example.com"
              required
            />
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label className="label">Phone</label>
            <input
              type="tel"
              name="phone"
              className="input-field"
              value={formData.phone}
              onChange={handleChange}
              placeholder="555-1234"
              required
            />
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label className="label">I am a...</label>
            <select
              name="defaultRole"
              className="input-field"
              value={formData.defaultRole}
              onChange={handleChange}
              required
            >
              <option value="seller">Seller</option>
              <option value="buyer">Buyer</option>
            </select>
            <p className="help-text">You can participate in transactions with different roles</p>
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label className="label">Password</label>
            <input
              type="password"
              name="password"
              className="input-field"
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
              required
            />
            <p className="help-text">Minimum 8 characters</p>
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label className="label">Confirm Password</label>
            <input
              type="password"
              name="confirmPassword"
              className="input-field"
              value={formData.confirmPassword}
              onChange={handleChange}
              placeholder="••••••••"
              required
            />
          </div>

          <button
            type="submit"
            className="btn-primary"
            disabled={loading}
            style={{ width: '100%' }}
          >
            {loading ? 'Creating account...' : 'Create Account'}
          </button>
        </form>

        <p style={{ textAlign: 'center', marginTop: '24px', fontSize: '14px', color: 'var(--realwiz-gray-600)' }}>
          Already have an account?{' '}
          <Link to="/login" style={{ color: 'var(--realwiz-blue)', textDecoration: 'none' }}>
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}

export default Register;
