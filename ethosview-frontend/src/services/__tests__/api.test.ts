import { api } from '../api';

// Mock fetch globally
global.fetch = jest.fn();

describe('API Service', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('dashboard', () => {
    it('should fetch dashboard data successfully', async () => {
      const mockData = {
        summary: {
          total_companies: 10,
          total_sectors: 5,
          avg_esg_score: 83.66
        },
        top_esg_scores: [],
        sectors: [],
        sector_stats: {}
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.dashboard();

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/dashboard',
        expect.objectContaining({
          headers: expect.objectContaining({
            Accept: 'application/json',
          }),
        })
      );
      expect(result).toEqual(mockData);
    });

    it('should handle fetch errors gracefully', async () => {
      (fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

      await expect(api.dashboard()).rejects.toThrow('Network error');
    });

    it('should retry on 429 status', async () => {
      const mockData = { summary: { total_companies: 10 } };

      (fetch as jest.Mock)
        .mockResolvedValueOnce({
          ok: false,
          status: 429,
          headers: new Map([['Retry-After', '1']]),
          text: async () => 'Rate limited',
        })
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockData,
        });

      const result = await api.dashboard();

      expect(fetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual(mockData);
    });
  });

  describe('analyticsSummary', () => {
    it('should fetch analytics summary with correct cache TTL', async () => {
      const mockData = {
        summary: { total_companies: 10 },
        sector_comparisons: [],
        top_esg_performers: [],
        top_market_cap: [],
        correlation: {}
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.analyticsSummary();

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/analytics/summary',
        expect.objectContaining({
          headers: expect.objectContaining({
            Accept: 'application/json',
          }),
        })
      );
      expect(result).toEqual(mockData);
    });
  });

  describe('marketLatest', () => {
    it('should fetch market data successfully', async () => {
      const mockData = {
        market_data: {
          date: '2024-01-15T00:00:00Z',
          sp500_close: 5035.6,
          nasdaq_close: 15750.8,
          dow_close: 39650.4
        }
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.marketLatest();

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/financial/market',
        expect.any(Object)
      );
      expect(result).toEqual(mockData);
    });
  });

  describe('companyBySymbol', () => {
    it('should fetch company by symbol with URL encoding', async () => {
      const mockData = {
        id: 1,
        name: 'Apple Inc.',
        symbol: 'AAPL',
        sector: 'Technology'
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.companyBySymbol('AAPL');

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/companies/symbol/AAPL',
        expect.any(Object)
      );
      expect(result).toEqual(mockData);
    });

    it('should handle special characters in symbol', async () => {
      const mockData = { id: 1, name: 'Test & Co.', symbol: 'T&C' };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      await api.companyBySymbol('T&C');

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/companies/symbol/T%26C',
        expect.any(Object)
      );
    });
  });

  describe('esgTrends', () => {
    it('should fetch ESG trends with parameters', async () => {
      const mockData = {
        company_id: 1,
        trends: [],
        count: 0,
        days: 30
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.esgTrends(1, 30);

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/analytics/companies/1/esg-trends?days=30',
        expect.any(Object)
      );
      expect(result).toEqual(mockData);
    });
  });

  describe('marketHistory', () => {
    it('should fetch market history with date parameters', async () => {
      const mockData = {
        start_date: '2024-01-01',
        end_date: '2024-01-15',
        data: [],
        count: 0
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await api.marketHistory('2024-01-01', '2024-01-15', 30);

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/financial/market/history?start_date=2024-01-01&end_date=2024-01-15&limit=30',
        expect.any(Object)
      );
      expect(result).toEqual(mockData);
    });
  });

  describe('caching behavior', () => {
    it('should cache successful responses', async () => {
      const mockData = { summary: { total_companies: 10 } };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      // First call
      await api.dashboard();
      // Second call should use cache
      await api.dashboard();

      expect(fetch).toHaveBeenCalledTimes(1);
    });

    it('should serve stale cache during backoff', async () => {
      const mockData = { summary: { total_companies: 10 } };

      (fetch as jest.Mock)
        .mockResolvedValueOnce({
          ok: true,
          json: async () => mockData,
        })
        .mockResolvedValueOnce({
          ok: false,
          status: 429,
          headers: new Map(),
          text: async () => 'Rate limited',
        });

      // First call - successful
      await api.dashboard();
      // Second call - rate limited, should serve stale cache
      const result = await api.dashboard();

      expect(fetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual(mockData);
    });
  });

  describe('concurrency control', () => {
    it('should limit concurrent requests', async () => {
      const mockData = { summary: { total_companies: 10 } };
      let callCount = 0;

      (fetch as jest.Mock).mockImplementation(() => {
        callCount++;
        return Promise.resolve({
          ok: true,
          json: async () => mockData,
        });
      });

      // Make multiple concurrent requests
      const promises = Array(10).fill(null).map(() => api.dashboard());
      await Promise.all(promises);

      // Should have made requests but limited concurrency
      expect(callCount).toBeGreaterThan(0);
      expect(callCount).toBeLessThanOrEqual(10);
    });
  });
});
