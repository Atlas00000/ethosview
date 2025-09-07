import React from "react";

export function WSStatus({ status }: { status: { status: string } }) {
  const color = status?.status === "running" ? "#1D9A6C" : status?.status === "degraded" ? "#FFB547" : "#E25555";
  return (
    <div className="flex items-center gap-2 text-sm" style={{ color: "#374151" }}>
      <span className="pulse-outline" style={{ background: color, width: 10, height: 10, borderRadius: 9999, display: "inline-block" }} />
      <span>Realtime</span>
    </div>
  );
}


