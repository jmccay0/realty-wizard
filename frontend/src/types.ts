export interface Project {
  id: string;
  created_at: string;
  updated_at: string;
  property_address: string;
  seller_names: string[];
  seller_email: string;
  seller_phone: string;
  has_agent: boolean;
  agent_name?: string;
  title_company?: string;
  target_list_date?: string;
  status: 'setup' | 'listing_prep' | 'under_contract' | 'closing';
}

export interface Property {
  id: string;
  project_id: string;
  year_built: number;
  legal_description: string;
  tax_id: string;
  is_homestead: boolean;
  has_hoa: boolean;
  hoa_name?: string;
  hoa_fee?: number;
  hoa_frequency?: 'monthly' | 'quarterly' | 'annual';
  has_survey: boolean;
  survey_date?: string;
  survey_conditions?: string;
}

export interface Disclosure {
  id: string;
  project_id: string;
  updated_at: string;
  flooding_history: boolean;
  flooding_details?: string;
  in_flood_plain: boolean;
  flood_insurance_required: boolean;
  foundation_issues: boolean;
  foundation_details?: string;
  roof_age: number;
  roof_issues: boolean;
  roof_details?: string;
  plumbing_issues: boolean;
  plumbing_details?: string;
  electrical_issues: boolean;
  electrical_details?: string;
  hvac_age: number;
  hvac_issues: boolean;
  hvac_details?: string;
  lead_based_paint: boolean;
  insurance_claims_history: boolean;
  claims_details?: string;
  other_material_defects?: string;
}

export interface ContractTerms {
  id: string;
  project_id: string;
  effective_date: string;
  sales_price: number;
  closing_date: string;
  option_fee: number;
  option_period_days: number;
  earnest_money: number;
  title_commitment_days: number;
  buyer_financing: boolean;
  seller_stays_post_close: boolean;
  seller_lease_days?: number;
  excluded_items?: string[];
  included_personal_items?: string[];
}

export interface Deadline {
  id: string;
  project_id: string;
  type: string;
  description: string;
  due_date: string;
  is_business_day: boolean;
  priority: 'critical' | 'high' | 'normal';
  completed: boolean;
  notes?: string;
}

export interface Document {
  id: string;
  project_id: string;
  type: string;
  form_number: string;
  status: 'pending' | 'draft' | 'ready' | 'delivered';
  generated_at?: string;
  file_path?: string;
}

export interface ProjectSummary {
  project: Project;
  property?: Property;
  disclosure?: Disclosure;
  contract?: ContractTerms;
  deadlines?: Deadline[];
  documents?: Document[];
}
