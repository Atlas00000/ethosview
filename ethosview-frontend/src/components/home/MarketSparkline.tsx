"use client";
import React from "react";
import type { MarketHistoryResponse, CompanyPricesResponse } from "../../types/api";
import { LineChart, Line, ResponsiveContainer, AreaChart, Area, YAxis, Tooltip } from "recharts";

export function MarketSparkline({ history }: { history: MarketHistoryResponse }) {
  const rows = Array.isArray(history?.data) ? history.data : [];
  const data = rows
    .slice()
    .reverse()
    .map((d) => ({ x: d.date, y: d.sp500_close ?? d.nasdaq_close ?? d.dow_close ?? 0 }));
  return (
    <div style={{ width: "100%", height: 56 }}>
      <ResponsiveContainer>
        <LineChart data={data} margin={{ top: 4, bottom: 0, left: 0, right: 0 }}>
          <Line type="monotone" dataKey="y" stroke="#1E6AE1" strokeWidth={1.5} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

export function CompanySparkline({ series }: { series: CompanyPricesResponse | null }) {
  const rows = Array.isArray(series?.prices) ? series!.prices : [];
  const data = rows
    .slice()
    .reverse()
    .map((d) => ({ x: d.date, y: d.close_price }));
  const ys = data.map((d) => d.y);
  const min = ys.length ? Math.min(...ys) : 0;
  const max = ys.length ? Math.max(...ys) : 0;
  const pad = (max - min) > 0 ? (max - min) * 0.05 : (max || 1) * 0.02;
  const domain: [number, number] = [min - pad, max + pad];
  const gradId = `spark-${series?.company_id ?? 'x'}`;
  return (
    <div style={{ width: "100%", height: 44 }}>
      <ResponsiveContainer>
        <AreaChart data={data} margin={{ top: 0, bottom: 0, left: 0, right: 0 }}>
          <defs>
            <linearGradient id={gradId} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#1D9A6C" stopOpacity={0.45} />
              <stop offset="100%" stopColor="#1D9A6C" stopOpacity={0.05} />
            </linearGradient>
          </defs>
          <YAxis domain={domain} hide />
          <Tooltip
            cursor={{ stroke: "rgba(13, 18, 28, 0.15)", strokeWidth: 1 }}
            contentStyle={{ padding: 6, borderRadius: 8 }}
            labelFormatter={(l) => String(l)}
            formatter={(v: any) => [Number(v).toFixed(2), "Price"]}
          />
          <Area type="monotone" dataKey="y" stroke="#1D9A6C" strokeWidth={1.25} fill={`url(#${gradId})`} isAnimationActive />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}


