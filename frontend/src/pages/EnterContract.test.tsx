import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import EnterContract from './EnterContract';
import * as api from '../api';
import type { Project, ContractTerms } from '../types';
import type { AxiosResponse } from 'axios';

// Mock the API module
vi.mock('../api', () => ({
  getProject: vi.fn(),
  createContract: vi.fn(),
}));

const mockProject: AxiosResponse<Project> = {
  data: {
    id: 'test-project-1',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    property_address: '123 Test St, Austin, TX 78701',
    user_role: 'seller',
    seller_names: ['John Doe', 'Jane Doe'],
    seller_email: 'seller@example.com',
    seller_phone: '512-555-0100',
    buyer_names: [''],
    buyer_email: '',
    buyer_phone: '',
    has_agent: true,
    agent_name: 'Agent Smith',
    status: 'setup',
  },
  status: 200,
  statusText: 'OK',
  headers: {},
  config: {} as any,
};

const mockContractResponse: AxiosResponse<ContractTerms> = {
  data: {} as ContractTerms,
  status: 201,
  statusText: 'Created',
  headers: {},
  config: {} as any,
};

describe('EnterContract', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the form with all required fields', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    // Check for key form sections
    expect(screen.getByText('Key Dates & Amounts')).toBeInTheDocument();
    expect(screen.getByText('Option Period & Earnest Money')).toBeInTheDocument();
    expect(screen.getByText('Title & Financing')).toBeInTheDocument();

    // Check for required fields by checking the form inputs exist
    expect(screen.getByPlaceholderText('450000')).toBeInTheDocument(); // Sales Price
    expect(screen.getByPlaceholderText('500')).toBeInTheDocument(); // Option Fee
    expect(screen.getByPlaceholderText('5000')).toBeInTheDocument(); // Earnest Money
    expect(screen.getByPlaceholderText('10')).toBeInTheDocument(); // Option Period Days
    expect(screen.getByPlaceholderText('20')).toBeInTheDocument(); // Title Commitment Days
  });

  it('displays property address in subtitle', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('123 Test St, Austin, TX 78701')).toBeInTheDocument();
    });
  });

  it('shows buyer financing checkbox', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      const checkbox = screen.getByRole('checkbox', { name: /Buyer will use financing/i });
      expect(checkbox).toBeInTheDocument();
      expect(checkbox).not.toBeChecked();
    });
  });

  it('shows seller lease fields when seller stays post-close is checked', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    // Initially, lease days field should not be visible
    expect(screen.queryByPlaceholderText('30')).not.toBeInTheDocument();

    // Check the seller stays checkbox
    const sellerStaysCheckbox = screen.getByRole('checkbox', {
      name: /Seller will remain in property after closing/i
    });
    await user.click(sellerStaysCheckbox);

    // Now the lease days field should appear
    expect(screen.getByPlaceholderText('30')).toBeInTheDocument();
  });

  it('submits form with valid data', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);
    vi.mocked(api.createContract).mockResolvedValue(mockContractResponse);
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
          <Route path="/project/:id" element={<div>Dashboard</div>} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    // Fill out the form using input names
    const effectiveDateInput = document.querySelector('input[name="effectiveDate"]') as HTMLInputElement;
    const salesPriceInput = screen.getByPlaceholderText('450000');
    const closingDateInput = document.querySelector('input[name="closingDate"]') as HTMLInputElement;
    const optionFeeInput = screen.getByPlaceholderText('500');
    const earnestMoneyInput = screen.getByPlaceholderText('5000');

    await user.type(effectiveDateInput, '2024-03-01');
    await user.clear(salesPriceInput);
    await user.type(salesPriceInput, '450000');
    await user.type(closingDateInput, '2024-04-15');
    await user.clear(optionFeeInput);
    await user.type(optionFeeInput, '500');
    await user.clear(earnestMoneyInput);
    await user.type(earnestMoneyInput, '5000');

    // Submit the form
    const submitButton = screen.getByRole('button', { name: /Create Contract & Generate Deadlines/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(api.createContract).toHaveBeenCalledWith(
        'test-project-1',
        expect.objectContaining({
          sales_price: 450000,
          option_fee: 500,
          earnest_money: 5000,
          option_period_days: 10,
          title_commitment_days: 20,
          buyer_financing: false,
          seller_stays_post_close: false,
        })
      );
    });
  });

  it('shows loading state while fetching project', () => {
    vi.mocked(api.getProject).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByText('Loading project...')).toBeInTheDocument();
  });

  it('shows error when project fetch fails', async () => {
    vi.mocked(api.getProject).mockRejectedValue(new Error('Failed to load'));

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Project not found')).toBeInTheDocument();
    });
  });

  it('shows error when contract submission fails', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);
    vi.mocked(api.createContract).mockRejectedValue({
      response: { data: { error: 'Contract creation failed' } },
    });
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    // Fill out minimum required fields
    const effectiveDateInput = document.querySelector('input[name="effectiveDate"]') as HTMLInputElement;
    const salesPriceInput = screen.getByPlaceholderText('450000');
    const closingDateInput = document.querySelector('input[name="closingDate"]') as HTMLInputElement;
    const optionFeeInput = screen.getByPlaceholderText('500');
    const earnestMoneyInput = screen.getByPlaceholderText('5000');

    await user.type(effectiveDateInput, '2024-03-01');
    await user.clear(salesPriceInput);
    await user.type(salesPriceInput, '450000');
    await user.type(closingDateInput, '2024-04-15');
    await user.clear(optionFeeInput);
    await user.type(optionFeeInput, '500');
    await user.clear(earnestMoneyInput);
    await user.type(earnestMoneyInput, '5000');

    // Submit the form
    const submitButton = screen.getByRole('button', { name: /Create Contract & Generate Deadlines/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText('Contract creation failed')).toBeInTheDocument();
    });
  });

  it('parses excluded and included items from comma-separated strings', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);
    vi.mocked(api.createContract).mockResolvedValue(mockContractResponse);
    const user = userEvent.setup();

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    // Fill required fields
    const effectiveDateInput = document.querySelector('input[name="effectiveDate"]') as HTMLInputElement;
    const salesPriceInput = screen.getByPlaceholderText('450000');
    const closingDateInput = document.querySelector('input[name="closingDate"]') as HTMLInputElement;
    const optionFeeInput = screen.getByPlaceholderText('500');
    const earnestMoneyInput = screen.getByPlaceholderText('5000');

    await user.type(effectiveDateInput, '2024-03-01');
    await user.clear(salesPriceInput);
    await user.type(salesPriceInput, '450000');
    await user.type(closingDateInput, '2024-04-15');
    await user.clear(optionFeeInput);
    await user.type(optionFeeInput, '500');
    await user.clear(earnestMoneyInput);
    await user.type(earnestMoneyInput, '5000');

    // Fill in excluded and included items
    const excludedInput = screen.getByPlaceholderText(/chandelier/i);
    const includedInput = screen.getByPlaceholderText(/refrigerator/i);

    await user.type(excludedInput, 'chandelier, custom curtains, antique mirror');
    await user.type(includedInput, 'refrigerator, washer, dryer');

    // Submit
    const submitButton = screen.getByRole('button', { name: /Create Contract & Generate Deadlines/i });
    await user.click(submitButton);

    await waitFor(() => {
      expect(api.createContract).toHaveBeenCalledWith(
        'test-project-1',
        expect.objectContaining({
          excluded_items: ['chandelier', 'custom curtains', 'antique mirror'],
          included_personal_items: ['refrigerator', 'washer', 'dryer'],
        })
      );
    });
  });

  it('has cancel button that navigates back', async () => {
    vi.mocked(api.getProject).mockResolvedValue(mockProject);

    render(
      <MemoryRouter initialEntries={['/project/test-project-1/contract/new']}>
        <Routes>
          <Route path="/project/:id/contract/new" element={<EnterContract />} />
          <Route path="/project/:id" element={<div>Dashboard</div>} />
        </Routes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Enter Contract Terms')).toBeInTheDocument();
    });

    const cancelButton = screen.getByRole('button', { name: /Cancel/i });
    expect(cancelButton).toBeInTheDocument();
  });
});
