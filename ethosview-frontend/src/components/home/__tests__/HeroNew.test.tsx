import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { HeroNew } from '../HeroNew';
import type { DashboardResponse, AnalyticsSummaryResponse, MarketLatestResponse, MarketHistoryResponse } from '../../../types/api';

// Mock the useCountUp hook
jest.mock('../useCountUp', () => ({
  useCountUp: (value: number) => value,
}));

// Mock Recharts components
jest.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div data-testid="responsive-container">{children}</div>,
  AreaChart: ({ children }: { children: React.ReactNode }) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  Tooltip: () => <div data-testid="tooltip" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  PieChart: ({ children }: { children: React.ReactNode }) => <div data-testid="pie-chart">{children}</div>,
  Pie: ({ children }: { children: React.ReactNode }) => <div data-testid="pie">{children}</div>,
  Cell: () => <div data-testid="cell" />,
  BarChart: ({ children }: { children: React.ReactNode }) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => <div data-testid="bar" />,
  RadialBarChart: ({ children }: { children: React.ReactNode }) => <div data-testid="radial-bar-chart">{children}</div>,
  RadialBar: () => <div data-testid="radial-bar" />,
}));

// Mock the API service
jest.mock('../../../services/api', () => ({
  api: {
    esgTrends: jest.fn().mockResolvedValue({
      company_id: 1,
      trends: [],
      count: 0,
      days: 30,
    }),
  },
}));

const mockDashboard: DashboardResponse = {
  summary: {
    total_companies: 10,
    total_sectors: 5,
    avg_esg_score: 83.66,
  },
  top_esg_scores: [
    {
      id: 1,
      company_id: 1,
      environmental_score: 85.5,
      social_score: 78.2,
      governance_score: 82.1,
      overall_score: 82.1,
      score_date: '2024-01-15T00:00:00Z',
      data_source: 'MSCI',
      created_at: '2024-01-15T00:00:00Z',
      updated_at: '2024-01-15T00:00:00Z',
      company_name: 'Apple Inc.',
      company_symbol: 'AAPL',
    },
  ],
  sectors: ['Technology', 'Healthcare', 'Financial Services'],
  sector_stats: {
    Technology: 3,
    Healthcare: 1,
    Financial Services: 1,
  },
};

const mockAnalytics: AnalyticsSummaryResponse = {
  summary: {
    total_companies: 10,
    total_sectors: 5,
    avg_esg_score: 78.66,
  },
  sector_comparisons: [
    {
      sector: 'Technology',
      company_count: 3,
      avg_esg_score: 82.23,
      avg_pe_ratio: 35.37,
      avg_market_cap: 2316666666666.67,
      total_market_cap: 6950000000000,
      best_esg_company: 'Microsoft Corporation',
      worst_esg_company: 'Alphabet Inc.',
    },
  ],
  top_esg_performers: [
    {
      company_id: 2,
      company_name: 'Microsoft Corporation',
      metric: 'ESG Score',
      value: 85.2,
      rank: 1,
      total_count: 10,
      percentile: 100,
      date: '2024-01-15T00:00:00Z',
    },
  ],
  top_market_cap: [
    {
      company_id: 2,
      company_name: 'Microsoft Corporation',
      metric: 'Market Cap',
      value: 3250000000000,
      rank: 1,
      total_count: 10,
      percentile: 100,
      date: '2025-08-10T00:00:00Z',
    },
  ],
  correlation: {
    avg_esg_score: 78.66,
    avg_market_cap: 893000000000,
    avg_pe_ratio: 24.41,
    avg_profit_margin: 0.282,
    avg_roe: 0.218,
    esg_market_cap_corr: 0.5114,
    esg_pe_corr: 0.4377,
    esg_roe_corr: -0.5558,
    esg_profit_corr: -0.3903,
    sample_size: 10,
  },
};

const mockMarket: MarketLatestResponse = {
  market_data: {
    id: 1,
    date: '2025-08-10T00:00:00Z',
    sp500_close: 5035.6,
    nasdaq_close: 15750.8,
    dow_close: 39650.4,
    vix_close: 9.4,
    treasury_10y: 4,
    created_at: '2025-08-10T00:00:00Z',
    updated_at: '2025-08-10T00:00:00Z',
  },
};

const mockHistory: MarketHistoryResponse = {
  start_date: '2024-01-01',
  end_date: '2024-01-15',
  data: [
    {
      date: '2024-01-01',
      sp500_close: 5000,
      nasdaq_close: 15500,
      dow_close: 39000,
      vix_close: 10,
      treasury_10y: 4.1,
    },
  ],
  count: 1,
};

