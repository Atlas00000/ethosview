"use client";
import React, { useMemo, useState } from "react";
import type { DashboardResponse } from "../../types/api";
import { useCountUp } from "./useCountUp";

export function BusinessPreview({ dashboard }: { dashboard: DashboardResponse }) {
	const s = dashboard.summary;
	const [expanded, setExpanded] = useState(false);
	const totalCompanies = useCountUp(s.total_companies);
	const totalSectors = useCountUp(s.total_sectors);
	const avgESG = useCountUp(s.avg_esg_score);

	const sectorEntries = useMemo(() => Object.entries(dashboard.sector_stats || {}), [dashboard.sector_stats]);
	const sortedSectors = useMemo(() => sectorEntries.sort((a, b) => b[1] - a[1]), [sectorEntries]);
	const visibleSectors = useMemo(() => sortedSectors.slice(0, expanded ? 12 : 6), [sortedSectors, expanded]);
	const maxCount = useMemo(() => Math.max(1, ...sectorEntries.map(([, v]) => v as number)), [sectorEntries]);

	const topESG = dashboard.top_esg_scores?.slice(0, expanded ? 8 : 5) || [];

	return (
		<section className="max-w-6xl mx-auto px-4 py-8 relative overflow-hidden">
			<div className="blob" style={{ top: -80, left: -60, width: 220, height: 220, background: "rgba(30,106,225,0.10)", borderRadius: 9999 }} />
			<div className="flex items-center justify-between mb-3">
				<h2 className="text-xl font-semibold text-gradient">Business dashboard preview</h2>
				<button className="glass-card px-3 py-1 btn-sheen" onClick={() => setExpanded(v => !v)}>{expanded ? "Show fewer" : "Show more"}</button>
			</div>
			<div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
				{/* KPI tiles */}
				<div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
					<div className="text-sm" style={{ color: "#374151" }}>Total companies</div>
					<div className="text-2xl font-semibold" style={{ color: "#0B2545" }}>{Math.round(totalCompanies).toLocaleString()}</div>
				</div>
				<div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
					<div className="text-sm" style={{ color: "#374151" }}>Sectors</div>
					<div className="text-2xl font-semibold" style={{ color: "#0B2545" }}>{Math.round(totalSectors)}</div>
				</div>
				<div className="glass-card p-4 hover-lift tilt-hover animate-fade-in-up">
					<div className="text-sm" style={{ color: "#374151" }}>Avg ESG</div>
					<div className="text-2xl font-semibold" style={{ color: "#1D9A6C" }}>{avgESG.toFixed(2)}</div>
					<div className="mt-2 h-1.5 bg-white/30 rounded">
						<div className="h-1.5 rounded" style={{ width: `${Math.max(0, Math.min(100, (s.avg_esg_score)))}%`, background: '#1D9A6C', transition: 'width 320ms var(--ease-enter)' }} />
					</div>
				</div>
			</div>

			<div className="mt-4 grid grid-cols-1 lg:grid-cols-3 gap-4">
				{/* Top ESG list */}
				<div className="glass-card p-4 animate-fade-in-up">
					<div className="text-sm font-medium mb-2" style={{ color: "#0B2545" }}>Top ESG companies</div>
					<div className="divide-y">
						{topESG.map((it) => (
							<div key={it.id} className="py-2 flex items-center justify-between hover-lift">
								<div className="truncate pr-3">
									<div className="text-sm font-medium truncate" style={{ color: "#0B2545" }}>{it.company_name || `#${it.company_id}`}</div>
									<div className="text-xs" style={{ color: "#374151" }}>{it.company_symbol ? it.company_symbol : '—'} • {new Date(it.score_date).toISOString().slice(0,10)}</div>
								</div>
								<div className="text-sm font-semibold" style={{ color: "#1D9A6C" }}>{it.overall_score.toFixed(2)}</div>
							</div>
						))}
					</div>
				</div>

				{/* Sector distribution mini-list */}
				<div className="glass-card p-4 animate-fade-in-up lg:col-span-2">
					<div className="text-sm font-medium mb-2" style={{ color: "#0B2545" }}>Sectors by companies</div>
					<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
						{visibleSectors.map(([name, count]) => {
							const pct = (Number(count) / maxCount) * 100;
							return (
								<div key={name as string} className="glass-card p-2 hover-lift tilt-hover">
									<div className="flex items-center justify-between text-sm">
										<div className="font-medium" style={{ color: "#0B2545" }}>{name}</div>
										<div className="text-xs" style={{ color: "#374151" }}>{String(count)}</div>
									</div>
									<div className="mt-1 h-1.5 bg-white/30 rounded">
										<div className="h-1.5 rounded" style={{ width: `${Math.round(pct)}%`, background: '#1E6AE1', transition: 'width 320ms var(--ease-enter)' }} />
									</div>
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</section>
	);
}


