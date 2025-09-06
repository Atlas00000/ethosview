"use client";
import React from "react";
import type { SectorComparison } from "../../types/api";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";

const COLORS = ["#1E6AE1", "#2AB3A6", "#3986FF", "#1DAA8E", "#9AD36A", "#FFB547"];

export function SectorPie({ sectors }: { sectors: SectorComparison[] }) {
  const data = sectors.map((s) => ({ name: s.sector, value: s.company_count }));
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Sector distribution</h2>
      <div className="glass-card p-3" style={{ width: "100%", height: 260 }}>
        <ResponsiveContainer>
          <PieChart>
            <Pie data={data} dataKey="value" nameKey="name" innerRadius={60} outerRadius={100} paddingAngle={3}>
              {data.map((_, idx) => (
                <Cell key={idx} fill={COLORS[idx % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </section>
  );
}