describe('HeroNew Component', () => {
  const defaultProps = {
    dashboard: mockDashboard,
    analytics: mockAnalytics,
    market: mockMarket,
    history: mockHistory,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders without crashing', () => {
    render(<HeroNew {...defaultProps} />);
    expect(screen.getByText('Welcome to EthosView, your real time ESG and financial intelligence hub.')).toBeInTheDocument();
  });

  it('displays dashboard summary data', () => {
    render(<HeroNew {...defaultProps} />);
    
    expect(screen.getByText('10')).toBeInTheDocument(); // Total companies
    expect(screen.getByText('5')).toBeInTheDocument(); // Total sectors
    expect(screen.getByText('83.66')).toBeInTheDocument(); // Average ESG score
  });

  it('displays sector information', () => {
    render(<HeroNew {...defaultProps} />);
    
    expect(screen.getByText('Technology')).toBeInTheDocument();
    expect(screen.getByText('Healthcare')).toBeInTheDocument();
    expect(screen.getByText('Financial Services')).toBeInTheDocument();
  });

  it('handles empty data gracefully', () => {
    const emptyProps = {
      dashboard: {
        summary: { total_companies: 0, total_sectors: 0, avg_esg_score: 0 },
        top_esg_scores: [],
        sectors: [],
        sector_stats: {},
      },
      analytics: {
        summary: { total_companies: 0, total_sectors: 0, avg_esg_score: 0 },
        sector_comparisons: [],
        top_esg_performers: [],
        top_market_cap: [],
        correlation: {
          avg_esg_score: 0,
          avg_market_cap: 0,
          avg_pe_ratio: 0,
          avg_profit_margin: 0,
          avg_roe: 0,
          esg_market_cap_corr: 0,
          esg_pe_corr: 0,
          esg_roe_corr: 0,
          esg_profit_corr: 0,
          sample_size: 0,
        },
      },
      market: mockMarket,
      history: mockHistory,
    };

    render(<HeroNew {...emptyProps} />);
    expect(screen.getByText('0')).toBeInTheDocument(); // Should show 0 for empty data
  });

  it('handles missing props gracefully', () => {
    const partialProps = {
      dashboard: mockDashboard,
      analytics: mockAnalytics,
      market: mockMarket,
      history: null,
    };

    render(<HeroNew {...partialProps} />);
    expect(screen.getByText('Welcome to EthosView, your real time ESG and financial intelligence hub.')).toBeInTheDocument();
  });

  it('displays market data correctly', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check if market data is displayed (the exact text might vary based on formatting)
    expect(screen.getByText(/5035\.6/)).toBeInTheDocument(); // S&P 500 value
  });

  it('shows ESG scores in the correct format', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for ESG score display
    expect(screen.getByText(/82\.1/)).toBeInTheDocument(); // Apple's ESG score
  });

  it('handles keyboard navigation', () => {
    render(<HeroNew {...defaultProps} />);
    
    const container = screen.getByRole('main') || document.body;
    
    // Test right arrow key
    fireEvent.keyDown(container, { key: 'ArrowRight' });
    // Test left arrow key
    fireEvent.keyDown(container, { key: 'ArrowLeft' });
    
    // Component should still be rendered
    expect(screen.getByText('Welcome to EthosView, your real time ESG and financial intelligence hub.')).toBeInTheDocument();
  });

  it('displays slide indicators', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Should have slide indicators (buttons)
    const indicators = screen.getAllByRole('button');
    expect(indicators.length).toBeGreaterThan(0);
  });

  it('handles slide navigation', () => {
    render(<HeroNew {...defaultProps} />);
    
    const indicators = screen.getAllByRole('button');
    if (indicators.length > 0) {
      fireEvent.click(indicators[0]);
      // Component should still be rendered after click
      expect(screen.getByText('Welcome to EthosView, your real time ESG and financial intelligence hub.')).toBeInTheDocument();
    }
  });

  it('displays charts when data is available', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check if chart containers are rendered
    expect(screen.getAllByTestId('responsive-container')).toHaveLength(4); // Multiple charts
  });

  it('formats numbers correctly', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for properly formatted large numbers
    expect(screen.getByText(/3,250,000,000,000/)).toBeInTheDocument(); // Market cap formatting
  });

  it('handles sector statistics display', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for sector statistics
    expect(screen.getByText('3')).toBeInTheDocument(); // Technology sector count
  });

  it('displays progress bars for ESG scores', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Look for progress bar elements (they should have data-width attributes)
    const progressBars = document.querySelectorAll('[data-width]');
    expect(progressBars.length).toBeGreaterThan(0);
  });

  it('handles analytics data display', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for analytics data
    expect(screen.getByText('Microsoft Corporation')).toBeInTheDocument();
    expect(screen.getByText(/85\.2/)).toBeInTheDocument(); // Top ESG performer score
  });

  it('displays correlation data', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for correlation data display
    expect(screen.getByText('10')).toBeInTheDocument(); // Sample size
    expect(screen.getByText(/78\.66/)).toBeInTheDocument(); // Average ESG score
  });

  it('handles market history data', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check if market history is used in charts
    expect(screen.getAllByTestId('area-chart')).toHaveLength(1);
  });

  it('displays company information correctly', () => {
    render(<HeroNew {...defaultProps} />);
    
    // Check for company names and symbols
    expect(screen.getByText('Apple Inc.')).toBeInTheDocument();
    expect(screen.getByText('AAPL')).toBeInTheDocument();
  });

  it('handles different data types gracefully', () => {
    const mixedProps = {
      dashboard: {
        ...mockDashboard,
        top_esg_scores: [], // Empty array
      },
      analytics: {
        ...mockAnalytics,
        sector_comparisons: [], // Empty array
      },
      market: mockMarket,
      history: {
        ...mockHistory,
        data: [], // Empty array
      },
    };

    render(<HeroNew {...mixedProps} />);
    expect(screen.getByText('Welcome to EthosView, your real time ESG and financial intelligence hub.')).toBeInTheDocument();
  });
});
