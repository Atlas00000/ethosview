import type { DashboardResponse, AnalyticsSummaryResponse, MarketLatestResponse } from "../types/api";
import { api } from "../services/api";
import { Hero } from "../components/home/Hero";
import { MarketBar } from "../components/home/MarketBar";
import { FeaturedCarousel } from "../components/home/FeaturedCarousel";
import { SectorHeatmap } from "../components/home/SectorHeatmap";
import { BusinessPreview } from "../components/home/BusinessPreview";
import { SymbolLookup } from "../components/home/SymbolLookup";
import { ESGHighlights } from "../components/home/ESGHighlights";
import { MarketSparkline } from "../components/home/MarketSparkline";
import { CorrelationTeaser } from "../components/home/CorrelationTeaser";
import { PELeaders } from "../components/home/PELeaders";
import { AlertsStrip } from "../components/home/AlertsStrip";
import { WSStatus } from "../components/home/WSStatus";
import { ESGTrendMini } from "../components/home/ESGTrendMini";
import { AdvancedInsights } from "../components/home/AdvancedInsights";
import { FinancialSnapshot } from "../components/home/FinancialSnapshot";
import { SectorPie } from "../components/home/SectorPie";
import { RiskTeaser } from "../components/home/RiskTeaser";
import { ESGFeed } from "../components/home/ESGFeed";

export const revalidate = 60;

export default async function HomePage() {
  const today = new Date();
  const start = new Date(today.getTime() - 7 * 24 * 3600 * 1000);
  const fmt = (d: Date) => d.toISOString().slice(0, 10);
  const [dashboard, analytics, market, history, corr, peTop, alerts, ws, esgTrends, advSummary, finInd, finPrice, risk, esgList]: [
    DashboardResponse,
    AnalyticsSummaryResponse,
    MarketLatestResponse,
    any,
    any,
    any,
    any,
    any,
    any,
    any,
    any,
    any,
    any
  ] = await Promise.all([
    api.dashboard(),
    api.analyticsSummary(),
    api.marketLatest(),
    api.marketHistory(fmt(start), fmt(today), 30),
    api.correlation(),
    api.topPERatio(10),
    api.alerts(false),
    api.wsStatus(),
    api.esgTrends(1, 30),
    api.advancedSummary(),
    api.financialIndicators(1),
    api.latestPrice(1),
    api.riskAssessment(1),
    api.esgScores(10, 0),
  ]);

  return (
    <main>
      <div id="hero">
        <Hero dashboard={dashboard} analytics={analytics} />
      </div>
      <div id="market">
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
      <div id="esg">
        <ESGHighlights top={analytics.top_esg_performers} />
        <CorrelationTeaser corr={corr} />
        <ESGTrendMini trends={esgTrends} />
      </div>
      <div id="featured">
        <FeaturedCarousel top={analytics.top_market_cap.slice(0, 10)} />
        <PELeaders top={peTop} />
      </div>
      <SymbolLookup />
      <div id="sectors">
        <SectorHeatmap sectors={analytics.sector_comparisons} />
        <SectorPie sectors={analytics.sector_comparisons} />
      </div>
      <BusinessPreview dashboard={dashboard} />
      <AdvancedInsights summary={advSummary} />
      <FinancialSnapshot ind={finInd} price={finPrice} />
      <RiskTeaser risk={risk} />
      <ESGFeed list={esgList} />
      <AlertsStrip alerts={alerts} />
    </main>
  );
}
