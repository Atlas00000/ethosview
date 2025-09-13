import type { DashboardResponse, AnalyticsSummaryResponse, MarketLatestResponse } from "../types/api";
import NextDynamic from "next/dynamic";
import { api } from "../services/api";
import { HeroNew } from "../components/home/HeroNew";
import { MarketBar } from "../components/home/MarketBar";
import { FeaturedCarousel } from "../components/home/FeaturedCarousel";
const SectorHeatmap = NextDynamic(() => import("../components/home/SectorHeatmap").then(m => m.SectorHeatmap), { loading: () => <div className="max-w-6xl mx-auto px-4 py-8"><div className="glass-card p-6 skeleton h-40" /></div> });
const BusinessPreview = NextDynamic(() => import("../components/home/BusinessPreview").then(m => m.BusinessPreview), { loading: () => <div className="max-w-6xl mx-auto px-4 py-8"><div className="glass-card p-6 skeleton h-40" /></div> });
import { SymbolLookup } from "../components/home/SymbolLookup";
import { ESGHighlightsPro } from "../components/home/ESGHighlightsPro";
import { MarketSparkline } from "../components/home/MarketSparkline";
import { CorrelationTeaser } from "../components/home/CorrelationTeaser";
import { PELeaders } from "../components/home/PELeaders";
import { AlertsStrip } from "../components/home/AlertsStrip";
import { WSStatus } from "../components/home/WSStatus";
import { ESGTrendMini } from "../components/home/ESGTrendMini";
const AdvancedInsights = NextDynamic(() => import("../components/home/AdvancedInsights").then(m => m.AdvancedInsights), { loading: () => <div className="max-w-6xl mx-auto px-4 py-8"><div className="glass-card p-6 skeleton h-40" /></div> });
const FinancialSnapshot = NextDynamic(() => import("../components/home/FinancialSnapshot").then(m => m.FinancialSnapshot), { loading: () => <div className="max-w-6xl mx-auto px-4 py-8"><div className="glass-card p-6 skeleton h-40" /></div> });
import { SectorPie } from "../components/home/SectorPie";
import { ESGFeed } from "../components/home/ESGFeed";
import { ScrollReveal } from "../components/home/ScrollReveal";
import { QuickWidget } from "../components/home/QuickWidget";

export const dynamic = "force-dynamic";
export const revalidate = 0;

