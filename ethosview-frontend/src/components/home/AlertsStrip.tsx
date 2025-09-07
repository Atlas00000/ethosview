import React from "react";
import type { AlertsResponse } from "../../types/api";

export function AlertsStrip({ alerts }: { alerts: AlertsResponse }) {
  if (!alerts.count) return null;
  return (
    <section className="animated-gradient">
      <div className="max-w-6xl mx-auto px-4 py-2 text-white text-sm whitespace-nowrap overflow-hidden">
        <div className="inline-block animate-fade-in-up" style={{ animationDuration: "400ms" }}>
          <span className="font-medium">Active alerts: {alerts.count}</span>
          <span className="mx-3 opacity-75">•</span>
          <span className="opacity-90">Tap a card to view more</span>
        </div>
      </div>
    </section>
  );
}


