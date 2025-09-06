"use client";
import React from "react";

const nav = [
  { href: "#hero", label: "Home" },
  { href: "#market", label: "Market" },
  { href: "#esg", label: "ESG" },
  { href: "#featured", label: "Featured" },
  { href: "#sectors", label: "Sectors" },
  { href: "/dashboard", label: "Dashboard" },
];

export function Header() {
  return (
    <header className="sticky top-0 z-40 backdrop-blur border-b" style={{ background: "rgba(255,255,255,0.7)", borderColor: "rgba(0,0,0,0.06)" }}>
      <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <a href="#hero" className="font-semibold text-gradient">EthosView</a>
        <nav className="flex items-center gap-5 text-sm">
          {nav.map((n) => (
            <a key={n.href} href={n.href} className="text-gray-700 hover:text-black transition-colors" style={{ color: "#374151" }}>
              {n.label}
            </a>
          ))}
        </nav>
      </div>
    </header>
  );
}