export default async function HomePage() {
  const today = new Date();
  const start = new Date(today.getTime() - 7 * 24 * 3600 * 1000);
  const fmt = (d: Date) => d.toISOString().slice(0, 10);
  // Load critical data first
  const [dashboard, analytics, market, history] = await Promise.all([
    api.dashboard().catch(() => ({ summary: { total_companies: 0, total_sectors: 0, avg_esg_score: 0 }, top_esg_scores: [], sectors: [], sector_stats: {} } as DashboardResponse)),
    api.analyticsSummary().catch(() => ({ summary: { total_companies: 0, total_sectors: 0, avg_esg_score: 0 }, sector_comparisons: [], top_esg_performers: [], top_market_cap: [], correlation: {} } as AnalyticsSummaryResponse)),
    api.marketLatest().catch(() => ({ market_data: { date: fmt(today) } } as MarketLatestResponse)),
    api.marketHistory(fmt(start), fmt(today), 30).catch(() => ({ start_date: fmt(start), end_date: fmt(today), data: [], count: 0 })),
  ]);

  // Load secondary data with delays to avoid rate limiting
  const secondaryData = await Promise.all([
    api.correlation().catch(() => ({ sample_size: 0, avg_esg_score: 0, avg_market_cap: 0, avg_pe_ratio: 0, avg_roe: 0, avg_profit_margin: 0, esg_market_cap_corr: 0, esg_pe_corr: 0, esg_roe_corr: 0, esg_profit_corr: 0 })),
    api.topPERatio(24).catch(() => ({ metric: "pe_ratio", top_performers: [], count: 0, limit: 24 })),
  ]);

  // Load low-priority data (these will gracefully fail)
  const lowPriorityData = await Promise.allSettled([
    api.alerts(false),
    api.wsStatus(),
    api.esgTrends(1, 30),
    api.advancedSummary(),
    api.financialIndicators(1),
    api.latestPrice(1),
    api.esgScores(20, 0, 0),
    api.companyFinancialSummary(1),
    api.stockPrices(1, 30),
  ]);

  const [corr, peTop] = secondaryData;
  const [alertsResult, wsResult, esgTrendsResult, advSummaryResult, finIndResult, finPriceResult, esgListResult, finSummaryResult, priceSeriesResult] = lowPriorityData.map(r => 
    r.status === 'fulfilled' ? r.value : null
  );

  const alerts = alertsResult || { alerts: [], count: 0, active_only: true };
  const ws = wsResult || { status: "degraded" };
  const esgTrends = esgTrendsResult || { company_id: 1, trends: [], count: 0, days: 30 };
  const advSummary = advSummaryResult || { summary: { portfolio_optimization: null, risk_summary: { total_companies_assessed: 0, average_risk_score: 0, risk_distribution: { low: 0, medium: 0, high: 0 } }, trend_summary: { esg_trends: { improving: 0, declining: 0, stable: 0 }, price_trends: { up: 0, down: 0, stable: 0 } }, message: "" } };
  const finInd = finIndResult || { company_id: 1, indicators: null };
  const finPrice = finPriceResult || { company_id: 1, price: null };
  const esgList = esgListResult || { scores: [], pagination: { limit: 20, offset: 0, count: 0 }, filters: { min_score: 0 } };
  const finSummary = finSummaryResult;
  const priceSeries = priceSeriesResult;

  return (
    <main>
      <ScrollReveal>
      <div id="hero" className="section-band band-a">
        <HeroNew dashboard={dashboard} analytics={analytics} market={market} history={history} />
      </div>
      <div id="market" className="section-band reveal-on-scroll">
        <div className="max-w-6xl mx-auto px-4">
          <div className="section-heading text-gradient"><span className="dot" /> Market</div>
        </div>
        <MarketBar market={market} />
        <div className="max-w-6xl mx-auto px-4 py-4 grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
          <div className="sm:col-span-2 glass-card p-3">
            <MarketSparkline history={history} />
          </div>
          <div className="glass-card p-3 flex items-center justify-center">
            <WSStatus status={ws} />
          </div>
        </div>
      </div>
      <div id="esg" className="section-band band-b reveal-on-scroll">
        <div className="max-w-6xl mx-auto px-4">
          <div className="section-heading text-gradient"><span className="dot" /> ESG</div>
        </div>
        <ESGHighlightsPro analytics={analytics} />
        <CorrelationTeaser corr={corr} />
        <ESGTrendMini trends={esgTrends} />
      </div>
      <div id="featured" className="section-band reveal-on-scroll">
        <div className="max-w-6xl mx-auto px-4">
          <div className="section-heading text-gradient"><span className="dot" /> Featured & P/E</div>
        </div>
        <FeaturedCarousel top={analytics.top_market_cap.slice(0, 24)} />
        <PELeaders top={peTop} />
      </div>
      <div className="section-band band-b">
        <div className="max-w-6xl mx-auto px-4">
          <div className="section-heading text-gradient"><span className="dot" /> Symbol Lookup</div>
        </div>
        <SymbolLookup />
      </div>
      <div id="sectors" className="section-band reveal-on-scroll">
        <div className="max-w-6xl mx-auto px-4">
          <div className="section-heading text-gradient"><span className="dot" /> Sectors</div>
        </div>
        <SectorHeatmap sectors={analytics.sector_comparisons} />
        <SectorPie sectors={analytics.sector_comparisons} />
      </div>
      <div className="section-band band-b reveal-on-scroll">
        <BusinessPreview dashboard={dashboard} />
      </div>
      <div className="section-band reveal-on-scroll">
        <AdvancedInsights summary={advSummary} />
      </div>
      <div className="section-band band-b reveal-on-scroll">
        <FinancialSnapshot ind={finInd} price={finPrice} summary={finSummary} series={priceSeries} />
      </div>
      <div className="section-band reveal-on-scroll">
        <ESGFeed list={esgList} />
      </div>
      <AlertsStrip alerts={alerts} />
      <QuickWidget />
      </ScrollReveal>
    </main>
  );
}
