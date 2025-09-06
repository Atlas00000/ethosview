"use client";
import React from "react";
import type { ESGTrendsResponse } from "../../types/api";
import { LineChart, Line, ResponsiveContainer } from "recharts";

export function ESGTrendMini({ trends }: { trends: ESGTrendsResponse }) {
  const data = (trends?.trends || []).slice().reverse().map((t) => ({ x: t.date, y: t.esg_score }));
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">ESG trend</h2>
      <div className="glass-card p-3" style={{ width: "100%", height: 120 }}>
        <ResponsiveContainer>
          <LineChart data={data} margin={{ top: 4, bottom: 0, left: 0, right: 0 }}>
            <Line type="monotone" dataKey="y" stroke="#1D9A6C" strokeWidth={1.5} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </section>
  );
}


