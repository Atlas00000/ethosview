"use client";
import React from "react";
import type { SectorComparison } from "../../types/api";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

export function SectorBar({ sectors }: { sectors: SectorComparison[] }) {
  const data = sectors.map((s) => ({ sector: s.sector, value: Number(s.avg_esg_score.toFixed(2)) }));
  return (
    <div style={{ width: "100%", height: 300 }}>
      <ResponsiveContainer>
        <BarChart data={data} margin={{ top: 8, right: 16, left: 0, bottom: 8 }}>
          <XAxis dataKey="sector" tick={{ fontSize: 11 }} interval={0} angle={-20} dy={10} height={60} />
          <YAxis tick={{ fontSize: 12 }} domain={[0, 100]} />
          <Tooltip cursor={{ fill: "rgba(0,0,0,0.04)" }} />
          <Bar dataKey="value" radius={[6, 6, 0, 0]} fill="#1D9A6C" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}


