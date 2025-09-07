"use client";
import React from "react";
import type { RiskAssessmentResponse } from "../../types/api";

export function RiskTeaser({ risk }: { risk: RiskAssessmentResponse }) {
  const level = (risk as any)?.assessment?.level || "unknown";
  const color = level === "low" ? "#1D9A6C" : level === "medium" ? "#FFB547" : level === "high" ? "#E25555" : "#6B7280";
  return (
    <section className="max-w-6xl mx-auto px-4 py-8">
      <h2 className="text-xl font-semibold mb-3 text-gradient">Risk assessment</h2>
      <div className="glass-card p-4 flex items-center gap-3 hover-lift tilt-hover animate-fade-in-up">
        <span className="pulse-outline" style={{ width: 10, height: 10, borderRadius: 9999, background: color, display: "inline-block" }} />
        <div className="text-sm" style={{ color: "#0B2545" }}>Level: {level}</div>
      </div>
    </section>
  );
}


