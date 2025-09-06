"use client";
import React from "react";
import type { MarketHistoryResponse } from "../../types/api";
import { LineChart, Line, ResponsiveContainer } from "recharts";

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


